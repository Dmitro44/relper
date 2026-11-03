package roles

import "fmt"

var ErrUnknownRole = fmt.Errorf("unknown role code")
var ErrLastAdmin = fmt.Errorf("cannot remove last admin role")
