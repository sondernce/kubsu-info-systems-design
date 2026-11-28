package main

import (
	"fmt"
	//"log"
	"transportation/driver"
)

func main() {

	d, _ := driver.NewDriver(1, "Иванов", "Иван", "Иванович", "+79161234567", 10)
	fmt.Println(d)
	mapper := driver.NewDriverMapper()
	short := mapper.ToShort(d) // указатель на общие данные
	fmt.Println(short)
	fmt.Println(short.Initials()) // И.И.
	d.SetFirstName("Пётр")        // пересчитывает initials внутри d.DriverShort
	fmt.Println(short.Initials()) // П.И. - изменилось само собой
	fmt.Println(d)

}
