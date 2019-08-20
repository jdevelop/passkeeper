package firmware

type Credentials struct {
	Id      string `json:"id,omitempty"`
	Service string `json:"service,omitempty"`
	Secret  string `json:"secret,omitempty"`
	Comment string `json:"comment,omitempty"`
}
