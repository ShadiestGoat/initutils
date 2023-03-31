# Init Utils

This is a package to help you with the initialization of several sub-packages.

The use case of this is if you have an application with a lot of sub-package that need to be initialized, but they require config variables, or they rely on variables that other sub-packages create. This package will help you manage those dependencies, and place an order for them.

## Usage

I recommend creating an `initializer` sub-package for your application.

```go
package initializer

import (
    "github.com/shadiestgoat/initutils"

    // ... these deps are used for context, so they will be different for your application
	"github.com/ShadiestGoat/donation-api-wrapper"
	"github.com/bwmarrin/discordgo"
)

const (
    // ... these are examples, they can be different for your application!
    MOD_CONFIG   initutils.Module = "Config"
    MOD_DISCORD  initutils.Module = "Discord"
    MOD_DONATION initutils.Module = "Donations"
)

type InitContext struct {
    // ... these are examples, they can be different for your application!
    Discord  *discordgo.Session
    Donation *donations.Client
}

var ctx = &InitContext{}

var priorityInit = initutils.NewInitializer[InitContext](ctx)
var normalInit   = initutils.NewInitializer[InitContext](ctx)

func RegisterPriority(m initutils.Module, h func(c *InitContext), dependencies ...initutils.Module) {
    priorityInit.Register(m, h, dependencies...)
}

func Register(m initutils.Module, h func(c *InitContext), dependencies ...initutils.Module) {
    normalInit.Register(m, h, dependencies...)
}

func Init() {
    err := priorityInit.Init()
    if err != nil {
        panic(err)
    }

    err = normalInit.Init()
    if err != nil {
        panic(err)
    }
}
```

Then, you need to register each module using ether `initializer.RegisterPriority` or `initializer.Register`. Finally, when you wish to initialize your modules, run `initializer.Init`. This will place all your packages in the correct order (based on dependencies)

This sub-package is needed so that you can define a type safe context, and define constants for your module names.
