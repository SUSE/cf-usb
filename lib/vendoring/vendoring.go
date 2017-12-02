// +build vendoring

package vendoring

// This package is for vendoring only; it is never built
// It exists so we can make godep vendor commands (instead of packages) and not
// worry about trying to actually build them
import (
	_ "github.com/jteeuwen/go-bindata"
	_ "github.com/jteeuwen/go-bindata/go-bindata"
)
