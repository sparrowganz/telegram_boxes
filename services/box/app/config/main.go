package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config interface {
	Getter
}

type Data struct {
	KeyboardsData Keyboards       `json:"keyboards"`
	CommandsData  map[string]Type `json:"commands"` // /tasks -> tasks
	TextsData     Texts           `json:"texts"`
	CountsData    Counts          `json:"counts"`
}

func CreateConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := &Data{}

	err = json.NewDecoder(f).Decode(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type Getter interface {
	Texts() Texts
	Counts() Counts
	Keyboards() Keyboards
	Commands() map[string]Type
	GetTypeForNameCommands(name string) (Type, error)
}

func (c *Data) Keyboards() Keyboards {
	return c.KeyboardsData
}

func (c *Data) Commands() map[string]Type {
	return c.CommandsData
}

func (c *Data) GetTypeForNameCommands(name string) (Type, error) {

	val, ok := c.CommandsData[name]
	if !ok {
		return "", errors.New(" Name not found ")
	}

	return val, nil
}

func (c *Data) Texts() Texts {
	return c.TextsData
}

func (c *Data) Counts() Counts {
	return c.CountsData
}

type Texts struct {
	Errors              Errors `json:"errors"`
	NotifyForReferral   string `json:"notifyForReferral"`
	StartText           string `json:"startText"`
	Balance             string `json:"balance"`
	Referrals           string `json:"referrals"`
	StatusReferralOk    string `json:"statusReferralOk"`
	StatusReferralFalse string `json:"statusReferralFalse"`
	Help                string `json:"help"`
}

type Errors struct {
	IncorrectCommand string `json:"incorrectCommand"`
	Error            string `json:"error"`
}

type Counts struct {
	CostForReferral int `json:"costForReferral"`
	VerifiedCount   int `json:"verifiedCount"`
}
