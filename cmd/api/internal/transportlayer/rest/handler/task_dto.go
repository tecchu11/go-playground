package handler

type (
	ReqPostTask struct {
		Content string `json:"content"`
	}
	ResPostTask struct {
		ID string `json:"id"`
	}
)

type ReqPutTask struct {
	Content string `json:"content"`
}
