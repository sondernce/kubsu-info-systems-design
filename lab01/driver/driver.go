package driver

// Driver - полный класс сущности Водитель с инкапсулированными полями
type Driver struct {
	id         int
	lastName   string
	firstName  string
	middleName string
	experience int
}

// геттеры
func (d *Driver) ID() int {
	return d.id
}

func (d *Driver) LastName() string {
	return d.lastName
}

func (d *Driver) FirstName() string {
	return d.firstName
}

func (d *Driver) MiddleName() string {
	return d.middleName
}

func (d *Driver) Experience() int {
	return d.experience
}