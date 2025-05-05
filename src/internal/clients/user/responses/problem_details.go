package reponses

type ProblemDetails struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int16  `json:"status"`
}
