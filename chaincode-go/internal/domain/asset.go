package domain

//  ABAC-managed asset. DDD entity
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type stubContextKey struct{}

var StubKey = stubContextKey{}
