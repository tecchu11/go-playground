package entity

// Item is pagination item type.
type Item[T comparable] interface {
	Token() T
}

// CursorPage is cursor paginatied items.
// Type T is cursor token type and type I is item type.
type CursorPage[T comparable, I Item[T]] struct {
	Items     []I  `json:"items"`
	HasNext   bool `json:"hasNext"`
	NextToken T    `json:"next"`
}

// NewCursorPage inits CursorPage.
// NextToken is zero if no next items.
func NewCursorPage[T comparable, I Item[T]](items []I, limit int32) CursorPage[T, I] {
	if len(items) >= int(limit)+1 {
		next := items[limit].Token()
		return CursorPage[T, I]{
			Items:     items[:limit],
			HasNext:   true,
			NextToken: next,
		}
	}
	return CursorPage[T, I]{
		Items:   items,
		HasNext: false,
	}
}
