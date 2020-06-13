package config

import "errors"

type KeyboardType string

var (
	Static KeyboardType = "static"
	Inline KeyboardType = "inline"
)

type Row []Type

type Keyboard struct {
	Type KeyboardType `json:"type"`
	Rows []Row        `json:"rows"`
}

type Keyboards struct {
	Buttons map[Type]string `json:"buttons"`
	Main    Keyboard        `json:"main"`
	Cancel  Keyboard        `json:"toMainMenu"`
}

type Result struct {
	Type  Type
	Value string
}

func (k *Keyboards) GetMain() (KeyboardType, [][]Result) {

	var out [][]Result

	for _, row := range k.Main.Rows {
		if len(row) > 0 {

			var resRow []Result

			for _, but := range row {
				if val, ok := k.Buttons[but]; ok {
					resRow = append(resRow, Result{Type: but, Value: val})
				}
			}
			out = append(out, resRow)
		}
	}

	return k.Main.Type, out
}

func (k *Keyboards) GetCancel() (KeyboardType, [][]Result) {

	var out [][]Result

	for _, row := range k.Cancel.Rows {
		if len(row) > 0 {

			var resRow []Result

			for _, but := range row {
				if val, ok := k.Buttons[but]; ok {
					resRow = append(resRow, Result{Type: but, Value: val})
				}
			}
			out = append(out, resRow)
		}
	}

	return k.Cancel.Type, out
}

func (k *Keyboards) GetTypeForName(name string) (Type, error) {
	for tp, val := range k.Buttons {
		if val == name {
			return tp, nil
		}
	}
	return "", errors.New(" Name not found ")
}
