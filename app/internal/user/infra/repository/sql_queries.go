package repository

const (
	queryGetUserById = `
	select * from client where login = $1;
`
	queryCreateUser = `
	insert into client(name, email, login) values ($1, $2, $3) returning id;
`
)
