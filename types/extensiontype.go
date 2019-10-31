package types

type ExtType string

func (kind ExtType) Code() uint16 {
	return extCodes[kind]
}

const (
	Void        ExtType = "void_t"
	ExtJsonType ExtType = "extension_json_type"
)

var extTypes = [...]ExtType{
	Void,
	ExtJsonType,
}

var extCodes map[ExtType]uint16

func init() {
	extCodes = make(map[ExtType]uint16, len(extTypes))
	for i, extType := range extTypes {
		extCodes[extType] = uint16(i)
	}
}
