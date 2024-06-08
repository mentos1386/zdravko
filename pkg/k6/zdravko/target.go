package zdravko

import "github.com/dop251/goja"

func (z *Zdravko) GetTarget() goja.Value {
	zdravkoContext := GetZdravkoContext(z.vu.Context())
	return z.vu.Runtime().ToValue(zdravkoContext.Target)
}
