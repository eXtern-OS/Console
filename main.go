package main

import (
	"./app"
	"./auth"
	"./db"
	"./publisher"
	"./web"
	"encoding/json"
	"fmt"
	"github.com/eXtern-OS/AMS"
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
	WebsiteURL     string `json:"host_url"`
	FlickrApi      string `json:"flickr-id"`
	FlickrSecret   string `json:"flickr-secret"`
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
	db.Init(c.MongoURI, c.FlickrApi, c.FlickrSecret)
	auth.Init()
	AMS.Init(c.MongoURI, "")
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

	// Those are needed paths for app icons and covers
	r.Static("/api/images/icons", "/temp/icons")
	r.Static("/api/images/covers", "/temp/covers")

	r.GET("/", func(c *gin.Context) {
		if tid, err := c.Cookie("devid"); err == nil {
			log.Println(tid)
			if t, uid := auth.AuthenticateCookie(tid); t {
				if x, _ := publisher.GetPublisherByUID(uid); x {
					c.HTML(http.StatusOK, "index.html", web.RenderIndex(uid))
					return
				} else {
					c.Redirect(http.StatusTemporaryRedirect, "/create")
					c.Abort()
					return
				}
			}
		}
		fmt.Println("REDIRECTED")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	})

	r.POST("/login", func(c *gin.Context) {
		login := c.DefaultPostForm("username", "")
		password := c.DefaultPostForm("password", "")

		if login == "" || password == "" {
			log.Println("Empty params")
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		log.Println(login, password)

		if t, uid := auth.GetUserIdByEmailAndPassword(login, password); t {
			c.SetCookie("devid", auth.NewCookie(uid), int(time.Now().Add(12*30*time.Hour).Unix()), "/", config.WebsiteURL, false, false)
			log.Println("Set a cookie")

			// Now let's check if he is a verified publisher

			if t, _ := publisher.GetPublisherByUID(uid); t {
				c.Redirect(http.StatusFound, "/")
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusFound, "/create")
				c.Abort()
			}
			return
		} else {
			log.Println("Error getting auth info")
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}
	})

	r.GET("/app", func(c *gin.Context) {
		appId := c.DefaultQuery("appId", "")

		if appId == "" {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}
		if cid, err := c.Cookie("devid"); err == nil {
			if t, uid := auth.AuthenticateCookie(cid); t {
				c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, ""))
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return

	})

	r.GET("/create", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			if t, _ := auth.AuthenticateCookie(cid); t {
				c.HTML(http.StatusOK, "create_team.html", gin.H{})
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	})

	r.POST("/create", func(c *gin.Context) {

		tname := c.PostForm("tname")
		tmail := c.PostForm("tmail")
		turl := c.PostForm("turl")
		taddr := c.PostForm("taddr")

		if tname == "" || tmail == "" || turl == "" || taddr == "" {
			fmt.Println(tname, tmail, turl, taddr)
			fmt.Println(c.Request.PostForm)
			c.Redirect(http.StatusFound, "/create")
			c.Abort()
			return
		} else {
			if tid, err := c.Cookie("devid"); err == nil {
				log.Println(tid)

				if t, uid := auth.AuthenticateCookie(tid); t {
					publisher.Create(tname, turl, taddr, tmail, uid)
					c.Redirect(http.StatusFound, "/")
					c.Abort()
					return
				} else {
					c.Status(http.StatusBadRequest)
					return
				}
			} else {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
		}
	})

	r.GET("/apps", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			if t, uid := auth.AuthenticateCookie(cid); t {
				c.JSON(http.StatusOK, web.RenderApplicationTables(uid))
				return
			}
		}
		log.Println("Error: failed to get cookie")
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/logout", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			auth.RemoveCookie(cid)
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	})

	r.GET("/newApp", func(c *gin.Context) {
		if tid, err := c.Cookie("devid"); err == nil {
			log.Println(tid)
			if t, uid := auth.AuthenticateCookie(tid); t {
				c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, "", ""))
				return
			}
		}
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	})

	r.POST("/newApp", func(c *gin.Context) {
		if tid, err := c.Cookie("devid"); err == nil {
			log.Println(tid)
			if t, uid := auth.AuthenticateCookie(tid); t {
				appname := c.PostForm("name")
				appdesc := c.PostForm("description")
				appicon, errI := c.FormFile("app_icon")
				appcover, errC := c.FormFile("app_cover")
				appversion := c.PostForm("app_version")
				appvdesc := c.PostForm("version_description")
				appvp, errP := c.FormFile("upload_package")
				appscreens := c.PostForm("screenshots")
				apppt := c.PostForm("pt")
				apppr := c.PostForm("price")

				if appname == "" || appdesc == "" || appversion == "" || appvdesc == "" {
					fmt.Println(appname, appdesc, appversion, appvdesc, appvp, appscreens, apppt, apppr)
					log.Println("SOMETHING IS MISSING")
					c.Redirect(http.StatusTemporaryRedirect, "/newApp")
					c.Abort()
					return
				}

				// App is free
				if apppr == "" {
					if errP != nil {
						log.Println("Upload is empty")
						c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, "Error uploading file, is it empty?", ""))
						c.Abort()
						return
					}
					resp := app.CreateFreeApp(appname, appdesc, appscreens, appversion, appvdesc, uid, appicon, appcover, appvp, c)
					if resp != "" {
						c.HTML(http.StatusOK, "create_app.html", web.RenderNewApplication(uid, resp, ""))
						c.Abort()
						return
					} else {
						c.Redirect(http.StatusFound, "/")
						c.Abort()
						return
					}
				} else { //App is paid
					log.Println("Create paid app")
				}

				if errI != nil {
					log.Println("------ERR ICON")
					log.Println(errI)
				}
				if errC != nil {
					log.Println("------ERR COVER")
					log.Println(errC)
				}
				fmt.Println(appname, appdesc, appversion, appvdesc, appvp, appscreens, apppt, apppr)
				return
			}
		}
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	})

	r.POST("/pushUpdate", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			if t, uid := auth.AuthenticateCookie(cid); t {
				appId := c.PostForm("appId")
				vIndex := c.PostForm("vindex")
				vDesc := c.PostForm("vdesc")
				vUp, errU := c.FormFile("vup")
				if errU != nil {
					log.Println(errU)
					c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, "Error processing file"))
					return
				}
				if publisher.VerifyPublisherOwnsApp(appId, uid) {
					resp := app.NewUpdate(uid, appId, vIndex, vDesc, vUp, c)
					log.Println(resp)
					if resp == "" {
						c.Redirect(http.StatusFound, "/")
						c.Abort()
						return
					} else {
						c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, resp))
						return
					}
				} else {
					c.HTML(http.StatusOK, "app.html", web.RenderApplicationPage(appId, uid, "We can't verify you own this app"))
					return
				}
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				c.Abort()
				return
			}
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	})

	r.Run(":80")
}
