package dproxy

import (
	"math"
	"testing"
)

func TestValueProxy_Bool(t *testing.T) {
	tests := []struct {
		in  interface{}
		out bool
		err ErrorCode
	}{
		// boolean types
		{
			in:  true,
			out: true,
		},
		{
			in:  false,
			out: false,
		},

		// integer types
		{
			in:  int(0),
			out: false,
		},
		{
			in:  int(1),
			out: true,
		},
		{
			in:  int8(0),
			out: false,
		},
		{
			in:  int8(1),
			out: true,
		},
		{
			in:  int16(0),
			out: false,
		},
		{
			in:  int16(1),
			out: true,
		},
		{
			in:  int32(0),
			out: false,
		},
		{
			in:  int32(1),
			out: true,
		},
		{
			in:  int64(0),
			out: false,
		},
		{
			in:  int64(1),
			out: true,
		},
		{
			in:  uint(0),
			out: false,
		},
		{
			in:  uint(1),
			out: true,
		},
		{
			in:  uint8(0),
			out: false,
		},
		{
			in:  uint8(1),
			out: true,
		},
		{
			in:  uint16(0),
			out: false,
		},
		{
			in:  uint16(1),
			out: true,
		},
		{
			in:  uint32(0),
			out: false,
		},
		{
			in:  uint32(1),
			out: true,
		},
		{
			in:  uint64(0),
			out: false,
		},
		{
			in:  uint64(1),
			out: true,
		},

		// strings
		{
			in:  "false",
			out: false,
		},
		{
			in:  "true",
			out: true,
		},

		// errors
		{
			in:  struct{}{},
			err: ErrorCodeType,
		},
		{
			in:  "foobar",
			err: ErrorCodeConvertFailure,
		},
	}

	for i, tt := range tests {
		proxy := New(tt.in)
		v, err := proxy.Bool()
		if tt.err == 0 {
			if v != tt.out {
				t.Errorf("%d: want %v, got %v", i, tt.out, v)
			}
		} else {
			myErr, ok := err.(Error)
			if !ok {
				t.Errorf("%d: want dproxy.Error, but not", i)
				continue
			}
			if myErr.ErrorCode() != tt.err {
				t.Errorf("%d: unexpected error type: want %s, got %s", i, tt.err, myErr.ErrorCode())
			}
		}
	}
}

func TestValueProxy_Int64(t *testing.T) {
	tests := []struct {
		in  interface{}
		out int64
		err ErrorCode
	}{
		// integer types
		{
			in:  int(math.MaxInt32),
			out: math.MaxInt32,
		},
		{
			in:  int(math.MinInt32),
			out: math.MinInt32,
		},
		{
			in:  int8(math.MaxInt8),
			out: math.MaxInt8,
		},
		{
			in:  int8(math.MinInt8),
			out: math.MinInt8,
		},
		{
			in:  int16(math.MaxInt16),
			out: math.MaxInt16,
		},
		{
			in:  int16(math.MinInt16),
			out: math.MinInt16,
		},
		{
			in:  int32(math.MaxInt32),
			out: math.MaxInt32,
		},
		{
			in:  int32(math.MinInt32),
			out: math.MinInt32,
		},
		{
			in:  int64(math.MaxInt64),
			out: math.MaxInt64,
		},
		{
			in:  int64(math.MinInt64),
			out: math.MinInt64,
		},
		{
			in:  uint(math.MaxUint32),
			out: math.MaxUint32,
		},
		{
			in:  uint8(math.MaxUint8),
			out: math.MaxUint8,
		},
		{
			in:  uint16(math.MaxUint16),
			out: math.MaxUint16,
		},
		{
			in:  uint32(math.MaxUint32),
			out: math.MaxUint32,
		},
		{
			in:  uint64(math.MaxInt64),
			out: math.MaxInt64,
		},

		// floats
		{
			in:  float64(1 << 53),
			out: 1 << 53,
		},

		// strings
		{
			in:  "9223372036854775807", // math.MaxInt64
			out: 9223372036854775807,
		},
		{
			in:  "-9223372036854775808", // math,MinInt64
			out: -9223372036854775808,
		},

		// errors
		{
			in:  struct{}{},
			err: ErrorCodeType,
		},
		{
			in:  "foobar",
			err: ErrorCodeConvertFailure,
		},
		{
			in:  uint64(math.MaxInt64 + 1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  float32(1 << 63),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  math.Float32frombits(math.Float32bits(-(1 << 63)) + 1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  float64(1 << 63),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  math.Float64frombits(math.Float64bits(-(1 << 63)) + 1),
			err: ErrorCodeConvertFailure,
		},
	}

	for i, tt := range tests {
		proxy := New(tt.in)
		v, err := proxy.Int64()
		if tt.err == 0 {
			if v != tt.out {
				t.Errorf("%d: want %v, got %v", i, tt.out, v)
			}
		} else {
			myErr, ok := err.(Error)
			if !ok {
				t.Errorf("%d: want dproxy.Error, but not", i)
				continue
			}
			if myErr.ErrorCode() != tt.err {
				t.Errorf("%d: unexpected error type: want %s, got %s", i, tt.err, myErr.ErrorCode())
			}
		}
	}
}

func TestValueProxy_Uint64(t *testing.T) {
	tests := []struct {
		in  interface{}
		out uint64
		err ErrorCode
	}{
		// integer types
		{
			in:  int(math.MaxInt32),
			out: math.MaxInt32,
		},
		{
			in:  int8(math.MaxInt8),
			out: math.MaxInt8,
		},
		{
			in:  int16(math.MaxInt16),
			out: math.MaxInt16,
		},
		{
			in:  int32(math.MaxInt32),
			out: math.MaxInt32,
		},
		{
			in:  int64(math.MaxInt64),
			out: math.MaxInt64,
		},
		{
			in:  uint(math.MaxUint32),
			out: math.MaxUint32,
		},
		{
			in:  uint8(math.MaxUint8),
			out: math.MaxUint8,
		},
		{
			in:  uint16(math.MaxUint16),
			out: math.MaxUint16,
		},
		{
			in:  uint32(math.MaxUint32),
			out: math.MaxUint32,
		},
		{
			in:  uint64(math.MaxUint64),
			out: math.MaxUint64,
		},

		// floats
		{
			in:  float64(1 << 53),
			out: 1 << 53,
		},

		// strings
		{
			in:  "18446744073709551615", // math.MaxUint64
			out: 18446744073709551615,
		},

		// errors
		{
			in:  struct{}{},
			err: ErrorCodeType,
		},
		{
			in:  int(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  int8(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  int16(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  int32(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  int64(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  float64(-1),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  float64(1 << 64),
			err: ErrorCodeConvertFailure,
		},
		{
			in:  "foobar",
			err: ErrorCodeConvertFailure,
		},
		{
			in:  "18446744073709551616", // math.MaxUint64 + 1
			err: ErrorCodeConvertFailure,
		},
	}

	for i, tt := range tests {
		proxy := New(tt.in)
		v, err := proxy.Uint64()
		if tt.err == 0 {
			if v != tt.out {
				t.Errorf("%d: want %v, got %v", i, tt.out, v)
			}
		} else {
			myErr, ok := err.(Error)
			if !ok {
				t.Errorf("%d: want dproxy.Error, but not", i)
				continue
			}
			if myErr.ErrorCode() != tt.err {
				t.Errorf("%d: unexpected error type: want %s, got %s", i, tt.err, myErr.ErrorCode())
			}
		}
	}
}
