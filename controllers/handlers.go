package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ank809/File-Sharing-System/database"
	"github.com/ank809/File-Sharing-System/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadGateway, "Cannot retrieve file")
		return
	}
	password := c.PostForm("password")

	log.Println(header.Filename)
	log.Println(header.Header)
	log.Println(header.Size)
	fmt.Println(password)
	log.Println("File retrieved successgully")

	var file model.File

	newpassword, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		fmt.Println("Error in generating hashed password")
	}
	url := generateURL(header.Filename)
	file = model.File{
		ID:       primitive.NewObjectID(),
		Filename: header.Filename,
		Password: string(newpassword),
		Url:      url,
	}
	collection := database.OpenCollection(database.Client, "Files")
	res, err := collection.InsertOne(context.TODO(), file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

}

func generateURL(filename string) string {
	return "http://localhost:8081/" + filename
}
