package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
)

type ModelsHandler struct{}

func NewModelsHandler() *ModelsHandler { return &ModelsHandler{} }

func (h *ModelsHandler) List(c *gin.Context) {
	var items []model.MLModel
	database.GetDB().Find(&items)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *ModelsHandler) Create(c *gin.Context) {
	var m model.MLModel
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	database.GetDB().Create(&m)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *ModelsHandler) Update(c *gin.Context) {
	var m model.MLModel
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	id := c.Param("id")
	database.GetDB().Model(&model.MLModel{}).Where("id = ?", id).Updates(m)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *ModelsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	database.GetDB().Delete(&model.MLModel{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
