package model

type GetUserByLoginReq struct {
	Login string `query:"login"`
}

type GetUserByLoginRes struct {
	Id    int64  `json:"Id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type CreateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type CreateUserRes struct {
	Id    int64  `json:"Id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}
