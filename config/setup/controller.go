package setup

import (
	"github.com/osfx/snail/config/parse"
	"github.com/osfx/snail/server"
)

type Controller struct {
	*server.Config
	parse.Dispenser
}
