package swagger

import (
	// stdlib
	"reflect"
)

type Response struct {
	BaseObject
}

// Parse a response structure for JSON generation
func (r *Response) Parse(sw Doc) {
	ParseRootType(r, sw)
}

func (r *Response) GetSchema() *Schema {
	return r.Schema
}

func (r *Response) SetTypeName(typeName string) {
	r.TypeName = typeName
}

func (r *Response) GetType() reflect.Kind {
	return r.Type
}

func (r *Response) SetFormat(format string) {
	r.Format = format
}
