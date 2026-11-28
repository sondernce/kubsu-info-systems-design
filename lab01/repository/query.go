package repository

// DriverFilter ограничивает выборку по точной фамилии и минимальному стажу.
type DriverFilter struct {
	LastName      string
	MinExperience int
}

type SortField uint8

const (
	SortByID SortField = iota
	SortByLastName
	SortByExperience
)

type DriverSort struct {
	Field      SortField
	Descending bool
}
