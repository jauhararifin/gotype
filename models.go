package gotype

type Type struct {
	PrimitiveType *PrimitiveType
	QualType      *QualType
	ChanType      *ChanType
	SliceType     *SliceType
	PtrType       *PtrType
	ArrayType     *ArrayType
	MapType       *MapType
	FuncType      *FuncType
	StructType    *StructType
	InterfaceType *InterfaceType
}

type PrimitiveKind string

const (
	PrimitiveKindBool       PrimitiveKind = "bool"
	PrimitiveKindByte       PrimitiveKind = "byte"
	PrimitiveKindInt        PrimitiveKind = "int"
	PrimitiveKindInt8       PrimitiveKind = "int8"
	PrimitiveKindInt16      PrimitiveKind = "int16"
	PrimitiveKindInt32      PrimitiveKind = "int32"
	PrimitiveKindInt64      PrimitiveKind = "int64"
	PrimitiveKindUint       PrimitiveKind = "uint"
	PrimitiveKindUint8      PrimitiveKind = "uint8"
	PrimitiveKindUint16     PrimitiveKind = "uint16"
	PrimitiveKindUint32     PrimitiveKind = "uint32"
	PrimitiveKindUint64     PrimitiveKind = "uint64"
	PrimitiveKindUintptr    PrimitiveKind = "uintptr"
	PrimitiveKindFloat32    PrimitiveKind = "float32"
	PrimitiveKindFloat64    PrimitiveKind = "float64"
	PrimitiveKindComplex64  PrimitiveKind = "complex64"
	PrimitiveKindComplex128 PrimitiveKind = "complex128"
	PrimitiveKindString     PrimitiveKind = "string"
	PrimitiveKindError      PrimitiveKind = "error"
)

type PrimitiveType struct {
	Kind PrimitiveKind
}

type QualType struct {
	Package string
	Name    string
}

type ChanTypeDir int

const (
	ChanTypeDirRecv = iota
	ChanTypeDirSend
	ChanTypeDirBoth
)

type ChanType struct {
	Dir  ChanTypeDir
	Elem Type
}

type SliceType struct {
	Elem Type
}

type PtrType struct {
	Elem Type
}

type ArrayType struct {
	Len  int
	Elem Type
}

type MapType struct {
	Key  Type
	Elem Type
}

type TypeField struct {
	Name string
	Type Type
}

type FuncType struct {
	Inputs     []TypeField
	Outputs    []TypeField
	IsVariadic bool
}

type StructType struct {
	Fields []TypeField
}

type InterfaceTypeMethod struct {
	Name string
	Func FuncType
}

type InterfaceType struct {
	Methods []InterfaceTypeMethod
}

func (t PrimitiveType) Type() Type { return Type{PrimitiveType: &t} }
func (t QualType) Type() Type      { return Type{QualType: &t} }
func (t ChanType) Type() Type      { return Type{ChanType: &t} }
func (t SliceType) Type() Type     { return Type{SliceType: &t} }
func (t PtrType) Type() Type       { return Type{PtrType: &t} }
func (t ArrayType) Type() Type     { return Type{ArrayType: &t} }
func (t MapType) Type() Type       { return Type{MapType: &t} }
func (t FuncType) Type() Type      { return Type{FuncType: &t} }
func (t StructType) Type() Type    { return Type{StructType: &t} }
func (t InterfaceType) Type() Type { return Type{InterfaceType: &t} }

func (t Type) IsPrimitive() bool { return t.PrimitiveType != nil }
func (t Type) IsQual() bool      { return t.QualType != nil }
func (t Type) IsChan() bool      { return t.ChanType != nil }
func (t Type) IsSlice() bool     { return t.SliceType != nil }
func (t Type) IsPtr() bool       { return t.PtrType != nil }
func (t Type) IsArray() bool     { return t.ArrayType != nil }
func (t Type) IsMap() bool       { return t.MapType != nil }
func (t Type) IsFunc() bool      { return t.FuncType != nil }
func (t Type) IsStruct() bool    { return t.StructType != nil }
func (t Type) IsInterface() bool { return t.InterfaceType != nil }

type TypeSpec struct {
	PackagePath string
	Name        string
}
