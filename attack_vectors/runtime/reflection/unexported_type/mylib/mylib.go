package library

type unexportedStruct struct {
	field string
}

func NewUnexportedStruct(value string) *unexportedStruct {
	return &unexportedStruct{
		field: value,
	}
}
