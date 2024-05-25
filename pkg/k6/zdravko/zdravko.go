package zdravko

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/zdravko", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
		// comparator is the exported type
		zdravko *Zdravko
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:      vu,
		zdravko: &Zdravko{vu: vu},
	}
}

type Target struct {
	Name     string
	Group    string
	Metadata map[string]interface{}
}

type Zdravko struct {
	vu      modules.VU
	Targets []Target
}

func (z *Zdravko) GetTarget() goja.Value {
	zdravkoContext := GetZdravkoContext(z.vu.Context())
	return z.vu.Runtime().ToValue(zdravkoContext.Target)
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.zdravko,
	}
}
