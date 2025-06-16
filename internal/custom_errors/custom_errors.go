package customerrors

import "errors"

var (
	ErrNoServersAvailable        = errors.New("no servers available")
	ErrNoHealthyServersAvailable = errors.New("no servers are healthy")
)
