package model

type SecretPostInfo struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	EnvID       int               `json:"envId"`
	Type        string            `json:"type"`
	Value       map[string]string `json:"value"`
}
