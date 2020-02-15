package ikasbox

import (
	"github.com/ikascrew/ikasbox/handler"
)

func Start() error {
	return handler.Listen()
}
