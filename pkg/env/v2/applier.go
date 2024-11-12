package env

import (
	"errors"
	"fmt"
	"net/url"
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

func (a *Applier) URL(key string) *url.URL {
	v, ok := a.lookup(key)
	if !ok {
		a.err = errors.Join(a.err, errMissing(key))
		return nil
	}
	u, err := url.Parse(v)
	if err != nil {
		a.err = errors.Join(a.err, err)
		return nil
	}
	return u
}

func (a *Applier) Err() error {
	return a.err
}
