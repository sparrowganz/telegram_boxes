package config

import "errors"

type KeyboardType string

var (
	Static KeyboardType = "static"
	Inline KeyboardType = "inline"
	Url    KeyboardType = "url"
)

type Row []Type

type Keyboard struct {
	Type KeyboardType `json:"type"`
	Rows []Row        `json:"rows"`
}

type Keyboards struct {
	Buttons map[Type]string     `json:"buttons"`
	Scheme  map[string]Keyboard `json:"scheme"`
}

type Result struct {
	Type  Type
	Text  string
	Value string
}

func (k *Keyboards) GetMain() (KeyboardType, [][]*Result) {
	return k.Scheme["main"].Type, k.getKeyboard(k.Scheme["main"].Rows)
}

func (k *Keyboards) GetCancel() (KeyboardType, [][]*Result) {
	return k.Scheme["cancel"].Type, k.getKeyboard(k.Scheme["cancel"].Rows)
}

func (k *Keyboards) GetOutput() (KeyboardType, [][]*Result) {
	return k.Scheme["output"].Type, k.getKeyboard(k.Scheme["output"].Rows)
}

func (k *Keyboards) GetFromType(tp string) (KeyboardType, [][]*Result) {
	return k.Scheme[tp].Type, k.getKeyboard(k.Scheme[tp].Rows)
}

func (k *Keyboards) GetNextTask() (KeyboardType, [][]*Result) {
	return k.Scheme["nextTask"].Type, k.getKeyboard(k.Scheme["nextTask"].Rows)
}

func (k *Keyboards) getKeyboard(rows []Row) (out [][]*Result) {
	for _, row := range rows {
		if len(row) > 0 {

			var resRow []*Result

			for _, but := range row {
				if val, ok := k.Buttons[but]; ok {
					resRow = append(resRow, &Result{Type: but, Text: val})
				}
			}
			out = append(out, resRow)
		}
	}
	return
}



func (k *Keyboards) GetTypeForName(name string) (Type, error) {
	for tp, val := range k.Buttons {
		if val == name {
			return tp, nil
		}
	}
	return "", errors.New(" Name not found ")
}
