package users

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,min=10,max=14"`
	Password string `json:"password" binding:"required,min=7,max=20"`
	Code     string `json:"code" binding:"max=10"`
}

type LoginByCodeRequest struct {
	Phone string `json:"phone" binding:"required,min=10,max=14"`
	Code  string `json:"code" binding:"required,min=1,max=10"`
}
