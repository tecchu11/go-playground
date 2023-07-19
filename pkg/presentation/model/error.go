package model

var statusMap = map[int]string{
	400: "Bad Request",
	401: "UnAuthorized",
	403: "Forbidden",
	404: "Resource Not Found",
	500: "Internal Server Error",
}

type ProblemDetail struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	Instant string `json:"instant"`
}

func NewProblemDetail(
	message string,
	path string,
	status int,
) *ProblemDetail {
	t, found := statusMap[status]
	if !found {
		t = string(rune(status))
	}
	return &ProblemDetail{
		Type:    "",
		Title:   t,
		Detail:  message,
		Instant: path,
	}
}
