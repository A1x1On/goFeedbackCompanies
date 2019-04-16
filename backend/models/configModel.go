package models

type ConfigModel struct {
	Proxy struct {
		Adress string `json:"adress"`
		Login  string `json:"login"`
		Pass   string `json:"pass"`
	} `json:"proxy"`
	API struct {
		BaseURL struct {
			Golang string `json:"golang"`
		} `json:"baseUrl"`
	} `json:"api"`
}