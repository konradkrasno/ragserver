package models

type Document struct {
	Text string `json:"text"`
}

type AddDocumentsRequest struct {
	Documents []Document `json:"documents"`
}

type QueryRequest struct {
	SessionId string `json:"sessionId"`
	Content   string `json:"content"`
}

type QueryResponse struct {
	Query  string `json:"query"`
	Answer string `json:"answer"`
}
