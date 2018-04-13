package storage

import (
	"github.com/go-xorm/xorm"
	_"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"strconv"
	"github.com/alezh/novel/config"
)

var Source *DataSource

type (

	Db interface {
		Mysql()(*xorm.Engine)
	}
	DbConfig struct {
		DbName  string                  //db name
		DbPort  int                     //3306
		DbUser  string                  //root
		DbPass  string                  //password
		Charset string                  //utf8
		Prefix  string                  //prefix_
		ShowSQL bool                    //true则会在控制台打印出生成的SQL语句；
	}
	DataSource struct {
		Mysql   *xorm.Engine
	}
)

func init()  {
	Source =  &DataSource{
		NewMysqlSource().Mysql(),
	}
}

func NewMysqlSource() Db {
	return &DbConfig{
		config.MYSQL_DB,
		config.MYSQL_PORT,
		config.MYSQL_USER,
		config.MYSQL_PASS,
		config.CHARSET,
		config.PREFIX,
		config.SHOWSQL,
	}
}

func (d *DbConfig)Mysql() (*xorm.Engine){

	dataSourceName := d.DbUser+":"+d.DbPass+"@:"+ strconv.Itoa(d.DbPort) +"/"+d.DbName+"?charset="+d.Charset

	engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err!=nil{
		panic("orm failed to initialized")
		//return nil,errors.New("orm failed to initialized")
	}
	//if errs := engine.Ping(); errs!=nil{
	//	panic("orm failed to initialized")
	//	//return nil,errors.New("orm failed to initialized")
	//}
	if d.Prefix == ""{
		engine.SetTableMapper(core.SnakeMapper{})
	}else{
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, d.Prefix)
		engine.SetTableMapper(tbMapper)
	}
	//日志打印SQL
	engine.ShowSQL(d.ShowSQL)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(20)

	return engine
}