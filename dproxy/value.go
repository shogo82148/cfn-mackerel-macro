package dproxy

import (
	"math"
	"strconv"
)

type valueProxy struct {
	value  interface{}
	parent frame
	label  string
}

// valueProxy implements Proxy.
var _ Proxy = (*valueProxy)(nil)

func (p *valueProxy) Nil() bool {
	return p.value == nil
}

func (p *valueProxy) Value() (interface{}, error) {
	return p.value, nil
}

func (p *valueProxy) Bool() (bool, error) {
	switch v := p.value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case string:
		w, err := strconv.ParseBool(v)
		if err != nil {
			return false, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	default:
		return false, typeError(p, Tbool, v)
	}
}

func (p *valueProxy) OptionalBool() (*bool, error) {
	v, err := p.Bool()
	if err, ok := err.(Error); ok && err.ErrorType() == Enotfound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

type int64er interface {
	Int64() (int64, error)
}

func (p *valueProxy) Int64() (int64, error) {
	switch v := p.value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return int64(v), nil
	case int64er:
		w, err := v.Int64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	case string:
		w, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	default:
		return 0, typeError(p, Tint64, v)
	}
}

func (p *valueProxy) OptionalInt64() (*int64, error) {
	v, err := p.Int64()
	if err, ok := err.(Error); ok && err.ErrorType() == Enotfound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

type uint64er interface {
	Uint64() (uint64, error)
}

func (p *valueProxy) Uint64() (uint64, error) {
	switch v := p.value.(type) {
	case int:
		if v < 0 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case float32:
		if v < 0 || v >= 1<<64 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case float64:
		if v < 0 || v >= 1<<64 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	case uint64er:
		w, err := v.Uint64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	case string:
		w, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	default:
		return 0, typeError(p, Tint64, v)
	}
}

func (p *valueProxy) OptionalUint64() (*uint64, error) {
	v, err := p.Uint64()
	if err, ok := err.(Error); ok && err.ErrorType() == Enotfound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

type float64er interface {
	Float64() (float64, error)
}

func (p *valueProxy) Float64() (float64, error) {
	switch v := p.value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case float64er:
		w, err := v.Float64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   "overflow",
			}
		}
		return float64(v), nil
	case int64er:
		w, err := v.Int64()
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return float64(w), nil
	case string:
		w, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, &errorProxy{
				errorType: EconvertFailure,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return w, nil
	default:
		return 0, typeError(p, Tfloat64, v)
	}
}

func (p *valueProxy) OptionalFloat64() (*float64, error) {
	v, err := p.Float64()
	if err, ok := err.(Error); ok && err.ErrorType() == Enotfound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (p *valueProxy) String() (string, error) {
	switch v := p.value.(type) {
	case string:
		return v, nil
	default:
		return "", typeError(p, Tstring, v)
	}
}

func (p *valueProxy) OptionalString() (*string, error) {
	v, err := p.String()
	if err, ok := err.(Error); ok && err.ErrorType() == Enotfound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (p *valueProxy) Array() ([]interface{}, error) {
	switch v := p.value.(type) {
	case []interface{}:
		return v, nil
	default:
		return nil, typeError(p, Tarray, v)
	}
}

func (p *valueProxy) Map() (map[string]interface{}, error) {
	switch v := p.value.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		return nil, typeError(p, Tmap, v)
	}
}

func (p *valueProxy) A(n int) Proxy {
	switch v := p.value.(type) {
	case []interface{}:
		a := "[" + strconv.Itoa(n) + "]"
		if n < 0 || n >= len(v) {
			return notfoundError(p, a)
		}
		return &valueProxy{
			value:  v[n],
			parent: p,
			label:  a,
		}
	default:
		return typeError(p, Tarray, v)
	}
}

func (p *valueProxy) M(k string) Proxy {
	switch v := p.value.(type) {
	case map[string]interface{}:
		a := "." + k
		w, ok := v[k]
		if !ok {
			return notfoundError(p, a)
		}
		return &valueProxy{
			value:  w,
			parent: p,
			label:  a,
		}
	default:
		return typeError(p, Tmap, v)
	}
}

func (p *valueProxy) P(q string) Proxy {
	return pointer(p, q)
}

func (p *valueProxy) ProxySet() ProxySet {
	switch v := p.value.(type) {
	case []interface{}:
		return &setProxy{
			values: v,
			parent: p,
		}
	default:
		return typeError(p, Tarray, v)
	}
}

func (p *valueProxy) Q(k string) ProxySet {
	w := findAll(p.value, k)
	return &setProxy{
		values: w,
		parent: p,
		label:  ".." + k,
	}
}

func (p *valueProxy) findJPT(t string) Proxy {
	switch v := p.value.(type) {
	case map[string]interface{}:
		return p.M(t)
	case []interface{}:
		n, err := strconv.ParseUint(t, 10, 0)
		if err != nil {
			return &errorProxy{
				errorType: EinvalidIndex,
				parent:    p,
				infoStr:   err.Error(),
			}
		}
		return p.A(int(n))
	default:
		return &errorProxy{
			errorType: EmapNorArray,
			parent:    p,
			actual:    detectType(v),
		}
	}
}

func (p *valueProxy) parentFrame() frame {
	return p.parent
}

func (p *valueProxy) frameLabel() string {
	return p.label
}

// Default return the v as default value, if p is not found error.
func Default(p Proxy, v interface{}) Proxy {
	if err, ok := p.(Error); ok && err.ErrorType() == Enotfound {
		return New(v)
	}
	return p
}
