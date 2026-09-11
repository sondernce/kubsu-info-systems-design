package main

import (
	"fmt"
	"log"
	"transportation/driver"
)

func main() {
	// доступ к полям только через геттеры
	// d.id или d.lastName - ошибка если вне пакета driver
	var d driver.Driver

	fmt.Println("ID:", d.ID())
	fmt.Println("Фамилия:", d.LastName())

	// из строки
	d1, err := driver.NewDriverFromString("1;Иванов;Иван;Иванович;5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("из строки:", d1.LastName(), d1.FirstName())

	// из джейсона
	jsonBytes := []byte(`{"id": 2, "last_name": "Петров", "first_name": "Пётр", "middle_name": "Петрович", "experience": 10}`)
	d2, err := driver.NewDriverFromJSON(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("из джейсона:", d2.LastName(), d2.FirstName())
}
