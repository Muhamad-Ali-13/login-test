package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password []byte // hashed
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
}

var users = map[string]User{}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// init dummy user
	password := "admin123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	users["admin"] = User{
		Username: "admin",
		Password: hash,
	}

	r.POST("/login", loginHandler)

	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "invalid request"})
		return
	}

	user, exists := users[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Message: "username not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{Message: "wrong password"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "login success"})
}
