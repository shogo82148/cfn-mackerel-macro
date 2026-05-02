package dproxy

import "strings"

var jptR = strings.NewReplacer("~1", "/", "~0", "~")

func unescapeJPT(t string) string {
	return jptR.Replace(t)
}

func pointer(p Proxy, q string) Proxy {
	if len(q) == 0 {
		return p
	}
	if q[0] != '/' {
		return &errorProxy{
			errorCode: ErrorCodeInvalidQuery,
			parent:    p,
			infoStr:   "not start with '/'",
		}
	}
	for t := range strings.SplitSeq(q[1:], "/") {
		p = p.findJPT(unescapeJPT(t))
	}
	return p
}

// Pointer returns a Proxy which pointed by JSON Pointer's query q
func Pointer(v any, q string) Proxy {
	return pointer(New(v), q)
}
