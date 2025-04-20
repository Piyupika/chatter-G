package handlers

import (
	"net/http"

	"github.com/Piyu-Pika/godzilla-go/internal/database"
	"github.com/Piyu-Pika/godzilla-go/internal/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	userList := make([]struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}, len(users))
	for i, user := range users {
		userList[i] = struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
		}{user.ID, user.Username}
	}
	c.JSON(http.StatusOK, userList)
}

func SaveData(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Uuid     string `json:"uuid"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Uuid:     data.Uuid,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User saved successfully"})
}

func GetUserData(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
