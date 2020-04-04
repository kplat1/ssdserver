package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goki/mat32"
	"encoding/json"
	_ "github.com/heroku/x/hmetrics/onload"
	// "io/ioutil"
	// "strconv"
)

type PlayerPosData struct {
	Username   string
	BattleName string
	Pos mat32.Vec3
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
	PlayerPos["serverTest"] = &PlayerPosData{"serverTest", "testBattle", mat32.Vec3{1,1,1}, 5}

	router.GET("/website", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/playerPosPost", func(c *gin.Context) {
		jsonStruct := &PlayerPosData{}
		decoder := json.NewDecoder(c.Request.Body)
		err := decoder.Decode(jsonStruct)
		if err != nil {
			log.Printf("Err: %v \n", err)
		}
		log.Printf("Struct: %v", jsonStruct)
		log.Printf("Battle Name: %v", jsonStruct.BattleName)
		// points, _ := strconv.Atoi(c.Param("points"))
		// posX, _ := strconv.ParseFloat(c.Param("posX"), 32)
		// posY, _ := strconv.ParseFloat(c.Param("posY"), 32)
		// posZ, _ := strconv.ParseFloat(c.Param("posZ"), 32)
		PlayerPos[jsonStruct.Username] = &PlayerPosData{jsonStruct.Username, jsonStruct.BattleName, jsonStruct.Pos, jsonStruct.Points}
		d := PlayerPos[jsonStruct.Username]
		log.Printf("Data: %v \n", d)
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
