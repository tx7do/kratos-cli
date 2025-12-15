//go:build wireinject
// +build wireinject

package {{.Package}}

import (
	"github.com/google/wire"
)

// ProviderSet is {{.Package}} providers.
var ProviderSet = wire.NewSet(
{{renderInParameters .NewFunctions}})
