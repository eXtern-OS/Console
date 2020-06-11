package server

import (
	"../publisher"
	"../web"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)
import "../auth"

func SetCompany(r *gin.Engine) {
	r.GET("/company", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			if t, uid := auth.AuthenticateCookie(cid); t {
				c.HTML(http.StatusOK, "company.html", web.RenderCompanyPage(uid, ""))
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	})

	r.POST("/company", func(c *gin.Context) {
		if cid, err := c.Cookie("devid"); err == nil {
			if t, uid := auth.AuthenticateCookie(cid); t {
				cname := c.PostForm("cname")
				cmail := c.PostForm("cmail")
				caddr := c.PostForm("caddr")
				cweb := c.PostForm("cweb")
				cIcon, errI := c.FormFile("cico")
				cCover, errC := c.FormFile("ccover")

				withIcon, withCover := errI == nil, errC == nil
				fmt.Println(cname, cmail, caddr, cweb, withIcon, withCover, errI, errC)
				c.HTML(http.StatusOK, "company.html", web.RenderCompanyPage(uid, publisher.CreateInfoUpdate(cname, cmail, caddr, cweb, uid, withIcon, withCover, cIcon, cCover, c)))
			}
		}
	})
}
