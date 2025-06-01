package response

type Response struct {
	Data   interface{} `json:"response_data"`
	Code   string      `json:"response_code"`
	Refnum string      `json:"response_refnum"`
	ID     string      `json:"response_id"`
	Desc   string      `json:"response_desc"`
}

func Init(id string) *Response {
	return &Response{
		Data: new(struct{}),
		ID:   id,
		Code: "XX",
		Desc: "General Error",
	}
}
