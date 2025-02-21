package consumer

import "context"

type Consumer interface {
	Start(context.Context) error
}
