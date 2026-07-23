package app

type AppConfigDTO struct {
	InheritClient     bool   `json:"inherit_client"`
	InheritClientCode string `json:"inherit_client_code"`
	InheritClientName string `json:"inherit_client_name"`
}
