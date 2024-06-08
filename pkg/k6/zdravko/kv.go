package zdravko

import "github.com/dop251/goja"

type Kv struct {
}

func (z *Zdravko) Kv() goja.Value {
	zdravkoContext := GetZdravkoContext(z.vu.Context())
	return z.vu.Runtime().ToValue(zdravkoContext.Target)
}
