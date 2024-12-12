package types

const ()

type Role struct {
	Name   string          `json:"name"`
	Config map[string]bool `json:"config"`
}

// SECTION - API

type APIAddRoleRequest struct {
	Logged
	Name string `json:"name"`
}

type APIConfigRoleRequest struct {
	Logged
	Id     string          `json:"id"`
	Config map[string]bool `json:"config"`
}
