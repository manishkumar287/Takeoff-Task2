package utils

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateEmployee(c *gin.Context) {
	// Get the employee ID from the request parameter
	d := c.Param("id")

	// Get the employee data from the request body
	employee := Employee{}
	err := c.BindJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error  at 1": err.Error()})
		return
	}

	id, err := strconv.Atoi(d)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error  at 12": err.Error()})
		return
	}

	// Open the CSV file for reading and writing
	f, err := os.OpenFile("./utils/employees.csv", os.O_RDWR, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error at 2": err.Error()})
		return
	}
	defer f.Close()

	// Create a CSV reader
	r := csv.NewReader(f)

	// Create a CSV writer
	// w := csv.NewWriter(f)

	// Read all the employee records
	records, err := r.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error at 3": err.Error()})
		return
	}

	// Create a buffer to store all updated records
	var updatedRecordsBuffer bytes.Buffer

	for _, record := range records {
		idd, err := strconv.Atoi(record[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error  at 123": err.Error()})
			return
		}

		// If the employee ID does not match the given ID, add the record to the buffer
		if idd != id {
			updatedRecordsBuffer.WriteString(strings.Join(record, ",") + "\n")
		} else {
			// Update the employee data in the record
			record[1] = employee.FirstName
			record[2] = employee.LastName
			record[3] = employee.Email
			record[4] = employee.PhoneNo
			record[5] = employee.Password
			record[6] = employee.Role
			record[7] = strconv.FormatFloat(employee.Salary, 'f', 2, 64)

			// Join the record elements with a comma
			updatedRecordsBuffer.WriteString(strings.Join(record, ",") + "\n")
		}
	}

	// Truncate the CSV file
	f.Truncate(0)
	f.Seek(0, 0)

	// Write the updated records from the buffer back to the file
	_, err = f.WriteString(updatedRecordsBuffer.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error at 4": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}
