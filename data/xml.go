package data

type IAsXML interface {
	AsXML() []byte
}

type XMLSource []byte

func (x XMLSource) AsXML() []byte {
	return x
}
