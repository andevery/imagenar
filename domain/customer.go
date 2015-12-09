package domain

type CustomerRepository interface {
	Store(c *Customer)
	Find(id int64) *Customer
}

type Customer struct {
	ID   int64
	Name string
}
