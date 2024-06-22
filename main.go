package main

import (
	// "github.com/ank809/File-Sharing-System/model"
	"github.com/ank809/File-Sharing-System/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/upload", controllers.UploadFile)
	r.GET("/download", controllers.DownloadFile)
	r.Run(":8081")
}
