package types

type Kind uint8

const (
	Error Kind = iota

	Bool   // Boolean type
	Int    // Integer type
	Float  // Floating-point type
	String // String type
	Char   // Character type

	Unknown // Unknown type (type inference not yet performed)
	Unit    // Unit type
	Nominal // Nominal type (user-defined types)

	// Can be extended with more types as needed
)

// TypeInfo represents the type information of a value in the type system.
type TypeInfo struct {
	Kind Kind
	Name string
}

// NominalType returns a TypeInfo representing a nominal (user-defined) type with the given name.
func NominalType(name string) *TypeInfo {
	return &TypeInfo{
		Kind: Nominal,
		Name: name,
	}
}

func (t *TypeInfo) Equals(other *TypeInfo) bool {
	if t == nil || other == nil {
		return false
	}

	return t.Kind == other.Kind && t.Name == other.Name
}

func IsAssignable(from, to *TypeInfo) bool {
	if from == nil || to == nil {
		return false
	}

	if from.Kind == Error || to.Kind == Error {
		return false
	}

	if from.IsPrimitive() {
		return from.Kind == to.Kind
	}

	if from.Kind == Nominal && to.Kind == Nominal {
		return from.Name == to.Name
	}

	return false
}
