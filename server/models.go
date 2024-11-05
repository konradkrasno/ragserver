package server

type document struct {
	Text string `json:"text"`
}

type addDocumentsRequest struct {
	Documents []document `json:"documents"`
}

type queryRequest struct {
	Content string `json:"content"`
}
