package response

type Meta struct {
	Success    bool   `json:"success"`
	Message    string `json:"message" default:"true"`
	StatusCode int    `json:"status_code"`
	Detail     any    `json:"detail,omitempty"`
}
