package model

type CreateUserReq struct {
	Name  string
	Email string
	Login string
}

type CreateUserRes struct {
	Id int64 `db:"id"`
}

type GetUserByLoginReq struct {
	Login string
}

type GetUserByLoginRes struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Login string `db:"login"`
}
