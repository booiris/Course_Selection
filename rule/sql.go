package rule

const (
	SQL_USER         = "root"
	SQL_PASSWORD     = "12345678"
	SQL_IPANDPORT    = "127.0.0.1:3306"
	SQL_DATABASENAME = "isuse"
)

const (
	SQL_DRIVER = "mysql"
	SQL_PATH   = SQL_USER + ":" + SQL_PASSWORD + "@tcp(" + SQL_IPANDPORT + ")/" + SQL_DATABASENAME + "?charset=utf8"
)
