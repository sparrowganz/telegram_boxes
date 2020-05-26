package models

type Info interface {
	InfoSetter
	InfoGetter
}

type infoData struct {
	Txt string `bson:"text"`
	Img string `bson:"image"`
	Vid string `bson:"video"`
}

func CreateInfo() Info {
	return &infoData{
		Txt: "",
		Img: "",
		Vid: "",
	}
}

type InfoSetter interface {
	SetText(txt string)
	SetVideo(id string)
	SetImage(id string)
}

func (i *infoData) SetText(txt string) {
	i.Txt = txt
}

func (i *infoData) SetVideo(id string) {
	i.Vid = id
}

func (i *infoData) SetImage(id string) {
	i.Img = id
}

type InfoGetter interface {
	Text() string
	Video() string
	Image() string
}

func (i *infoData) Text() string {
	return i.Txt
}

func (i *infoData) Video() string {
	return i.Vid
}

func (i *infoData) Image() string {
	return i.Img
}
