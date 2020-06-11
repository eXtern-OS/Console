package server

import (
	"../app"
	"../auth"
	"../publisher"
	"../web"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SetApp(r *gin.Engine) {
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
}
