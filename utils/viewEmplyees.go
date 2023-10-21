package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ViewEmployees(c *gin.Context) {
	employees, err := ReadEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, employees)
}
