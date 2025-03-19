package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	if len(albums) == 0 {
		log.Warn("No albums found")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No albums found"})
		return
	}

	log.Info("Successfully retrieved albums")
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		log.Error("Failed to bind JSON: ", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	albums = append(albums, newAlbum)

	log.Infof("New album added: %+v", newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			log.Info("Album found: ", a)
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	log.Warn("Album not found with ID: ", id)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}
