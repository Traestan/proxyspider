package models

// Proxy структура для прокси одной
type Proxy struct {
	IP     string `json:"ip"`
	Source string `json:"source"`
	Type   string `json:"type"`
}
