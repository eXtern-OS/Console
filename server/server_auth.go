package server

import (
	"../auth"
	"../publisher"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func SetAuth(r *gin.Engine, config Config) {

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

	r.GET("/logout", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			auth.RemoveCookie(cid)
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
}
