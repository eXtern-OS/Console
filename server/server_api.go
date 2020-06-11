package server

import (
	"../auth"
	"../web"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SetApi(r *gin.Engine) {
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
}
