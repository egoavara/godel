package godel

// Support
// + string
// + int<int8, int16, int32, int64>, uint<uint8, uint16, uint32, uint64>
// + []Identifier
type Identifier interface {
}
func validateIdentifier(identifier Identifier) Identifier {
	switch ident := identifier.(type) {
	case int:
	case string:
	case int8:
		return int(ident)
	case int16:
		return int(ident)
	case int32:
		return int(ident)
	case int64:
		return int(ident)
	case uint:
		return int(ident)
	case uint8:
		return int(ident)
	case uint16:
		return int(ident)
	case uint32:
		return int(ident)
	case uint64:
		return int(ident)
	default:
		return nil
	}
	return identifier
}