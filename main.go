package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

var dbMessage *redis.Client
var messageList []Message
var countDeleted int

// Database

func initDB() {
	address := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	dbMessage = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
}

func ping(c *gin.Context) {
	pong, err := dbMessage.Ping(context.TODO()).Result()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
			"time":    time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": pong,
		"time":    time.Now(),
	})
}

// Dummy Data

func createDummyData() {
	countMessage := len(messageList)

	message0 := Message{
		ID:       countMessage,
		Name:     "Winter Boo",
		Password: "super secret",
		Data:     "Semangatt!!!",
	}

	message1 := Message{
		ID:       countMessage + 1,
		Name:     "Summer Boo",
		Password: "super duper secret",
		Data:     "Yeayyy",
	}

	message2 := Message{
		ID:       countMessage + 2,
		Name:     "Boo",
		Password: "very secret",
		Data:     "Need sleep 😴",
	}

	messageList = append(messageList, message0, message1, message2)
}

// Type

type CreateMessageRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Data     string `json:"data"`
}

// Utils

func getMessageById(id int) (*Message, int, error) {
	for idx, message := range messageList {
		if message.ID == id {
			return &message, idx, nil
		}
	}
	return nil, -1, errors.New("Message not found!")
}

func isValidRequest(c *gin.Context) (*CreateMessageRequest, error) {
	name := c.PostForm("name")
	if name == "" {
		return nil, errors.New("The required parameter, name, is missing")
	}

	password := c.PostForm("password")
	if password == "" {
		return nil, errors.New("The required parameter, password, is missing")
	}

	data := c.PostForm("data")
	if data == "" {
		return nil, errors.New("The required parameter, data, is missing")
	}

	return &CreateMessageRequest{
		Name:     name,
		Password: password,
		Data:     data,
	}, nil
}

func isAUthenticated(oldMessage *Message, name string, password string) bool {
	if oldMessage.Name == name && oldMessage.Password == password {
		return true
	}
	return false
}

// Request Handler

func getAllMessageHandler(c *gin.Context) {
	c.JSON(200, gin.H{"data": messageList, "error": nil})
}

func getMessageByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"data": nil, "error": "ID must be a number!"})
		return
	}

	message, _, err := getMessageById(id)
	if err != nil {
		c.JSON(404, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": message, "error": nil})
}

func createMessageHandler(c *gin.Context) {
	request, err := isValidRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"data": nil, "error": err.Error()})
		return
	}

	newMessage := Message{
		ID:       len(messageList) + countDeleted,
		Name:     request.Name,
		Password: request.Password,
		Data:     request.Data,
	}

	messageList = append(messageList, newMessage)
	c.JSON(201, gin.H{"data": newMessage, "error": nil})
}

func updateMessageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"data": nil, "error": "ID must be a number!"})
		return
	}

	request, err := isValidRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"data": nil, "error": err.Error()})
		return
	}

	_, idx, err := getMessageById(id)
	if err != nil {
		c.JSON(404, gin.H{"data": nil, "error": err.Error()})
		return
	}

	message := &messageList[idx]

	if isAUthenticated(message, request.Name, request.Password) {
		message.Data = request.Data
		c.JSON(200, gin.H{"data": message, "error": nil})
		return
	}

	c.JSON(401, gin.H{"data": nil, "error": "Unauthorized - invalid credentials!"})
}

func deleteMessageHandler(c *gin.Context) {
	password := c.GetHeader("Authorization")
	if password == "" {
		c.JSON(400, gin.H{"data": nil, "error": "Missing authorization data!"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"data": nil, "error": "ID must be a number!"})
		return
	}

	message, idx, err := getMessageById(id)
	if err != nil {
		c.JSON(404, gin.H{"data": nil, "error": err.Error()})
		return
	}

	if message.Password != password {
		c.JSON(401, gin.H{"data": nil, "error": "Unauthorized - invalid credentials!"})
		return
	}

	messageList = append(messageList[:idx], messageList[idx+1:]...)
	countDeleted++
	c.JSON(200, gin.H{"data": messageList})
}

func main() {
	log.Println("LAW - Starting simple CRUD...")

	log.Println("Init database connection")
	initDB()

	log.Println("Creating Dummy Data")
	createDummyData()

	r := gin.Default()
	r.SetTrustedProxies([]string{"localhost", "herokuapp.com"})

	r.GET("/ping", ping)

	messageAPI := r.Group("/message")
	{
		messageAPI.GET("/list", getAllMessageHandler)
		messageAPI.GET("/:id", getMessageByIdHandler)
		messageAPI.POST("/create", createMessageHandler)
		messageAPI.PUT("/:id", updateMessageHandler)
		messageAPI.DELETE("/:id", deleteMessageHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
