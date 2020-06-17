package types

import "errors"

type Types interface {
	Getter
}

type typeData struct {
	//todo remove debug structure
	storage []*Type
}

//todo remove debug structure
type Type struct {
	ID   string
	Name string
}


func CreateType() Types {
	//todo remove debug structure
	return &typeData{
		storage: []*Type{
			{"1", "TYPE 1"},
			{"2", "TYPE 2"},
			{"3", "TYPE 3"},
			{"4", "TYPE 4"},
		},
	}
}

type Getter interface {
	GetType(id string) (*Type, error)
}

func (t *typeData) GetType(id string) (*Type, error) {
	for _, tp := range t.storage {
		if tp.ID == id {
			return tp, nil
		}
	}
	return nil, errors.New(" Not found ")
}
