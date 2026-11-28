package main

import (
	"fmt"
	"log"

	"transportation/driver"
)

func main() {
	d, err := driver.NewDriver(1, "Иванов", "Иван", "Иванович", "+79161234567", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)
	mapper := driver.NewDriverMapper()
	short := mapper.ToShort(d) // указатель на общие данные
	fmt.Println(short)
	fmt.Println(short.Initials()) // И.И.
	if err := d.SetFirstName("Пётр"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(short.Initials()) // П.И. - изменилось само собой
	fmt.Println(d)
}
