package env

import (
	"errors"
	"fmt"
	"os"
)

type Applier struct {
	err    error
	lookup func(string) (string, bool)
}

func New(lookup func(string) (string, bool)) *Applier {
	if lookup == nil {
		lookup = os.LookupEnv
	}
	return &Applier{lookup: lookup}
}

func errMissing(key string) error {
	return fmt.Errorf("missing env %s", key)
}

func (a *Applier) String(key string) string {
	v, ok := a.lookup(key)
	if ok {
		return v
	}
	a.err = errors.Join(a.err, errMissing(key))
	return ""
}

func (a *Applier) Err() error {
	return a.err
}
