package main

import (
	"context"
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

	userDatabase := c.MustGet("userDatabase").(*mongo.Collection)
	if err := c.BindJSON(&searchUser); err != nil {
		return
	}
	//filter := bson.M{}

	cursor, err := userDatabase.Find(context.TODO(), bson.D{{}})
	//fmt.Printf("want to send: %+v\n", filteredUser)

	if err != nil {
		panic(err)
	}

	var userFound []bson.M
	if err = cursor.All(context.TODO(), &userFound); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, userFound)
}

func addUser(c *gin.Context) {
	userDatabase := c.MustGet("userDatabase").(*mongo.Collection)

	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	//users = append(users, newUser)
	insertResult, err := userDatabase.InsertOne(context.TODO(), newUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("insert result : %s", insertResult)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func deleteUser(c *gin.Context) {

}

func connectMongo(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func main() {
	client := connectMongo("mongodb://localhost:27017/test")
	userDatabase := client.Database("company").Collection("users")

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userDatabase", userDatabase)
		c.Next()
	})

	r.GET("/ping", ping)
	r.GET("/get-user", getUser)
	r.POST("/add-user", addUser)
	r.DELETE("/delete-User", deleteUser)

	err := r.Run()
	if err != nil {
		fmt.Printf("%e\n", err)
	}

	//defer func() {
	//	if err := client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//coll := client.Database("sample_mflix").Collection("movies")
	//title := "Back to the Future"
	//
	//var result bson.M
	//err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	//if errors.Is(err, mongo.ErrNoDocuments) {
	//	fmt.Printf("No document was found with the title %s\n", title)
	//	return
	//}
	//if err != nil {
	//	panic(err)
	//}
	//
	//jsonData, err := json.MarshalIndent(result, "", "    ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s\n", jsonData)
}
