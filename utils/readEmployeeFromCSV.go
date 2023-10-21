package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadEmployees() ([]Employee, error) {
	
	f, err := os.Open("./utils/employees.csv")
	if err != nil {
		fmt.Println("error in opening the file")
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
		// fmt.Println("record 0 ", id, sal)

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
