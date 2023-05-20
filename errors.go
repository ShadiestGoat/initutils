package initutils

import (
	"errors"
	"fmt"
)

type ErrDepCycle struct {
	Module1 Module
	Module2 Module
}

func (e ErrDepCycle) Error() string {
	return fmt.Sprintf(`Dependency cycle between '%s' and '%s'`, e.Module1, e.Module2)
}

type ErrUnknownDep struct {
	Module Module
	Dep    Module
}

func (e ErrUnknownDep) Error() string {
	return fmt.Sprintf("Module '%s' requires module '%s' but module '%s' was never registered", e.Module, e.Dep, e.Module)
}

var ErrAlreadyInitialized = errors.New("the initializer has already been called")