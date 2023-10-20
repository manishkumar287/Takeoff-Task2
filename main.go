package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	PhoneNo   string
	Role      string
	Salary    float64
}

func ReadEmployees() ([]Employee, error) {
	f, err := os.Open("./utils/employees.csv")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	var employees []Employee
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// fmt.Println("record 0 ",record[0])
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		sal, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, err
		}
		fmt.Println("record 0 ", id, sal)

		employee := Employee{
			ID:        id,
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			Password:  record[4],
			PhoneNo:   record[5],
			Role:      record[6],
			Salary:    sal,
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func ViewEmployees(c *gin.Context) {
	employees, err := ReadEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, employees)
}

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

func CreateEmployee(c *gin.Context) {
	// Get the employee data from the request body
	employee := Employee{}
	err := c.BindJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		err = w.Write([]string{"id", "first_name", "last_name", "email", "phone", "password", "role", "salary"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer w.Flush()
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
	err = w.Write([]string{fmt.Sprintf("%d", employee.ID), employee.FirstName, employee.LastName, employee.Email, employee.PhoneNo, employee.Password, employee.Role, fmt.Sprintf("%f", employee.Salary)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer w.Flush()

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully"})
}

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
	w := csv.NewWriter(f)

	// Read all the employee records
	records, err := r.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error at 3": err.Error()})
		return
	}

	// Create a new slice to store the updated employee records
	var newRecords []string
	for _, record := range records {
		idd, err := strconv.Atoi(record[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error  at 123": err.Error()})
			return
		}

		// If the employee ID does not match the given ID, add the record to the new slice
		if idd != id {
			newRecords = append(newRecords, record...)
		} else {
			// Update the employee data in the record
			record[1] = employee.FirstName
			record[2] = employee.LastName
			record[3] = employee.Email
			record[4] = employee.PhoneNo
			record[5] = employee.Password
			record[6] = employee.Role
			record[7] = strconv.FormatFloat(employee.Salary, 'f', 2, 64)

			// Add the updated record to the new slice
			newRecords = append(newRecords, record...)
		}
	}

	// Truncate the CSV file
	f.Truncate(0)
	f.Seek(0, 0)

	// Write the updated employee records to the CSV file
	for _, record := range newRecords {
		recordSlice := strings.Split(record, ",")
		err = w.Write(recordSlice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error at 4": err.Error()})
			return
		}
	}

	// Flush the CSV writer
	w.Flush()

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})

}

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
	w := csv.NewWriter(f)

	// Read all the employee records
	records, err := r.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new slice to store the updated employee records
	var newRecords []string
	for _, record := range records {
		if record[0] != id {
			newRecords = append(newRecords, record...)
		}
	}

	// Truncate the CSV file
	f.Truncate(0)
	f.Seek(0, 0)

	// Write the updated employee records to the CSV file
	for _, record := range newRecords {
		recordSlice := strings.Split(record, ",")
		err = w.Write(recordSlice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error at 4": err.Error()})
			return
		}
	}

	// Flush the CSV writer
	w.Flush()

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

func main() {
	router := gin.Default()
	router.GET("/employees", ViewEmployees)
	router.POST("/employees", CreateEmployee)
	router.GET("/search-employee", SearchEmployee)
	router.GET("/employee/:id", GetEmployeeByID)
	router.PUT("/employee/:id", UpdateEmployee)
	router.DELETE("/employee/:id", DeleteEmployee)

	router.Run("localhost:4040")
}
