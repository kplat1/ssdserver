package main

import (
	"log"
	"net/http"
	"os"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/goki/mat32"
	_ "github.com/heroku/x/hmetrics/onload"
	"sync"
	// "io/ioutil"
	// "strconv"
)

type PlayerPosData struct {
	Username   string
	BattleName string
	Pos        mat32.Vec3
	Points     int
}

type FireEvent struct {
	Creator string
	Origin  mat32.Vec3
	Dir     mat32.Vec3
	Damage  int
	BattleName string
}

type PlayerPosMap map[string]*PlayerPosData

type FireEventMap map[int]*FireEvent

var TheBattleMaps map[string]PlayerPosMap

var TheFireEvents map[string]FireEventMap

var ServerMutex sync.Mutex

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	TheBattleMaps = make(map[string]PlayerPosMap)
	TheFireEvents = make(map[string]FireEventMap)
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("webFiles/*.html")
	router.Static("/static", "static")
	// PlayerPos["serverTest"] = &PlayerPosData{"serverTest", "testBattle", mat32.Vec3{1,1,1}, 5}

	router.GET("/website", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/playerPosPost", func(c *gin.Context) {
		ServerMutex.Lock()
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
		ppmap, ok := TheBattleMaps[jsonStruct.BattleName]
		if !ok || ppmap == nil {
			ppmap = make(PlayerPosMap)
			TheBattleMaps[jsonStruct.BattleName] = ppmap
		}
		ppmap[jsonStruct.Username] = jsonStruct
		TheBattleMaps[jsonStruct.BattleName] = ppmap
		ServerMutex.Unlock()
	})

	router.GET("/playerPosGet", func(c *gin.Context) {
		ServerMutex.Lock()
		battleName := c.Query("battleName")
		if battleName == "" {
			log.Printf("Didn't get battle name!")
			c.String(422, "text/text", "Did not get battle name, fail")
			ServerMutex.Unlock()
			return
		}
		ppmap, ok := TheBattleMaps[battleName]
		log.Printf("Pppmap: %v Ok: %v", ppmap, ok)
		if !ok || ppmap == nil {
			log.Printf("Battle maps nil")
			c.String(422, "text/text", "Battle map nil")
			ServerMutex.Unlock()
			return
		}
		c.JSON(http.StatusOK, ppmap)
		ServerMutex.Unlock()
	})

	router.POST("/fireEventsPost", func(c *gin.Context) {
		ServerMutex.Lock()
		jsonStruct := &FireEvent{}
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
		femap, ok := TheFireEvents[jsonStruct.BattleName]
		if !ok || femap == nil {
			femap = make(FireEventMap)
			TheFireEvents[jsonStruct.BattleName] = femap
		}
		femap[len(femap)] = jsonStruct
		TheFireEvents[jsonStruct.BattleName] = femap
		ServerMutex.Unlock()
	})

	router.GET("/fireEventsGet", func(c *gin.Context) {
		ServerMutex.Lock()
		battleName := c.Query("battleName")
		if battleName == "" {
			log.Printf("Didn't get battle name!")
			c.String(422, "text/text", "Did not get battle name, fail")
			ServerMutex.Unlock()
			return
		}
		femap, ok := TheFireEvents[battleName]
		log.Printf("Femap: %v Ok: %v", femap, ok)
		if !ok || femap == nil {
			log.Printf("Battle maps nil")
			c.String(422, "text/text", "Battle map nil")
			ServerMutex.Unlock()
			return
		}
		c.JSON(http.StatusOK, femap)
		ServerMutex.Unlock()
	})

	// router.POST("/playerPos", func(c *gin.Context) {
	//
	// })

	router.Run(":" + port)
}
