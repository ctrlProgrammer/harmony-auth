package types

const ()

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// SECTION - API

type APIAddUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type APIUpdateUserRole struct {
	Logged
	Email string `json:"email"`
	Role  string `json:"role"`
}
