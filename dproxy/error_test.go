package dproxy

import "testing"

func TestValueProxy_OptionalBool(t *testing.T) {
	in := New(map[string]interface{}{})
	v, err := in.M("notexists").OptionalBool()
	if err != nil {
		t.Error(err)
	}
	if v != nil {
		t.Errorf("want nil, got %v", v)
	}
}

func TestValueProxy_OptionalInt64(t *testing.T) {
	in := New(map[string]interface{}{})
	v, err := in.M("notexists").OptionalInt64()
	if err != nil {
		t.Error(err)
	}
	if v != nil {
		t.Errorf("want nil, got %v", v)
	}
}

func TestValueProxy_OptionalUint64(t *testing.T) {
	in := New(map[string]interface{}{})
	v, err := in.M("notexists").OptionalUint64()
	if err != nil {
		t.Error(err)
	}
	if v != nil {
		t.Errorf("want nil, got %v", v)
	}
}

func TestValueProxy_OptionalFloat64(t *testing.T) {
	in := New(map[string]interface{}{})
	v, err := in.M("notexists").OptionalFloat64()
	if err != nil {
		t.Error(err)
	}
	if v != nil {
		t.Errorf("want nil, got %v", v)
	}
}

func TestValueProxy_OptionalString(t *testing.T) {
	in := New(map[string]interface{}{})
	v, err := in.M("notexists").OptionalString()
	if err != nil {
		t.Error(err)
	}
	if v != nil {
		t.Errorf("want nil, got %v", v)
	}
}
