package repos

import (
	"context"
	stdhttp "net/http"
)

// WebsocketService ...
type WebsocketService interface {
	Ws(context.Context, stdhttp.ResponseWriter, *stdhttp.Request) error
}
