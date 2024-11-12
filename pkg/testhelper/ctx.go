package testhelper

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SpyContext struct {
	context.Context
	mock.Mock
}

func (mck *SpyContext) Value(key any) (value any) {
	args := mck.Called(key)
	return args.Get(0)
}

var _ context.Context = (*SpyContext)(nil)
