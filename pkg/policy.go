package pkg

import (
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		_ = Policy(WithClassesRequired(false), WithCodesRequired(false))
	})
}

var defaultPolicy *PolicyConfig

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
	p := &PolicyConfig{
		ClassesRequired: false,
		CodesRequired:   false,
	}
	for _, opt := range opts {
		opt.apply(p)
	}
	defaultPolicy = p
	return p
}
