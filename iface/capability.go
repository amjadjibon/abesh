package iface

type ICapability interface {
	Name() string
	Version() string

	Source() string
	Runtime() string
	Category() string
	ContractId() string
	Values() map[string]string

	Setup() error
	SetValues(map[string]string) error

	New() ICapability
}
