package storage

import (
	"github.com/go-xorm/xorm"
	_"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"strconv"
	"github.com/alezh/novel/config"
	"github.com/globalsign/mgo"
	"fmt"
)

var Source *DataSource

type (

	MySqlDb interface {
		Mysql()(*xorm.Engine)
	}
	MgoDB interface {
		MongoDb() *MongoDb
	}
	DbConfig struct {
		DbIp    string
		DbName  string                  //db name
		DbPort  int                     //3306
		DbUser  string                  //root
		DbPass  string                  //password
		Charset string                  //utf8
		Prefix  string                  //prefix_
		ShowSQL bool                    //true则会在控制台打印出生成的SQL语句；
	}
	MgoConfig struct {
		MgoIP string
		MgoDB string
		MgoUser string
		MgoPass string
		MgoPort int
		Collection string
		PoolLimit int
	}
	MongoDb struct {
		Session *mgo.Session
		Database *mgo.Database
		Collection *mgo.Collection
	}
	DataSource struct {
		Mysql   *xorm.Engine
		MongoDb *MongoDb
	}
)


func init()  {
	Source = &DataSource{
		NewMysqlSource().Mysql(),
		NewMongoDBSource().MongoDb(),
	}
}

func NewMysqlSource() MySqlDb {
	return &DbConfig{
		config.MYSQL_IP,
		config.MYSQL_DB,
		config.MYSQL_PORT,
		config.MYSQL_USER,
		config.MYSQL_PASS,
		config.CHARSET,
		config.PREFIX,
		config.SHOWSQL,
	}
}

func NewMongoDBSource() MgoDB {
	return &MgoConfig{
		config.MGO_URL,
		config.MGO_DB,
		config.MGO_USER,
		config.MGO_PASS,
		config.MGO_PORT,
		config.COLLECTION,
		config.MGO_POOL,
	}
}

func (d *DbConfig)Mysql() (*xorm.Engine){

	if d.DbIp == ""{
		return nil
	}

	dataSourceName := d.DbUser+":"+d.DbPass+"@tcp("+d.DbIp+":"+ strconv.Itoa(d.DbPort) +")/"+d.DbName+"?charset="+d.Charset

	engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err!=nil{
		return nil
		//panic("orm failed to initialized")
		//return nil,errors.New("orm failed to initialized")
	}
	if errs := engine.Ping(); errs != nil{
		return nil
		//panic(errs.Error())
		//return nil,errors.New("orm failed to initialized")
	}
	if d.Prefix == ""{
		engine.SetTableMapper(core.SnakeMapper{})
	}else{
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, d.Prefix)
		engine.SetTableMapper(tbMapper)
	}
	//日志打印SQL
	engine.ShowSQL(d.ShowSQL)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(1024)
	//设置最大打开连接数
	engine.SetMaxOpenConns(2048)

	return engine
}

func (m * MgoConfig)MongoDb() *MongoDb {
	if m.MgoIP != ""{
		//fmt.Printf("mongodb://%d:$d@%s:%d/%d?minPoolSize=%s&maxIdleTimeMS=%s", m.MgoIP, m.MgoPort)
		//fmt.Printf("mongodb://%d:$d@%s:%d,%s:%d/%d?minPoolSize=%s&maxIdleTimeMS=%s", m.MgoIP, m.MgoPort)
		//connection := "mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb?minPoolSize=0&maxIdleTimeMS=3000"
		session, err := mgo.Dial(fmt.Sprintf("%s:%d", m.MgoIP, m.MgoPort))
		//p, err := mgop.DialStrongPool(connection, 5)
		if err != nil {
			return nil
			//panic(err.Error())
		}
		session.SetPoolLimit(5000)
		//session := p.AcquireSession()
		database := session.DB(m.MgoDB)
		collection := database.C(m.Collection)
		return &MongoDb{session,database,collection}
	}
	return &MongoDb{nil,nil,nil}
}