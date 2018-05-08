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

//配置文件
const (
	CONFIG      string = "config.ini"       // 配置文件路径
)

//mongodb
const (
	MGO_URL     string = "127.0.0.1"
	MGO_DB      string = "BookDb"
	MGO_USER    string = ""
	MGO_PASS    string = ""
	MGO_PORT    int    = 27017
	COLLECTION  string = "BookCover"
	MGO_POOL    int    = 2048
)

// 运行状态
const (
	STOPPED = iota - 1
	STOP
	RUN
	PAUSE
)
// 运行模式
const (
	UNSET int = iota - 1
	OFFLINE
	SERVER
	CLIENT
)

// 数据头部信息
const (
	// 任务请求Header
	REQTASK = iota + 1
	// 任务响应流头Header
	TASK
	// 打印Header
	LOG
)

const SessionIDKey = "UserID"
