package main

import (
	"task2/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/employees", utils.ViewEmployees)
	router.POST("/employee", utils.CreateEmployee)
	router.GET("/search-employee", utils.SearchEmployee)
	router.GET("/employee/:id", utils.GetEmployeeByID)
	router.PUT("/employee/:id", utils.UpdateEmployee)
	router.DELETE("/employee/:id", utils.DeleteEmployee)

	router.Run("localhost:4040")
}
