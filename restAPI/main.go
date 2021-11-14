package restAPI

import (
	"net/http"

    "github.com/gin-gonic/gin"
)



type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title:"Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "1", Title:"Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Kind Of Blue", Artist: "Miles Davis", Price: 120.23},
}

func Run() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func getAlbums(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, albums)
}

func postAlbums(ctx *gin.Context)  {
	var newAlbum album

	if err := ctx.ShouldBindJSON(&newAlbum); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  // H(elpful comment)
        return
    }

	albums = append(albums, newAlbum)
	ctx.JSON(http.StatusOK, gin.H{"data": newAlbum})
}