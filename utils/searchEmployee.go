package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SearchEmployee(c *gin.Context) {
	// Get the search criteria from the request parameters
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	email := c.Query("email")
	role := c.Query("role")

	// Read the employee data from the CSV file
	employees, err := ReadEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a slice to store the search results
	var searchResults []Employee

	// Iterate over the employee data and find the employees that match the search criteria
	for _, employee := range employees {
		if (firstName == "" || employee.FirstName == firstName) &&
			(lastName == "" || employee.LastName == lastName) &&
			(email == "" || employee.Email == email) &&
			(role == "" || employee.Role == role) {
			searchResults = append(searchResults, employee)
		}
	}

	// Return the search results in JSON format
	c.JSON(http.StatusOK, searchResults)
}
