package response

type LoginResponseBody struct {
	Token string `json:"token"`
}

func NewLoginResponseBody(token string) *LoginResponseBody {
	return &LoginResponseBody{
		Token: token,
	}
}
