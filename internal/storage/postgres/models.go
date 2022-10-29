package postgres

type UserInfo struct {
	Id      string
	Balance int64
}

type Order struct {
	IdUser    string
	IdService string
	IdOrder   string
	Cost      int64
	Accepted  bool
}
