package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var users []User

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func filterUser(users []User, userName string) []User {
	var filteredUser []User
	for _, user := range users {
		if user.Name == userName {
			filteredUser = append(filteredUser, user)
		}
	}
	return filteredUser
}

func getUser(c *gin.Context) {
	var searchUser User
	if err := c.BindJSON(&searchUser); err != nil {
		return
	}
	filteredUsers := filterUser(users, searchUser.Name)
	c.IndentedJSON(http.StatusOK, filteredUsers)
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func deleteUser(c *gin.Context) {

}

func main() {
	r := gin.Default()

	r.GET("/ping", ping)

	r.GET("/get-user", getUser)
	r.POST("/add-user", addUser)
	r.DELETE("/delete-User", deleteUser)

	err := r.Run()
	if err != nil {
		fmt.Printf("%e\n", err)
	}
}
