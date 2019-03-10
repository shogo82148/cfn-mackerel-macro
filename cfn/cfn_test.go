package cfn

import (
	"testing"

	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
)

func TestProxyOptionalFloat64(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var d dproxy.Drain
		in := dproxy.New(map[string]interface{}{})
		ret := proxyOptionalFloat64(&d, in.M("foobar"))
		if ret != nil {
			t.Errorf("want nil, got %v", ret)
		}
		if err := d.CombineErrors(); err != nil {
			t.Error(err)
		}
	})
	t.Run("float64 value", func(t *testing.T) {
		var d dproxy.Drain
		in := dproxy.New(map[string]interface{}{
			"foobar": 50,
		})
		ret := proxyOptionalFloat64(&d, in.M("foobar"))
		if ret == nil {
			t.Error("want not nil, got nil")
			return
		}
		if *ret != 50 {
			t.Errorf("want 50.0, got %f", *ret)
		}
		if err := d.CombineErrors(); err != nil {
			t.Error(err)
		}
	})
	t.Run("string value", func(t *testing.T) {
		var d dproxy.Drain
		in := dproxy.New(map[string]interface{}{
			"foobar": "50",
		})
		ret := proxyOptionalFloat64(&d, in.M("foobar"))
		if ret == nil {
			t.Error("want not nil, got nil")
			return
		}
		if *ret != 50 {
			t.Errorf("want 50.0, got %f", *ret)
		}
		if err := d.CombineErrors(); err != nil {
			t.Error(err)
		}
	})
}
