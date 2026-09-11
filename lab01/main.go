package main

import (
	"fmt"
	"transportation/driver"
)

func main() {
	// доступ к полям только через геттеры
	// d.id или d.lastName - ошибка если вне пакета driver
	var d driver.Driver

	fmt.Println("ID:", d.ID())
	fmt.Println("Фамилия:", d.LastName())
}
