package dproxy

import "testing"

func TestValueProxy_Bool(t *testing.T) {
	tests := []struct {
		in  interface{}
		out bool
		err ErrorType
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
			err: Etype,
		},
		{
			in:  "foobar",
			err: EconvertFailure,
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
			if myErr.ErrorType() != tt.err {
				t.Errorf("%d: unexpected error type: want %s, got %s", i, tt.err, myErr.ErrorType())
			}
		}
	}
}
