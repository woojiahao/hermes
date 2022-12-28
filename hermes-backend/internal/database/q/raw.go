package q

type RawQuery struct {
	query string
}

func Raw(query string) *RawQuery {
	return &RawQuery{query}
}

func (r *RawQuery) Generate() string {
	return r.query
}
