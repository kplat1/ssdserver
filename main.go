package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/goki/mat32"
)

type PlayerPosData struct {
	Username   string
	BattleName string
	Pos        mat32.Vec3
	Points     int
}

var PlayerPos map[string]*PlayerPosData

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	PlayerPos = make(map[string]*PlayerPosData)
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("webFiles/*.html")
	router.Static("/static", "static")
	PlayerPos["serverTest"] = &PlayerPosData{"serverTest", "testBattle", mat32.Vec3{1, 1, 1}, 5}

	router.GET("/website", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/playerPos", func(c *gin.Context) {
		for _, d := range PlayerPos {
			c.JSON(http.StatusOK, gin.H{"username": d.Username, "battleName": d.BattleName, "pos": d.Pos, "points": d.Points})
		}
	})
	// router.POST("/playerPos", func(c *gin.Context) {
	//
	// })

	router.Run(":" + port)
}
