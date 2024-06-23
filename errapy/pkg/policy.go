package pkg

type PolicyConfig struct {
	ClassesRequired bool
	CodesRequired   bool
}

type PolicyOption interface {
	apply(*PolicyConfig)
}

type optionFunc func(*PolicyConfig)

func (f optionFunc) apply(c *PolicyConfig) {
	f(c)
}

func WithClassesRequired(classesRequired bool) PolicyOption {
	return optionFunc(func(c *PolicyConfig) {
		c.ClassesRequired = classesRequired
	})
}
func WithCodesRequired(codesRequired bool) PolicyOption {
	return optionFunc(func(c *PolicyConfig) {
		c.CodesRequired = codesRequired
	})
}

func Policy(opts ...PolicyOption) *PolicyConfig {
	policy := &PolicyConfig{
		ClassesRequired: false,
		CodesRequired:   false,
	}
	for _, opt := range opts {
		opt.apply(policy)
	}
	return policy
}
