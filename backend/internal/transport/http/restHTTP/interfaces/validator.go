package interfaces

import "context"

type Validator interface {
	StructCtx(ctx context.Context, s interface{}) error
}
