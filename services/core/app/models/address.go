package models

type Address interface {
	AddressGetter
	AddressSetter
}

type AddressData struct {
	Ip      string `bson:"ip"`
	PortNum string `bson:"port"`
}

func CreateAddress(ip, port string) *AddressData {
	return &AddressData{
		Ip:      ip,
		PortNum: port,
	}
}

type AddressGetter interface {
	IP() string
	Port() string
	Addr() string
}

func (a *AddressData) IP() string {
	return a.Ip
}

func (a *AddressData) Port() string {
	return a.PortNum
}

func (a *AddressData) Addr() string {
	return a.Ip + a.PortNum
}

type AddressSetter interface {
	SetIP(ip string)
	SetPort(port string)
}

func (a *AddressData) SetIP(ip string) {
	a.Ip = ip
}

func (a *AddressData) SetPort(port string) {
	a.PortNum = port
}
