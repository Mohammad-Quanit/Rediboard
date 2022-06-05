package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbInstance "github.com/mohammad-quanit/rediboard/db"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
)

// var database *dbInstance.Database

func main() {
	database, err := dbInstance.NewConnection(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	router := initRouter(database)
	router.Run(ListenAddr)
}

func initRouter(database *dbInstance.Database) *gin.Engine {
	r := gin.Default()
	r.GET("/", Hello)
	r.POST("/points", PostPoints)
	return r
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "Hello World!"})
}

func PostPoints(c *gin.Context) {
	var database dbInstance.Database
	var userJson dbInstance.User
	if err := c.ShouldBindJSON(&userJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := database.SaveUser(&userJson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userJson})
}
