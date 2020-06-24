package models

type Statistics interface {
	StatisticsGetter
	StatisticsSetter
}

type StatisticsData struct {
	All     int
	Blocked int
}

func CreateStatistics() *StatisticsData {
	return &StatisticsData{
		All:     0,
		Blocked: 0,
	}
}

type StatisticsGetter interface {
	GetAll() int
	GetBlocked() int
}

func (s *StatisticsData) GetAll() int {
	return s.All
}

func (s *StatisticsData) GetBlocked() int {
	return s.Blocked
}

type StatisticsSetter interface {
	SetAll(count int)
	SetBlocked(count int)
}

func (s *StatisticsData) SetAll(count int) {
	s.All = count
}

func (s *StatisticsData) SetBlocked(count int) {
	s.Blocked = count
}
