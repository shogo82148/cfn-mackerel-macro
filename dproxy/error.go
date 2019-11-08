package dproxy

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrorCode is type of errors
type ErrorCode int

const (
	// ErrorCodeType means expected type is not matched with actual.
	ErrorCodeType ErrorCode = iota + 1

	// ErrorCodeNotFound means key or index doesn't exist.
	ErrorCodeNotFound

	// ErrorCodeMapNorArray means target is not a map nor an array (for JSON Pointer)
	ErrorCodeMapNorArray

	// ErrorCodeConvertFailure means value conversion is failed.
	ErrorCodeConvertFailure

	// ErrorCodeInvalidIndex means token is invalid as index (for JSON Pointer)
	ErrorCodeInvalidIndex

	// ErrorCodeInvalidQuery means query is invalid as JSON Pointer.
	ErrorCodeInvalidQuery
)

func (et ErrorCode) String() string {
	switch et {
	case ErrorCodeType:
		return "ErrorCodeType"
	case ErrorCodeNotFound:
		return "ErrorCodeNotFound"
	case ErrorCodeMapNorArray:
		return "ErrorCodeMapNorArray"
	case ErrorCodeConvertFailure:
		return "ErrorCodeConvertFailure"
	case ErrorCodeInvalidIndex:
		return "ErrorCodeInvalidIndex"
	case ErrorCodeInvalidQuery:
		return "ErrorCodeInvalidQuery"
	default:
		return "Unknown"
	}
}

// Error get detail information of the error.
type Error interface {
	// ErrorCode returns type of error.
	ErrorCode() ErrorCode

	// FullAddress returns query string where cause first error.
	FullAddress() string
}

type errorProxy struct {
	errorCode ErrorCode
	parent    frame
	label     string

	expected Type
	actual   Type
	infoStr  string
}

// errorProxy implements error, Proxy and ProxySet.
var (
	_ error    = (*errorProxy)(nil)
	_ Proxy    = (*errorProxy)(nil)
	_ ProxySet = (*errorProxy)(nil)
)

func (p *errorProxy) Nil() bool {
	return false
}

func (p *errorProxy) Value() (interface{}, error) {
	return nil, p
}

func (p *errorProxy) Bool() (bool, error) {
	return false, p
}

func (p *errorProxy) OptionalBool() (*bool, error) {
	if p.errorCode == ErrorCodeNotFound {
		return nil, nil
	}
	return nil, p
}

func (p *errorProxy) Int64() (int64, error) {
	return 0, p
}

func (p *errorProxy) OptionalInt64() (*int64, error) {
	if p.errorCode == ErrorCodeNotFound {
		return nil, nil
	}
	return nil, p
}

func (p *errorProxy) Uint64() (uint64, error) {
	return 0, p
}
func (p *errorProxy) OptionalUint64() (*uint64, error) {
	if p.errorCode == ErrorCodeNotFound {
		return nil, nil
	}
	return nil, p
}

func (p *errorProxy) Float64() (float64, error) {
	return 0, p
}

func (p *errorProxy) OptionalFloat64() (*float64, error) {
	if p.errorCode == ErrorCodeNotFound {
		return nil, nil
	}
	return nil, p
}

func (p *errorProxy) String() (string, error) {
	return "", p
}

func (p *errorProxy) OptionalString() (*string, error) {
	if p.errorCode == ErrorCodeNotFound {
		return nil, nil
	}
	return nil, p
}

func (p *errorProxy) Array() ([]interface{}, error) {
	return nil, p
}

func (p *errorProxy) Map() (map[string]interface{}, error) {
	return nil, p
}

func (p *errorProxy) A(n int) Proxy {
	return p
}

func (p *errorProxy) M(k string) Proxy {
	return p
}

func (p *errorProxy) P(q string) Proxy {
	return p
}

func (p *errorProxy) Empty() bool {
	return true
}

func (p *errorProxy) Len() int {
	return 0
}

func (p *errorProxy) BoolArray() ([]bool, error) {
	return nil, p
}

func (p *errorProxy) Int64Array() ([]int64, error) {
	return nil, p
}

func (p *errorProxy) Float64Array() ([]float64, error) {
	return nil, p
}

func (p *errorProxy) StringArray() ([]string, error) {
	return nil, p
}

func (p *errorProxy) ArrayArray() ([][]interface{}, error) {
	return nil, p
}

func (p *errorProxy) MapArray() ([]map[string]interface{}, error) {
	return nil, p
}

func (p *errorProxy) ProxyArray() ([]Proxy, error) {
	return nil, p
}

func (p *errorProxy) ProxySet() ProxySet {
	return p
}

func (p *errorProxy) Q(k string) ProxySet {
	return p
}

func (p *errorProxy) Qc(k string) ProxySet {
	return p
}

func (p *errorProxy) findJPT(t string) Proxy {
	return p
}

func (p *errorProxy) parentFrame() frame {
	return p.parent
}

func (p *errorProxy) frameLabel() string {
	return p.label
}

func (p *errorProxy) Error() string {
	switch p.errorCode {
	case ErrorCodeType:
		return fmt.Sprintf("not matched types: expected=%s actual=%s: %s",
			p.expected.String(), p.actual.String(), p.FullAddress())
	case ErrorCodeNotFound:
		return fmt.Sprintf("not found: %s", p.FullAddress())
	case ErrorCodeMapNorArray:
		// FIXME: better error message.
		return fmt.Sprintf("not map nor array: actual=%s: %s",
			p.actual.String(), p.FullAddress())
	case ErrorCodeConvertFailure:
		return fmt.Sprintf("convert error: %s: %s", p.infoStr, p.FullAddress())
	case ErrorCodeInvalidIndex:
		// FIXME: better error message.
		return fmt.Sprintf("invalid index: %s: %s", p.infoStr, p.FullAddress())
	case ErrorCodeInvalidQuery:
		// FIXME: better error message.
		return fmt.Sprintf("invalid query: %s: %s", p.infoStr, p.FullAddress())
	default:
		return fmt.Sprintf("unexpected: %s", p.FullAddress())
	}
}

func (p *errorProxy) ErrorCode() ErrorCode {
	return p.errorCode
}

func (p *errorProxy) FullAddress() string {
	return fullAddress(p)
}

func typeError(p frame, expected Type, actual interface{}) *errorProxy {
	return &errorProxy{
		errorCode: ErrorCodeType,
		parent:    p,
		expected:  expected,
		actual:    detectType(actual),
	}
}

func elementTypeError(p frame, index int, expected Type, actual interface{}) *errorProxy {
	q := &simpleFrame{
		parent: p,
		label:  "[" + strconv.Itoa(index) + "]",
	}
	return typeError(q, expected, actual)
}

func notfoundError(p frame, address string) *errorProxy {
	return &errorProxy{
		errorCode: ErrorCodeNotFound,
		parent:    p,
		label:     address,
	}
}

// IsError checks whether p is an error and its error code is code.
func IsError(p Proxy, code ErrorCode) bool {
	err, ok := p.(Error)
	return ok && err.ErrorCode() == code
}

// IsErrorCode checks whether the error code is code.
func IsErrorCode(err error, code ErrorCode) bool {
	var myErr Error
	if errors.As(err, &myErr) {
		return myErr.ErrorCode() == code
	}
	return false
}
