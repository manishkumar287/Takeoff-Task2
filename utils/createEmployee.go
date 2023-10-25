package utils

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateEmployee(c *gin.Context) {
	// Get the employee data from the request body
	employee := Employee{}
	err := c.BindJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validation for FirstName
	if len(employee.FirstName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First Name is required"})
		return
	}

	// Validation for LastName
	if len(employee.LastName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Last Name is required"})
		return
	}

	// Validation for Email
	if !isValidEmail(employee.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email"})
		return
	}

	// Validation for PhoneNo (Assuming a valid phone number format)
	if len(employee.PhoneNo) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone Number is required"})
		return
	}

	// Validation for Password (Add more complex checks as needed)
	if len(employee.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password should be at least 8 characters long"})
		return
	}

	// Validation for Role (Add more specific checks as needed)
	if employee.Role != "admin" && employee.Role != "employee" && employee.Role != "manager" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Role please enter admin or employee or manager role only in lower case"})
		return
	}

	// Validation for Salary (Assuming a valid numeric format)
	if employee.Salary <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Salary"})
		return
	}

	// Check if the CSV file exists
	_, err = os.Stat("./utils/employees.csv")
	if err != nil && os.IsNotExist(err) {
		// Create the CSV file
		f, err := os.Create("./utils/employees.csv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		// Write the header row to the CSV file
		w := csv.NewWriter(f)
		err = w.Write([]string{"first_name", "last_name", "email", "phone", "password", "role", "salary"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer w.Flush()
	}

	// Check if the user already exists
	employees, err := ReadEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, emp := range employees {
		if emp.Email == employee.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
	}

	// Open the CSV file for writing
	f, err := os.OpenFile("./utils/employees.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Write the employee data to the CSV file
	w := csv.NewWriter(f)
	err = w.Write([]string{fmt.Sprintf("%d", time.Now().Nanosecond()), employee.FirstName, employee.LastName, employee.Email, employee.PhoneNo, employee.Password, employee.Role, fmt.Sprintf("%f", employee.Salary)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer w.Flush()

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully"})
}

func isValidEmail(email string) bool {
	// Define a regular expression for a basic email address format.
	// Note that this is a simple example, and email validation can be quite complex.
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	valid := regexp.MustCompile(pattern).MatchString(email)

	return valid
}
