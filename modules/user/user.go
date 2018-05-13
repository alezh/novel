package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	User     string    `xorm:"not null VARCHAR(32)"`
	Pass      string   `xorm:"not null VARCHAR(225)"`
	Admin     int      `xorm:"INT(1)"`
	Created  time.Time `xorm:"DATE"`
	//Sex      string    `xorm:"CHAR(1)"`
	//Address  string    `xorm:"VARCHAR(256)"`
}

var Table = "User"

type UserMgo struct {
	Id        bson.ObjectId     `bson:"_id" json:"id"`
	User      string
	Pass      string
	Admin     int
	Created   bson.MongoTimestamp //创建时间戳
}

func Verify(username, password string)  {

}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}