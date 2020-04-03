package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goki/mat32"
	_ "github.com/heroku/x/hmetrics/onload"
	// "strconv"
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
	router.POST("/playerPosPost", func(c *gin.Context) {
		// // points, _ := strconv.Atoi(c.Param("points"))
		// // posX, _ := strconv.ParseFloat(c.Param("posX"), 32)
		// // posY, _ := strconv.ParseFloat(c.Param("posY"), 32)
		// // posZ, _ := strconv.ParseFloat(c.Param("posZ"), 32)
		// // PlayerPos[c.Param("username")] = &PlayerPosData{c.Param("username"), c.Param("battleName"), mat32.Vec3{float32(posX), float32(posY), float32(posZ)}, points}
		// // d := PlayerPos[c.Param("username")]
		// c.JSON(http.StatusOK, gin.H{"username": d.Username, "battleName": d.BattleName, "pos": d.Pos, "points": d.Points})
		PlayerPos["postTest"] = &PlayerPosData{"postTest", "hello", mat32.Vec3{1,1,1}, 2}
		d := PlayerPos["postTest"]
		c.JSON(http.StatusOK, gin.H{"username": d.Username, "battleName": d.BattleName, "pos": d.Pos, "points": d.Points})
	})

	router.GET("/playerPosGet", func(c *gin.Context) {
		for _, d := range PlayerPos {
			c.JSON(http.StatusOK, gin.H{"username": d.Username, "battleName": d.BattleName, "pos": d.Pos, "points": d.Points})
		}
	})

	// router.POST("/playerPos", func(c *gin.Context) {
	//
	// })

	router.Run(":" + port)
}
