package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEmployeeByID(c *gin.Context) {
	// Get the employee ID from the request parameter
	_id := c.Param("id")

	id, err := strconv.Atoi(_id)
	if err != nil {
		return
	}

	// Find the employee with the given ID
	employee := FindEmployeeByID(id)

	// If the employee is not found, return a 404 Not Found error
	if employee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Employee not found with ID: %s", id)})
		return
	}

	// Return the employee details in JSON format
	c.JSON(http.StatusOK, employee)
}


func FindEmployeeByID(id int) *Employee {
	// Read the employee data from the CSV file
	employees, err := ReadEmployees()
	if err != nil {
		return nil
	}

	// Find the employee with the given ID
	for _, employee := range employees {
		if employee.ID == id {
			return &employee
		}
	}

	// If the employee is not found, return nil
	return nil
}
