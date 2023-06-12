package atlassian

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/atlassian", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		*Atlassian
	}

	Atlassian struct {
		vu      modules.VU
		exports *goja.Object
	}

	Option func(*Atlassian) error
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

func New() *RootModule {
	return &RootModule{}
}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	runtime := vu.Runtime()
	moduleInstance := &ModuleInstance{
		Atlassian: &Atlassian{
			vu:      vu,
			exports: runtime.NewObject(),
		},
	}

	mustExport := func(name string, value interface{}) {
		if err := moduleInstance.exports.Set(name, value); err != nil {
			common.Throw(runtime, err)
		}
	}

	// Export the constructors and functions from the Atlassian module to the JS code.
	// The Writer is a constructor and must be called with new, e.g. new Jira(...).
	mustExport("Jira", moduleInstance.jiraClass)
	// The Confluence is a constructor and must be called with new, e.g. new Confluence(...).
	mustExport("Confluence", moduleInstance.confluenceClass)

	return moduleInstance
}

func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.Atlassian.exports,
	}
}

func isEmpty(object interface{}) bool {
	// check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	}

	return false
}
