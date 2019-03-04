package model


type Error struct {
	Code string  `json:"code"`
	Message string  `json:"message"`
	Failed bool   `json:"failed"`
}
