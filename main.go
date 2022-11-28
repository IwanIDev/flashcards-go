package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type flashcard struct {
	gorm.Model
	Front string `json:"front"`
	Back  string `json:"back"`
}

var flashcards []flashcard

func main() {
	router := gin.Default()
	router.GET("/flashcards", getFlashcards)
	router.POST("/flashcards", postFlashcards)
	router.GET("/flashcards/:id", getCardByID)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database!")
	}
	db.AutoMigrate(&flashcard{})
	router.Run("localhost:8080")
}

func getFlashcards(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, flashcards)
}

func postFlashcards(c *gin.Context) {
	var newCard flashcard

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newCard); err != nil {
		return
	}

	// Add the new album to the slice.
	flashcards = append(flashcards, newCard)
	c.IndentedJSON(http.StatusCreated, newCard)
}

func getCardByID(c *gin.Context) {
	jsonId := c.Param("id")
	id64, _ := strconv.ParseUint(jsonId, 10, 64)
	id := uint(id64)

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range flashcards {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "flashcard not found"})
}
