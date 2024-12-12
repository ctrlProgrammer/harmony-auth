package types

const ()

type Session struct {
	User        User   `json:"user"`
	SessionCode string `json:"sessionCode"`
	Date        int64  `json:"date"`
}

// SECTION - API

type APILogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIValidateSession struct {
	Logged
}
