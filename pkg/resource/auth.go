package resource

type Auth struct {
	Code    *int    `json:"code,omitempty"`
	Data    *Data   `json:"data,omitempty"`
	Message *string `json:"msg,omitempty"`
}

type Data struct {
	Token *string `json:"token,omitempty"`
}
