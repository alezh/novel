package config

//mysql
const (
	MYSQL_IP    string = "127.0.0.1"
	MYSQL_DB    string = "novel"
	MYSQL_USER  string = "root"
	MYSQL_PASS  string = "123456"
	MYSQL_PORT  int    = 3306
	CHARSET     string = "utf8mb4"
	PREFIX      string = "go_"
	SHOWSQL     bool   = false
)

//mongodb
const (
	MGO_URL     string = ""
	MGO_DB      string = ""
	MGO_USER    string = ""
	MGO_PASS    string = ""
	MGO_PORT    int    = 27017
	COLLECTION  string = ""
	MGO_POOL    int    = 200
)
