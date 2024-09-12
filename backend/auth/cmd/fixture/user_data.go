package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/config"
	mysql_repo "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/repositories/mysql"
	service "gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/internal/services"
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
)

func main() {
	cfg := config.MustLoadConfig()
	conn := service.NewStorage(cfg.ConfigDatabase)

	officeRepo := mysql_repo.NewOfficeRepo(conn.Conn)
	officeService := service.NewOfficeService(officeRepo)
	userRepo := mysql_repo.NewUserRepo(conn.Conn)
	userService := service.NewUserService(userRepo)

	AddUsersFromCSV(userService, officeService)
}

func AddUsersFromCSV(userSrv service.UserService, officeSrv service.OfficeService) {
	fixturesPath := fetchConfigPath()
	users := ParseUserDataFromCSV(fixturesPath, officeSrv)

	for _, u := range users {
		err := userSrv.Create(&u)
		if err != nil {
			fmt.Println("Не удалось создать пользователя: ", err)
		} else {
			fmt.Println("Пользователей успешно создан: ", u)
		}
	}
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "path", "./storage/fixtures/UserData.csv", "path to file with user data") // Сначала пробуем считать путь из консоли
	flag.Parse()

	return res
}

func ParseUserDataFromCSV(path string, officeService service.OfficeService) []models.User {
	if path == "" {
		panic("Путь до файла с фикстурами пустой")
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

	var data []models.User

	for _, r := range records {
		var user models.User

		if r[0] == "Administrator" {
			user.RoleID = 1
		} else {
			user.RoleID = 2
		}

		user.Email = r[1]
		user.Password = r[2]
		user.FirstName = r[3]
		user.LastName = r[4]

		office, err := officeService.GetByTitle(r[5])
		if err != nil {
			panic(err)
		}
		user.OfficeID = office.ID

		parsedDate, err := time.Parse("1/2/2006", r[6])
		if err != nil {
			panic(err)
		}
		user.Birthday = parsedDate

		user.Active = r[7] == "1"

		data = append(data, user)
	}

	return data
}
