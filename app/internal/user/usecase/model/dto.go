package model

type GetUserByLoginReq struct {
	Login string
}

type CreateUserReq struct {
	Name  string
	Email string
	Login string
}

type CreateUserRes struct {
	Id int64
}
