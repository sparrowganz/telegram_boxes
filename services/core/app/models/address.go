package models

type Address interface {
	AddressGetter
	AddressSetter
}

type addressData struct {
	Ip      string `bson:"ip"`
	PortNum string `bson:"port"`
}

func CreateAddress(ip, port string) Address {
	return &addressData{
		Ip:      ip,
		PortNum: port,
	}
}

type AddressGetter interface {
	IP() string
	Port() string
	Addr() string
}

func (a *addressData) IP() string {
	return a.Ip
}

func (a *addressData) Port() string {
	return a.PortNum
}

func (a *addressData) Addr() string {
	return a.Ip + a.PortNum
}

type AddressSetter interface {
	SetIP(ip string)
	SetPort(port string)
}

func (a *addressData) SetIP(ip string) {
	a.Ip = ip
}

func (a *addressData) SetPort(port string) {
	a.PortNum = port
}
