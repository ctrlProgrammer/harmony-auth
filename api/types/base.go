package types

type Logged struct {
	FromUser    string `json:"fromUser"`
	SessionCode string `json:"sessionCode"`
}
