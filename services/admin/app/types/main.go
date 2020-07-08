package types

import "errors"

type Types interface {
	Getter
	Remover
}

type typeData struct {
	storage []*Type
}


type Type struct {
	ID      string
	Name    string
	IsCheck bool
}

func CreateType() Types {
	return &typeData{
		storage: []*Type{
			{"channel", "Канал", false},
			{"checkChannel", "Канал (проверка)", true},
			{"subscribeInstagram", "Instagram", false},
			{"likeInstagram", "Лайки Instagram", false},
			{"openWeb", "Ссылка", false},
			{"activateBot", "Бот", false},
		},
	}
}

type Getter interface {
	GetAllTypes() []*Type
	GetType(id string) (*Type, error)
	WithCheck(id string) bool
}

func (t *typeData) WithCheck(id string) bool {
	tp, err := t.GetType(id)
	if err != nil {
		return false
	}
	return tp.IsCheck
}

func (t *typeData) GetAllTypes() []*Type {
	return t.storage
}

func (t *typeData) GetType(id string) (*Type, error) {
	for _, tp := range t.storage {
		if tp.ID == id {
			return tp, nil
		}
	}
	return nil, errors.New(" Not found ")
}

type Remover interface {
	Delete(id string) error
}

func (t *typeData) Delete(id string) error {
	var newStorage []*Type
	var found bool

	for _, tp := range t.storage {
		if tp.ID != id {

			newStorage = append(newStorage, tp)

		} else {
			found = true
		}
	}

	if !found {
		return errors.New(" Not found ")
	}

	t.storage = newStorage

	return nil
}
