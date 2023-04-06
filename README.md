# Init Utils

[![Go Reference](https://pkg.go.dev/badge/github.com/shadiestgoat/initutils.svg)](https://pkg.go.dev/github.com/shadiestgoat/initutils)

This is a package to help you with the initialization of several sub-packages.

The use case of this is if you have an application with a lot of sub-package that need to be initialized, but they require config variables, or they rely on variables that other sub-packages create. This package will help you manage those dependencies, and place an order for them.

## Features

- ✅ Context Type Safety
- ✅ Module Dependencies
- ✅ Module Pre-hooks (or *before* dependencies)
- ✅ Error Checks
- ✅ Individual Initializers, allowing for priority initializers

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

var initMgr = initutils.NewInitializer(ctx)

func Register(m initutils.Module, h func(c *InitContext), preHooks []initutils.Module, dependencies ...initutils.Module) {
    initMgr.Register(m, h, preHooks, dependencies...)
}

func Init() {
    err = initMgr.Init()
    if err != nil {
        panic(err)
    }
}
```

Then, you need to register each module using `initializer.Register`. Finally, when you wish to initialize your modules, run `initializer.Init`. This will place all your packages in the correct order (based on dependencies).

Using this approach you can also create a priority initializer, which runs before the normal one. Simply make a new initializer (you can share a context if you want to), then add another `Register()` function (like `RegisterPriority(...)`). After that, in the aliased `Init()` function, add the priority initializer to `Init()` before the normal initializer. 

This sub-package is needed so that you can define a type safe context, and define constants for your module names.
