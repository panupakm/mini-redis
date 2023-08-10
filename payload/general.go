package payload

type General struct {
	typ ValueType
	len uint32
	buf []byte
}

func NewGeneral(typ ValueType, buff []byte) *General {
	return &General{
		typ: typ,
		len: uint32(len(buff)),
		buf: buff,
	}
}
