package types

type SignUpBody struct {
	Name  string `json:"name" example:"8978456532"`
	Phone string `json:"phone" example:"1234"`
	Email string `json:"email" example:"98453835574"`
	DOB   string `json:"dob" example:"1234"`
}

type SignUpResponse struct {
	OTP string `json:"otp" example:"0123"`
}
type VerifyResponse struct {
	OTP string `json:"otp" example:"1123"`
}

type VerifyBody struct {
	OTP string `json:"otp" example:"0123"`
}
type SetPasswordBody struct {
	Password string `json:"password" example:"1234@Abc"`
}

type SetUserNameBody struct {
	UserName string `json:"handle" example:"sansan"`
}

type LoginBody struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type FollowBody struct {
	ID string `json:"Id" example:"123456789"`
}
