package auth

type createRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createResponse struct {
	Successs bool   `json:"success"`
	Username string `json:"username"`
}


type logoutResponse struct {
	Success bool `json:"success"`
}