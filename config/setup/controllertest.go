package setup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/osfx/snail/config/parse"
	"github.com/osfx/snail/middleware"
	"github.com/osfx/snail/server"
)

// NewTestController creates a new *Controller for
// the input specified, with a filename of "Testfile"
func NewTestController(input string) *Controller {
	return &Controller{
		Config:    &server.Config{},
		Dispenser: parse.NewDispenser("Testfile", strings.NewReader(input)),
	}
}

// EmptyNext is a no-op function that can be passed into
// middleware.Middleware functions so that the assignment
// to the Next field of the Handler can be tested.
var EmptyNext = middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
})

// SameNext does a pointer comparison between next1 and next2.
func SameNext(next1, next2 middleware.Handler) bool {
	return fmt.Sprintf("%p", next1) == fmt.Sprintf("%p", next2)
}
