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

func checkUsernameIsLegal(username *string) (bool, error) {
	if match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", *username); !match {
		return false, fmt.Errorf("Username can only contain letters and numbers")
	}
	*username = strings.ToLower(strings.TrimSpace(*username))
	if !unicode.IsLetter(rune((*username)[0])) {
		return false, fmt.Errorf("Username must start with a letter")
	}
	if len(*username) < 4 || len(*username) > 16 {
		return false, fmt.Errorf("Username should be between 5 - 16 characters in length")
	}
	return true, nil
}

func (u *User) checkUserInfoIsLegal() error {
	if _, err := checkUsernameIsLegal(&u.Username); err != nil {
		return err
	}
	if len(u.Password) < 8 || len(u.Password) > 32 {
		return fmt.Errorf("Password should be between 8 - 32 characters in length")
	}
	if u.Role != "user" && u.Role != "" {
		return fmt.Errorf("No permission to create roles other than user")
	}

	u.Role = "user"
	u.Registration_date = time.Now().Unix()
	u.Points_balance = 0

	return nil
}

func (u *User) Insert() error {
	query := fmt.Sprintf(
		"INSERT INTO %s (username, password, registration_date, points_balance, role) VALUES (?, ?, ?, ?, ?)", USERTABLE)
	_, err := utils.DB.Exec(query, u.Username, u.Password, u.Registration_date, u.Points_balance, u.Role)
	return err
}

func checkUsernameExists(username string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = ?", USERTABLE)

	var count int
	rows, err := utils.DB.Query(query, username)
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

	if exists, err := checkUsernameExists(user.Username); err != nil {
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

func CheckUsernameHandler(c *gin.Context) {
	data := struct {
		Username string `json:"username"`
	}{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := checkUsernameIsLegal(&data.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exists, err := checkUsernameExists(data.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
}
