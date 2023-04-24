package initutils

import (
	"sort"
)

// Please note, T is the init context. T must be a struct, and **not** a pointer.
type Initializer[T any] struct {
	deps     map[Module][]Module
	handlers map[Module]func(c *T)
	ctx      *T
}

// Create a new Initializer instance. This will create a new dependency manager, any dependencies here will not be shared with other initializers.
// Please note, T is the init context. T must be a struct, and **not** a pointer.
func NewInitializer[T any](ctx *T) *Initializer[T] {
	if ctx == nil {
		ctx = new(T)
	}

	return &Initializer[T]{
		deps:     map[Module][]Module{},
		handlers: map[Module]func(c *T){},
		ctx:      ctx,
	}
}

// Register a module, along with it's dependencies
// You can also register what this module must be executed before, with a pre-hook, using the preHook argument.
func (i *Initializer[T]) Register(m Module, h func(c *T), preHooks []Module, dependencies ...Module) {
	i.handlers[m] = h

	if i.deps[m] == nil {
		i.deps[m] = []Module{}
	}

	i.deps[m] = append(i.deps[m], dependencies...)

	for _, bfr := range preHooks {
		if i.deps[bfr] == nil {
			i.deps[bfr] = []Module{}
		}

		i.deps[bfr] = append(i.deps[bfr], m)
	}
}

// Outputs the order in which the modules will be loaded in.
// Error this can output are *ErrUnknownDep, in case there is a dependency that is never registered and *ErrDepCycle, in case there is a dependency cycle.
func (i *Initializer[T]) Plan() ([]Module, error) {
	newDeps := map[Module][]Module{}

	for m := range i.deps {
		i.resolve(m, newDeps)
	}

	depMaps := map[Module]map[Module]bool{}
	modules := []Module{}

	for m, reqs := range newDeps {
		depMaps[m] = map[Module]bool{}
		modules = append(modules, m)

		for _, dep := range reqs {
			if _, ok := i.handlers[dep]; !ok {
				return nil, &ErrUnknownDep{
					Module: m,
					Dep:    dep,
				}
			}
			depMaps[m][dep] = true
		}
	}

	for m, deps := range newDeps {
		for _, r := range deps {
			if depMaps[r][m] {
				return nil, ErrDepCycle{
					Module1: m,
					Module2: r,
				}
			}
		}
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i] > modules[j]
	})
	
	sort.SliceStable(modules, func(i, j int) bool {
		return !depMaps[modules[i]][modules[j]]
	})

	return modules, nil
}

// Call all the initialization functions, in the correct order (based on the dependencies each module requires)
// This function outputs the same errors that Plan() can output. 
// If there is an error, the application should panic, as this error is most likely baked in.
func (i *Initializer[T]) Init() error {
	modules, err := i.Plan()

	if err != nil {
		return err
	}
	
	for _, m := range modules {
		i.handlers[m](i.ctx)
	}

	return nil
}

// The name of a module. This should be used for constants in your initializer sub package.
type Module string

func (i *Initializer[T]) resolve(m Module, cache map[Module][]Module) []Module {
	if reqs, ok := cache[m]; ok {
		return reqs
	}

	reqs := []Module{}

	if len(i.deps[m]) == 0 {
		cache[m] = []Module{}

		return []Module{}
	}

	for _, req := range i.deps[m] {
		reqs = append(reqs, i.resolve(req, cache)...)
		reqs = append(reqs, req)
	}

	dupeMap := map[Module]bool{}
	deDupedReqs := []Module{}

	for _, r := range reqs {
		if dupeMap[r] {
			continue
		}

		dupeMap[r] = true
		deDupedReqs = append(deDupedReqs, r)
	}

	cache[m] = deDupedReqs

	return reqs
}
