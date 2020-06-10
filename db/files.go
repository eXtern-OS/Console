package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

func UploadFile(m *multipart.FileHeader, c *gin.Context, path string, ftype string) string {
	fmt.Println("We are uploading")

	if !(strings.Contains(m.Filename, ftype)) {
		fmt.Println(m.Filename, "NOT CONTAINS")
		return ""
	}
	pseudo := strconv.Itoa(int(time.Now().UnixNano())) + ftype
	rname := path + "/" + pseudo
	err := c.SaveUploadedFile(m, rname)
	if err != nil {
		return ""
	}
	return pseudo
}
