package models

type Flight struct {
	Origin      string
	Destination string
	DepartDate  string
	ReturnDate  string
	Price       int
	Gate        string
	FoundAt     string
	Currency    string
}
