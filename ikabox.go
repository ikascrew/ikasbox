package ikabox

import (
	"github.com/ikascrew/ikabox/handler"
)

func Start() error {
	return handler.Listen()
}
