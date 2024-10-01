package entity

// Base is entity base interface for this app.
type Base interface {
	EncodeCursor() (string, error)
}

// Page is pagination of entities.
type Page[E Base] struct {
	Items     []E    `json:"items"`
	HasNext   bool   `json:"hasNext"`
	NextToken string `json:"next"`
}

// NewPage converts Page by E.
func NewPage[E Base](s []E, limit int32) (Page[E], error) {
	if len(s) >= int(limit)+1 {
		next, err := s[limit].EncodeCursor()
		if err != nil {
			return Page[E]{}, err
		}
		return Page[E]{
			Items:     s[:limit],
			HasNext:   true,
			NextToken: next,
		}, nil
	}
	return Page[E]{
		Items:   s,
		HasNext: false,
	}, nil
}
