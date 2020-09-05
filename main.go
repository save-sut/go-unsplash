package main

import (
	"encoding/json"
	"github.com/peteretelej/jsonbox"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Photo struct {
	Id    string `json:"_id"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Great!"})
}

func GetAllPhoto(c *gin.Context) {
	dbKey := os.Getenv("MY_JSONBOX_KEY")
	cl, _ := jsonbox.NewClient("https://jsonbox.io/")

	out, err := cl.Read(dbKey)
	if err != nil {
		panic(err)
	}

	photosByte := []byte(out)
	photos := []Photo{}
	_ = json.Unmarshal(photosByte, &photos)

	c.JSON(200, gin.H{"Status": "SUCCESS", "Message": "Success", "Photos": photos})
}

func AddPhoto(c *gin.Context) {
	label := c.PostForm("label")
	if label == "" {
		c.JSON(400, gin.H{"Status": "ERROR", "Message": "Missing label"})
		return
	}

	url := c.PostForm("url")
	if url == "" {
		c.JSON(400, gin.H{"Status": "ERROR", "Message": "Missing url"})
		return
	}

	dbKey := os.Getenv("MY_JSONBOX_KEY")
	cl, _ := jsonbox.NewClient("https://jsonbox.io/")
	urlObj := []byte(`{"label": "` + label + `", "url": "` + url + `"}`)
	out, _ := cl.Create(dbKey, urlObj)
	if out == nil {
		c.JSON(400, gin.H{"Status": "ERROR", "Message": "Create this photo has error."})
		return
	} else {
		c.JSON(200, gin.H{"Status": "SUCCESS", "Message": "This photo has been created."})
		return
	}
}

func Search(c *gin.Context) {
	keyword := c.Param("keyword")

	dbKey := os.Getenv("MY_JSONBOX_KEY")
	cl, _ := jsonbox.NewClient("https://jsonbox.io/")
	out, err := cl.Read(``+dbKey+`?q=label:`+keyword+``)
	if err != nil {
		panic(err)
	}

	photosByte := []byte(out)
	photos := []Photo{}
	_ = json.Unmarshal(photosByte, &photos)

	c.JSON(200, gin.H{"Status": "SUCCESS", "Message": "Success", "Photos": photos})
}

func RemoveById(c *gin.Context) {
	password := c.Param("password")
	photoId := c.PostForm("photo_id")

	dbKey := os.Getenv("MY_JSONBOX_KEY")
	cl, _ := jsonbox.NewClient("https://jsonbox.io/")

	if password != "" && password == "password" {
		err := cl.Delete(dbKey, photoId)
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{"Status": "SUCCESS", "Message": "Deleted"})
		return
	} else {
		c.JSON(400, gin.H{"Status": "ERROR", "Message": "Wrong password"})
		return
	}
}

func ClearDB(c *gin.Context) {
	password := c.Param("password")

	dbKey := os.Getenv("MY_JSONBOX_KEY")
	cl, _ := jsonbox.NewClient("https://jsonbox.io/")

	if password != "" && password == "save-sut" {
		err := cl.DeleteAll(dbKey)
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{"Status": "SUCCESS", "Message": "All data has been deleted"})
		return
	} else {
		c.JSON(400, gin.H{"Status": "FAILED", "Message": "Wrong password or Password is not entered"})
		return
	}
}

func main()  {
	r := gin.Default()
	r.Use(cors.Default())
	api := r.Group("/api")
	{
		api.GET("/", Ping)
		api.GET("/photo/all", GetAllPhoto)
		api.POST("/photo/add", AddPhoto)
		api.GET("/photo/search/:keyword", Search)
		api.DELETE("photo/:password", RemoveById)
		api.GET("clear-db/:password", ClearDB)
	}
	r.Run()
}