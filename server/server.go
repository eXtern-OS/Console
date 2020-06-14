package server

import (
	"externos.io/console/auth"
	"externos.io/console/publisher"
	"externos.io/console/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Init(r *gin.Engine, c Config) {
	SetRoot(r)
	SetApi(r)
	SetApp(r)
	SetAuth(r, c)
	SetCompany(r)
}

func SetRoot(r *gin.Engine) {
	r.GET("/alive", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

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
}
