package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	//err := r.Run()
	//if err != nil {
	//	fmt.Printf("%e\n", err)
	//}

	// make this follwing files connects to local mongodb

	uri := "mongodb://localhost:27017/test"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("sample_mflix").Collection("movies")
	title := "Back to the Future"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
