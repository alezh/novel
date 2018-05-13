package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	User     string    `xorm:"not null VARCHAR(32)"`
	Created  time.Time `xorm:"DATE"`
	//Sex      string    `xorm:"CHAR(1)"`
	//Address  string    `xorm:"VARCHAR(256)"`
}

type UserMgo struct {
	Id        bson.ObjectId `bson:"_id"`
	User      string
	Pass      string
	Admin     int
	Created   bson.MongoTimestamp //创建时间戳
}

func Verify(username, password string)  {

}