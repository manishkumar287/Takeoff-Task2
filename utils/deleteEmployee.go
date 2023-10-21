package utils

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeleteEmployee(c *gin.Context) {
	// Get the employee ID from the request parameter
	id := c.Param("id")

	// Open the CSV file for reading and writing
	f, err := os.OpenFile("./utils/employees.csv", os.O_RDWR, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a buffer to store the updated content
	var updatedContentBuffer bytes.Buffer

	for _, record := range records {
		if record[0] != id {
			// Add the record to the buffer (without the deleted record)
			updatedContentBuffer.WriteString(strings.Join(record, ",") + "\n")
		}
	}

	// Truncate the CSV file
	f.Truncate(0)
	f.Seek(0, 0)

	// Write the updated content from the buffer back to the file
	_, err = f.WriteString(updatedContentBuffer.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error at 4": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
