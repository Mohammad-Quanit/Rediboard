package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammad-quanit/rediboard/db"
)

var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
)

// var database *dbInstance.Database

func main() {
	database, err := db.NewConnection(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	router := initRouter(database)
	router.Run(ListenAddr)
}

func initRouter(database *db.Database) *gin.Engine {
	r := gin.Default()
	r.GET("/", Hello)
	r.POST("/points", PostUserPoints)
	r.GET("/points/:username", GetUserPoints)
	r.GET("/leaderboard", GetLeaderboardPoints)
	return r
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "Hello World endpoint!"})
}

func PostUserPoints(c *gin.Context) {
	var database db.Database
	var userJson db.User
	if err := c.ShouldBindJSON(&userJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := database.SaveUser(&userJson)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userJson})
}

func GetUserPoints(c *gin.Context) {
	var database db.Database
	username := c.Param("username")
	user, err := database.GetUser(username)
	if err != nil {
		if err == db.ErrNil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + username})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetLeaderboardPoints(c *gin.Context) {
	var database db.Database
	leaderboard, err := database.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}
