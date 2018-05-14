package user

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	Id        int64       `xorm:"not null pk autoincr BIGINT(11)"`
	User      string    `xorm:"not null VARCHAR(32)"`
	Pass      string   `xorm:"not null VARCHAR(225)"`
	Admin     int      `xorm:"INT(1)"`
	Created   time.Time `xorm:"DATE"`
	//Sex      string    `xorm:"CHAR(1)"`
	//Address  string    `xorm:"VARCHAR(256)"`
}

var Table = "User"

type UserMgo struct {
	Id        bson.ObjectId     `bson:"_id" json:"_id"`
	User      string
	Pass      string
	Admin     int
	Created   bson.MongoTimestamp //创建时间戳
}

func (u *UserMgo)GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func (u *UserMgo)ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}