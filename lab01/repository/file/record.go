package file

// driverRecord задаёт одинаковые имена полей в JSON и YAML.
type driverRecord struct {
	ID          int    `json:"id" yaml:"id"`
	LastName    string `json:"last_name" yaml:"last_name"`
	FirstName   string `json:"first_name" yaml:"first_name"`
	MiddleName  string `json:"middle_name" yaml:"middle_name"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number"`
	Experience  int    `json:"experience" yaml:"experience"`
}
