package swagger

import (
	// stdlib
	"reflect"
	"strconv"
	"strings"
)

type (
	// Description the REST-method of endpoint
	Method struct {
		Description string   `json:"description,omitempty"`
		Consumes    []string `json:"consumes,omitempty"`
		Produces    []string `json:"produces,omitempty"`
		Summary     string   `json:"summary,omitempty"`
		OperationID string   `json:"operationId,omitempty"`
		// The parameters of requests
		Parameters []*Parameter `json:"parameters,omitempty"`
		// The endpoint responses
		Responses map[string]*Response `json:"responses,omitempty"`
		// Counters for same codes for multiple responses
		ResponseCode map[int]int `json:"-"`
	}

	// List of methods (GET, POST,...) for endpoint
	Methods map[string]IMethod
)

// First, programmer set Consumes/Produces/Response
type IMethod interface {
	// SetConsumes - sets the MIME types of accept data for the endpoint
	SetConsumes(c ...string) Producer
	// SetConsumes - sets the MIME types of return data for the endpoint
	SetProduces(p ...string) Consumer
	// AddResponse - adds a response
	AddResponse(сode int, description string, schema interface{}) Responser
}

type Responser interface {
	// AddParameter - adds a response
	AddResponse(сode int, description string, schema interface{}) Responser
}

type AdderInPathHeaderCookieQueryParameter interface {
	// AddInBodyParameter - adds a request in body parameter
	AddInPathParameter(name, description string, t reflect.Kind) AdderInPathHeaderCookieQueryParameter
	// AddInQueryParameter - adds a request in query parameter
	AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
	// AddInHeaderParameter - adds a request in header parameter
	AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
	// AddInCookieParameter - adds a request in header parameter
	AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
	Responser
}

type ParameterAndResponser interface {
	// AddParameter - adds a response
	AddResponse(Code int, description string, schema interface{}) Responser
	// AddInBodyParameter - adds a request in body parameter
	AddInBodyParameter(name, description string, t interface{}, required bool) AdderInPathHeaderCookieQueryParameter
	// AddInPathParameter - adds a request in path parameter
	AddInPathParameter(name, description string, t reflect.Kind) AdderInPathHeaderCookieQueryParameter
	// AddInQueryParameter - adds a request in query parameter
	AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
	// AddInHeaderParameter - adds a request in header parameter
	AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
	// AddInCookieParameter - adds a request in header parameter
	AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter
}

type Summarer interface {
	// SetDescription - sets a description of endpoint
	SetSummary(d string) ParameterAndResponser
}

type Consumer interface {
	// SetDescription - sets a description of endpoint
	SetDescription(d string) Summarer
}

type Producer interface {
	// SetProduces - sets the MIME types of return data for the endpoint
	SetProduces(p ...string) Consumer
	// SetDescription - sets a description of endpoint
	SetDescription(d string) Summarer
}

// NewMethod - create a new instance of the Method
func NewMethod() *Method {
	return &Method{
		Description: "Unnamed handler",
	}
}

func (m *Method) SetConsumes(c ...string) Producer {
	if m == nil {
		return nil
	}
	m.Consumes = c
	return m
}

func (m *Method) SetProduces(p ...string) Consumer {
	if m == nil {
		return nil
	}
	m.Produces = p
	return m
}

func (m *Method) SetDescription(d string) Summarer {
	if m == nil {
		return nil
	}
	m.Description = d
	return m
}

func (m *Method) SetSummary(s string) ParameterAndResponser {
	if m == nil {
		return nil
	}
	m.Summary = s
	return m
}

func (m *Method) AddInPathParameter(name, description string, t reflect.Kind) AdderInPathHeaderCookieQueryParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, &Parameter{
		BaseObject: BaseObject{
			Name:        name,
			Description: description,
			Type:        t,
		},
		Req: true,
		IN:  InPath,
	})
	return m
}

func (m *Method) AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, &Parameter{
		BaseObject: BaseObject{
			Name:        name,
			Description: description,
			Type:        t,
		},
		Req: required,
		IN:  InQuery,
	})
	return m
}

func (m *Method) AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, &Parameter{
		BaseObject: BaseObject{
			Name:        name,
			Description: description,
			Type:        t,
		},
		Req: required,
		IN:  InCookie,
	})
	return m
}

func (m *Method) AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInPathHeaderCookieQueryParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, &Parameter{
		BaseObject: BaseObject{
			Name:        name,
			Description: description,
			Type:        t,
		},
		Req: required,
		IN:  InHeader,
	})
	return m
}

func (m *Method) AddInBodyParameter(name, description string, t interface{}, required bool) AdderInPathHeaderCookieQueryParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, &Parameter{
		BaseObject: BaseObject{
			Name:        name,
			Description: description,
			Schema: &Schema{
				Type: t,
			},
		},
		Req: required,
		IN:  InBody,
	})
	return m
}

func (m *Method) AddResponse(сode int, description string, schema interface{}) Responser {
	if m == nil {
		return nil
	}
	if m.Responses == nil {
		m.Responses = make(map[string]*Response)
	}
	response := &Response{
		BaseObject: BaseObject{
			Description: description,
		},
	}

	if schema != nil {
		response.Schema = &Schema{
			Type: schema,
		}
	}

	if m.ResponseCode == nil {
		m.ResponseCode = make(map[int]int)
	}
	if _, ok := m.ResponseCode[сode]; !ok {
		m.ResponseCode[сode] = 1
		m.Responses[strconv.Itoa(сode)] = response
	} else {
		m.ResponseCode[сode]++
		m.Responses[strconv.Itoa(сode)+"("+strconv.Itoa(m.ResponseCode[сode])+")"] = response
	}

	return m
}

func (m *Method) Parse(path, methodName string, sw Doc) {
	// Parse parameters
	for _, p := range m.Parameters {
		p.Parse(sw)
	}
	// Parse responses
	for _, r := range m.Responses {
		r.Parse(sw)
	}
	if sw.Paths[path] == nil {
		sw.Paths[path] = make(Methods)
	}
	sw.Paths[path][strings.ToLower(methodName)] = m
}
