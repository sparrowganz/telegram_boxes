package models

type Data interface {
	DataGetter
	DataSetter
}

type dataData struct {
	TypeId int    `bson:"typeId"`
	Url    string `bson:"link"`
	Val    int    `bson:"count"`
}

func CreateData(typeID int, count int, link string) Data {
	return &dataData{
		TypeId: typeID,
		Url:    link,
		Val:    count,
	}
}

type DataGetter interface {
	TypeID() int
	Link() string
	Count() int
}

func (d *dataData) TypeID() int {
	return d.TypeId
}

func (d *dataData) Link() string {
	return d.Url
}

func (d *dataData) Count() int {
	return d.Val
}

type DataSetter interface {
	SetTypeID(id int)
	SetLink(link string)
	SetCount(val int)
}

func (d *dataData) SetTypeID(id int) {
	d.TypeId = id
}

func (d *dataData) SetLink(link string) {
	d.Url = link
}

func (d *dataData) SetCount(val int) {
	d.Val = val
}
