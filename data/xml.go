package data

type AsXML interface {
	AsXML() []byte
}

type XMLSource []byte

func (x XMLSource) AsXML() []byte {
	return x
}
