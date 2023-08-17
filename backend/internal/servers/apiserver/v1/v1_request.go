package v1

// ListQueryParameters is a generic struct that stores the query parameters for all GET list endpoints.
type ListQueryParameters struct {
	L *int64 `schema:"limit"`
	O *int64 `schema:"offset"`
}

func (p ListQueryParameters) Limit() *int64 {
	return p.L
}

func (p ListQueryParameters) Offset() *int64 {
	return p.O
}
