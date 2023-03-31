package initutils

import (
	"fmt"
	"sort"
)

// Please note, T is the init context. T must be a struct, and **not** a pointer.
type Initializer[T any] struct {
	deps map[Module][]Module
	handlers  map[Module]func(c *T)
	ctx *T
}

// Create a new Initializer instance. This will create a new dependency manager, any dependencies here will not be shared with other initializers.
// Please note, T is the init context. T must be a struct, and **not** a pointer.
func NewInitializer[T any](ctx *T) *Initializer[T] {
	if ctx == nil {
		ctx = new(T)
	}

	return &Initializer[T]{
		deps: map[Module][]Module{},
		handlers:  map[Module]func(c *T){},
		ctx: ctx,
	}
}

func (i *Initializer[T]) Register(m Module, h func(c *T), dependencies ...Module) {
	i.handlers[m] = h
	i.deps[m] = dependencies
}

// Call all the initialization functions, in the correct order (based on the dependencies each module requires)
// This method returns ErrDepCycle. If this happens, your application itself is having dependency cycles.
// If such a thing happens, you should panic.
func (i *Initializer[T]) Init() error {
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
				return &ErrUnknownDep{
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
				return ErrDepCycle{
					Module1: m,
					Module2: r,
				}
			}
		}
	}

	// Less reports whether the element with index i
	// must sort before the element with index j.
	sort.SliceStable(modules, func(i, j int) bool {
		return !depMaps[modules[i]][modules[j]]
	})

	for _, m := range modules {
		i.handlers[m](i.ctx)
	}

	return nil
}

type Module string

type ErrDepCycle struct {
	Module1 Module
	Module2 Module
}

func (e ErrDepCycle) Error() string {
	return fmt.Sprintf(`Dependency cycle between '%s' and '%s'`, e.Module1, e.Module2)
}

type ErrUnknownDep struct {
	Module Module
	Dep Module
}
func (e ErrUnknownDep) Error() string {
	return fmt.Sprintf("Module '%s' requires module '%s' but module '%s' was never registered", e.Module, e.Dep, e.Module)
}

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
