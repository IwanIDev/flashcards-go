package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

type flashcard struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
	Front     string `json:"front"`
	Back      string `json:"back"`
}

type createFlashcardInput struct {
	Front string `json:"front" binding:"required"`
	Back  string `json:"back"  binding:"required"`
}

var db *gorm.DB

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})
	router.GET("/flashcards", getFlashcards)
	router.POST("/flashcards", postFlashcards)
	router.GET("/flashcards/:id", getCardByID)

	ConnectDatabase()

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&flashcard{})
	if err != nil {
		return
	}

	db = database
}

func getFlashcards(c *gin.Context) {
	var flashcards []flashcard
	db.Find(&flashcards)
	c.IndentedJSON(http.StatusOK, flashcards)
}

func postFlashcards(c *gin.Context) {
	var input createFlashcardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flashcard := flashcard{
		Front: input.Front,
		Back:  input.Back,
	}
	db.Create(&flashcard)

	c.JSON(http.StatusOK, flashcard)
}

func getCardByID(c *gin.Context) {
	var card flashcard

	if err := db.Where("id = ?", c.Param("id")).First(&card).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flashcard not found!"})
		return
	}

	c.JSON(http.StatusOK, card)
}
