package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type UserData struct {
	RoleID    int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Country   string
	Birthdate time.Time
	Active    bool
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "path", "./storage/fixtures/UserData.csv", "path to file with user data") // Сначала пробуем считать путь из консоли
	flag.Parse()

	return res
}

func ParseUserDataFromCSV() []UserData {
	path := fetchConfigPath()

	if path == "" {
		panic("Путь до конфига пустой")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Файл не найден")
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var data []UserData

	for i, r := range records {
		var user UserData

		if r[0] == "Administrator" {
			user.RoleID = 1
		} else {
			user.RoleID = 2
		}

		user.Email = r[1]
		user.Password = r[2]
		user.FirstName = r[3]
		user.LastName = r[4]
		user.Country = r[5]

		parsedDate, err := time.Parse("1/2/2006", r[6])
		if err != nil {
			fmt.Printf("Не удалось считать дату рождения в строке %d: %v\n", i, err.Error())
		}
		user.Birthdate = parsedDate

		user.Active = r[7] == "1"

		data = append(data, user)
	}

	return data
}

func main() {
	ParseUserDataFromCSV()
}
