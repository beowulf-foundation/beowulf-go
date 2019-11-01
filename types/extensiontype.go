package types

type ExtType string

func (kind ExtType) Code() uint16 {
	return ExtCodes[kind]
}

const (
	Void        ExtType = "void_t"
	ExtJsonType ExtType = "extension_json_type"
)

var extTypes = [...]ExtType{
	Void,
	ExtJsonType,
}

var ExtCodes map[ExtType]uint16

func init() {
	ExtCodes = make(map[ExtType]uint16, len(extTypes))
	for i, extType := range extTypes {
		ExtCodes[extType] = uint16(i)
	}
}

func GetExtCodes(s string) uint16 {
	switch s {
	case string(Void):
		return ExtCodes[Void]
	case string(ExtJsonType):
		return ExtCodes[ExtJsonType]
	}
	return ExtCodes[Void]
}
