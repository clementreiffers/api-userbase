package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getUser(c *gin.Context) {
	var searchUser User

	userDatabase := c.MustGet("userDatabase").(*mongo.Collection)
	if err := c.BindJSON(&searchUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cursor, err := userDatabase.Find(context.TODO(), bson.D{{}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertResult, err := userDatabase.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("insert result : %s", insertResult)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func connectMongo(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func main() {
	mongoUri := flag.String("mongoUri", "", "the uri to access to mongodb")

	flag.Parse()

	if *mongoUri == "" {
		panic("you need to specify the -mongoUri flag")
	}

	client := connectMongo(*mongoUri)
	userDatabase := client.Database("company").Collection("users")

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userDatabase", userDatabase)
		c.Next()
	})

	r.GET("/ping", ping)
	r.GET("/get-user", getUser)
	r.POST("/add-user", addUser)

	err := r.Run()
	if err != nil {
		fmt.Printf("%e\n", err)
	}
}
