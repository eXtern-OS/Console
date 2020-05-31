package main

import (
	"./auth"
	"./db"
	"./publisher"
	"./web"
	"encoding/json"
	"fmt"
	beatrix "github.com/eXtern-OS/Beatrix"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	BeatrixToken   string `json:"beatrix-token"`
	BeatrixChannel string `json:"beatrix-channelID"`
	MongoURI       string `json:"mongo-uri"`
	CookieSecret   string `json:"cookie_secret"`
}

func LoadConfig() Config {
	var config Config
	configFile, err := os.Open("credentials.json")
	if err != nil {
		log.Panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	fmt.Println(config)
	err = configFile.Close()
	if err != nil {
		log.Panic(err)
	}
	return config
}

func Init(c Config) {
	db.Init(c.MongoURI)
	publisher.Init(c.MongoURI)
	auth.Init()
	beatrix.Init("CONSOLE", c.BeatrixToken, c.BeatrixChannel)
}

func main() {
	r := gin.Default()

	config := LoadConfig()

	Init(config)

	r.LoadHTMLGlob("static/*.html")
	r.Static("/assets", "./static/assets")
	store := cookie.NewStore([]byte(config.CookieSecret))
	r.Use(sessions.Sessions("devsession", store))

	r.GET("/", func(c *gin.Context) {
		if tid, err := c.Cookie("devid"); err == nil {
			log.Println(tid)
			c.HTML(http.StatusOK, "index.html", web.RenderIndex())
			return
		} else {
			fmt.Println("REDIRECTED", err)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
		return
	})

	r.POST("/login", func(c *gin.Context) {
		login := c.DefaultPostForm("username", "")
		password := c.DefaultPostForm("password", "")

		if login == "" || password == "" {
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		if t, uid := auth.GetUserIdByEmailAndPassword(login, password); t {
			c.SetCookie("devid", auth.NewCookie(uid), int(time.Now().Add(12*30*time.Hour).Unix()), "/", "localhost", false, false)
			log.Println("Set a cookie")
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}
	})

	r.GET("/app", func(c *gin.Context) {
		appId := c.DefaultQuery("appId", "")

		if appId == "" {
			c.Status(http.StatusNotFound)
			return
		}

		c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage())
		return
	})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_team.html", gin.H{})
	})

	r.GET("/newApp", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_app.html", gin.H{})
	})

	r.Group("/api")
	r.GET("/apps", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"App Name":            "abc",
			"App Icon url":        "abc.com",
			"App Downloads Count": "15",
			"App Revenue":         "12$",
			"App Version":         "v2.3",
		})
	})

	r.Run()
}
