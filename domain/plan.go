package domain

type PlanRepository interface {
	Store(p *Plan)
	Find(int64) *Plan
}

type Plan struct {
	ID    int64
	Name  string
	Limit int
	Price float64
}
