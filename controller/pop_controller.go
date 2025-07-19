package controller

import (
	"net/http"
	"pop-calculator/model"
	"pop-calculator/service"

	"github.com/gin-gonic/gin"
)

func CalculatePoP(c *gin.Context) {
	var optionData model.PopRequest

	if err := c.ShouldBindJSON(&optionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Updated function call with new params
	pop := service.CalculatePoPValue(optionData.Spot, optionData.DaysToExpiry, optionData.Expiry, optionData.Symbol, optionData.OptionList)

	c.JSON(http.StatusOK, model.PopResponse{Pop: pop})
}
