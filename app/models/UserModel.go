package model

type UserRegisterRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10,max=16"`
	Name        string `json:"name" binding:"required,min=5,max=50"`
	Password    string `json:"password" binding:"required,min=5,max=15"`
}

type UserRegisterResponse struct {
	UserId      string `json:"userID"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10,max=16"`
	Password    string `json:"password" binding:"required,min=5,max=15"`
}

type GetCustomerParams struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}
