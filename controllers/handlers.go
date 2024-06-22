package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ank809/File-Sharing-System/database"
	"github.com/ank809/File-Sharing-System/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func UploadFile(c *gin.Context) {
	fileHeader, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadGateway, "Cannot retrieve file")
		return
	}
	defer fileHeader.Close()
	password := c.PostForm("password")
	filename := c.PostForm("filename")

	log.Println(header.Filename)
	log.Println(header.Header)
	log.Println(header.Size)
	fmt.Println(password)
	log.Println("File retrieved successfully")

	var file model.File

	newpassword, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		fmt.Println("Error in generating hashed password")
		c.JSON(http.StatusInternalServerError, "Error generating hashed password")
		return
	}

	url := generateURL(header.Filename)
	file = model.File{
		ID:       primitive.NewObjectID(),
		Filename: filename,
		Password: string(newpassword),
		Url:      url,
	}

	collection := database.OpenCollection(database.Client, "Files")
	res, err := collection.InsertOne(context.TODO(), file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error inserting file metadata")
		return
	}
	fmt.Println(res)

	db := database.Client.Database("GridFS")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error creating GridFS bucket")
		return
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error opening upload stream")
		return
	}
	defer uploadStream.Close()
	fileSize, err := io.Copy(uploadStream, fileHeader)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error uploading file")
		return
	}

	log.Printf("File uploaded successfully. File size: %d bytes\n", fileSize)
	c.JSON(http.StatusOK, "File uploaded successfully")
}

func generateURL(filename string) string {
	return "http://localhost:8081/download?file=" + filename
}

func DownloadFile(c *gin.Context) {
	fileName := c.Query("file")
	password := c.Query("password")

	collection := database.OpenCollection(database.Client, "Files")
	var file model.File
	err := collection.FindOne(context.TODO(), bson.M{"filename": fileName}).Decode(&file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(file.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Incorrect Password")
		return
	}

	db := database.Client.Database("GridFS")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error creating GridFS bucket")
		return
	}
	var buf []byte
	dStream, err := bucket.OpenDownloadStreamByName(fileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error opening download stream")
		return
	}
	defer dStream.Close()
	buf, err = io.ReadAll(dStream)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Error reading file from GridFS")
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", buf)
	fmt.Println(buf)
}
