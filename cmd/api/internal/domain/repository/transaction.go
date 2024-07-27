package repository

import (
	"context"
)

type TransactionRepository interface {
	Do(context.Context, func(context.Context) error) error
}
