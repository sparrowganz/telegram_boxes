package models

type Stats interface {
	StatsSetter
	StatsGetter
}

type StatsData struct {
	Success   int `json:"success"`
	WithError int `json:"withError"`
	All       int `json:"all"`
}

func CreateStats() Stats {
	return &StatsData{
		Success:   0,
		WithError: 0,
		All:       0,
	}
}

type StatsSetter interface {
	SetOK(count int)
	SetBad(count int)
	SetAll(count int)
}

func (s *StatsData) SetOK(count int) {
	s.Success += count
}

func (s *StatsData) SetBad(count int) {
	s.WithError += count
}

func (s *StatsData) SetAll(count int) {
	s.All += count
}

type StatsGetter interface {
	OK() int
	Bad() int
	Count() int
}

func (s *StatsData) OK() int {
	return s.Success
}

func (s *StatsData) Bad() int {
	return s.WithError
}

func (s *StatsData) Count() int {
	return s.All
}
