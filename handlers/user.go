package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/daimaxiaofeng/user-management/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const USERTABLE = "users"

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type User struct {
	Id                uint64 `json:"id"`
	Username          string `json:"username" binding:"required"`
	Password          string `json:"password" binding:"required"`
	Registration_date int64  `json:"registration_date"`
	Points_balance    uint64 `json:"points_balance"`
	Role              string `json:"role"`
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	return err
}

// for this case, u.Password is hashed, password parameter is provided by user
func (u User) verifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) checkUserInfoIsLegal() error {
	if match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", u.Username); !match {
		return fmt.Errorf("Username can only contain letters and numbers")
	}
	u.Username = strings.ToLower(strings.TrimSpace(u.Username))
	if !unicode.IsLetter(rune(u.Username[0])) {
		return fmt.Errorf("Username must start with a letter")
	}
	if len(u.Username) < 5 || len(u.Username) > 16 {
		return fmt.Errorf("Username should be between 5 - 16 characters in length")
	}
	if len(u.Password) < 5 || len(u.Password) > 16 {
		return fmt.Errorf("Password should be between 5 - 16 characters in length")
	}
	if u.Role != "user" && u.Role != "" {
		return fmt.Errorf("No permission to create roles other than user")
	}

	u.Role = "user"
	u.Registration_date = time.Now().Unix()
	u.Points_balance = 0

	return nil
}

func (u *User) checkUsernameExists() (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = ?", USERTABLE)

	var count int
	rows, err := utils.DB.Query(query, u.Username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (u *User) Insert() error {
	query := fmt.Sprintf(
		"INSERT INTO %s (username, password, registration_date, points_balance, role) VALUES (?, ?, ?, ?, ?)", USERTABLE)
	_, err := utils.DB.Exec(query, u.Username, u.Password, u.Registration_date, u.Points_balance, u.Role)
	return err
}

func RegisterHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.checkUserInfoIsLegal(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if exists, err := user.checkUsernameExists(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	if err := user.hashPassword(); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if err := user.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
