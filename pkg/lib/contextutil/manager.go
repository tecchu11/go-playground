package contextutil

import "context"

// ContextManager interface for retrve value of type T from context.
//
// How to use:
//
//	 type FooCtxManager struct {}
//
//	 type fooCtxKey struct {}
//
//	 func (manager *FooCtxManager) Get(ctx context.Context) (*Foo, errror) {
//			v, ok := ctx.Value(fooCtxKey{}).(*fooCtxKey)
//	     if !ok || v == nil {
//				return nil,  fmt.Errorf("error")
//	     }
//	     return v, nil
//	 }
//
//	 func (manager *FooCtxManager) set(ctx context.Context, v Foo) context.Context {
//			ctx := context.WithValue(ctx, fooCtxKey{}, v)
//			return ctx
//	 }
type ContextManager[T any] interface {
	Get(ctx context.Context) (T, error)
}
