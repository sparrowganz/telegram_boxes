package models

type Statistics interface {
	StatisticsGetter
	StatisticsSetter
}

type StatisticsData struct {
	All     int64
	Blocked int64
}

func CreateStatistics() *StatisticsData {
	return &StatisticsData{
		All:     0,
		Blocked: 0,
	}
}

type StatisticsGetter interface {
	GetAll() int64
	GetBlocked() int64
}

func (s *StatisticsData) GetAll() int64 {
	return s.All
}

func (s *StatisticsData) GetBlocked() int64 {
	return s.Blocked
}

type StatisticsSetter interface {
	SetAll(count int64)
	SetBlocked(count int64)
}

func (s *StatisticsData) SetAll(count int64) {
	s.All = count
}

func (s *StatisticsData) SetBlocked(count int64) {
	s.Blocked = count
}
