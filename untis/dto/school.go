package dto

type SchoolResponse struct {
	Result struct {
		Schools []struct {
			Server      string `json:"server"`
			Address     string `json:"address"`
			DisplayName string `json:"displayName"`
			LoginName   string `json:"loginName"`
		} `json:"schools"`
	} `json:"result"`
}
