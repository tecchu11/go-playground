package env

func (a *Applier) Lookup() func(string) (string, bool) {
	return a.lookup
}

var ErrMissing = errMissing
