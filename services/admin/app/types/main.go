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
	ID   string
	Name string
}

func CreateType() Types {
	return &typeData{
		storage: []*Type{
			{"channel", "Телеграмм подписка"},
			{"checkChannel", "Телеграмм подписка (check)"},
			{"subscribeInstagram", "Подписка инстаграм"},
			{"likeInstagram", "Лайки инстаграм"},
			{"openWeb","Открыть ссылку"},
			{"activateBot","Активировать бота"},
		},
	}
}

type Getter interface {
	GetAllTypes() []*Type
	GetType(id string) (*Type, error)
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
