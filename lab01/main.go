package main

import (
	"fmt"
	"log"
	"transportation/driver"
)

func main() {
	// // доступ к полям только через геттеры
	// // d.id или d.lastName - ошибка если вне пакета driver
	// var d driver.Driver

	// fmt.Println("ID:", d.ID())
	// fmt.Println("Фамилия:", d.LastName())

	// из строки
	d1, err := driver.NewDriver("1;Иванов;Иван;Иванович;+79956127674;23")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("из строки:", d1)

	// // из джейсона
	// jsonBytes := []byte(`{"id": 2, "last_name": "Петров", "first_name": "Пётр", "middle_name": "Петрович", "phone_number": "+79956127767", "experience": 10}`)
	// d2, err := driver.NewDriver(jsonBytes)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("из джейсона:", d2)

	// d3, _ := driver.NewDriver(2, "Петров", "Пётр", "Петрович", 10)
	// fmt.Println("прост конструктор:", d3.LastName(), d3.FirstName())

	// //7
	// // вывод
	// fmt.Println("полный вывод (String):", d1.FullPrint())
	// fmt.Println("краткий вывод (ShortString):", d1.ShortPrint())

	// // Сравнение
	// fmt.Println("d1 равен d2?", d1.Equals(d2)) // true
	// fmt.Println("d1 равен d3?", d1.Equals(d3)) // false
}
