package repositories

import (
	"sync"
	"github.com/alezh/novel/storage"
	."github.com/alezh/novel/modules/user"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/alezh/novel/config"
)

type (
    //Query func(User) bool

	UserRepository interface {
		Verify(username, password string) (interface{},bool)
		SelectById(id,pSlice interface{}) bool
		Select(interface{},interface{}) bool
		SelectMany(interface{},string,interface{}) bool
		InsertOrUpdate(interface{}) bool
		Insert(interface{}) bool
		InsertUser(name,userPassword string,admin int) bool
		Update(id,pData interface{}) bool
		//Delete(query Query, limit int) (deleted bool)
	}
	userRepository struct {
		source *storage.DataSource
		mu     sync.RWMutex
		User   User
	}
	userMgoRepository struct {
		source *storage.DataSource
		mu     sync.RWMutex
		User   UserMgo
	}
)

func NewUserRepository(source *storage.DataSource) UserRepository {
	switch config.Task.OutType {
	case "mgo":
		return &userMgoRepository{source: source}
	case "mysql":
		return &userRepository{source: source}
	default:
		return &userRepository{source: source}
	}
}
func (r *userRepository)Verify(username, password string) (interface{},bool) {
	return nil,false
}
func (r *userRepository)SelectById(id,pSlice interface{}) bool  {
	return false
}
func (r *userRepository)Select(search interface{},pSlice interface{}) bool {
	return false
}
func (r *userRepository)SelectMany(search interface{},sortKey string, pSlice interface{}) bool {
	return false
}

func (r *userRepository)InsertOrUpdate(pSlice interface{}) bool {
	return false
}

func (r *userRepository)Insert(pSlice interface{}) bool{
	return false
}
func (r *userRepository)Update(id,pData interface{}) bool {
	return false
}
func (r *userRepository)InsertUser(name,userPassword string,admin int) bool{
	return false
}


//--------------------------------mgo----------------------------------------------------------------------------------

func (r *userMgoRepository)Verify(username, password string) (interface{},bool) {
	user := new(UserMgo)
	search := bson.M{"user": username}
	if r.Select(search,user){
		hash := []byte(user.Pass)
		if ok, _ := user.ValidatePassword(password,hash);ok{
			return user.Id.Hex() ,true
		}
	}
	return nil, false
}

func (r *userMgoRepository)InsertUser(name,userPassword string,admin int) bool{
	user := new(UserMgo)
	objId := bson.NewObjectId()
	hashed, _ := user.GeneratePassword(userPassword)
	user.Id   = objId
	user.User = name
	user.Admin = admin
	user.Pass = string(hashed)
	return r.Insert(user)
}

func (r *userMgoRepository)SelectById(id,pSlice interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.FindId(id).One(pSlice)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find findId: %v", id)
		} else {
			fmt.Println(err.Error())
		}
		return false
	}
	return true
}

func (r *userMgoRepository)Select(search interface{}, pSlice interface{})  bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.Find(search).One(pSlice)
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not Find findall: %v", search)
		} else {
			fmt.Println(err.Error())
		}
		return false
	}
	return true
}

func (r *userMgoRepository)SelectMany(search interface{} , sortKey string, pSlice interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	var err error
	if sortKey == ""{
		err = coll.Find(search).All(pSlice)
	}else{
		err = coll.Find(search).Sort(sortKey).All(pSlice)
	}
	if err != nil {
		if err == mgo.ErrNotFound {
			fmt.Printf("Not FindAll findall: %v", search)
		} else {
			fmt.Println(err.Error())
		}
		return false
	}
	return true
}

func (r *userMgoRepository)Update(id,pData interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.UpdateId(id, pData)
	if err != nil {
		//fmt.Printf("UpdateSync error: %v \r\ntable: %s \r\nid: %v \r\ndata: %v \r\n",
		//	err.Error(), table, id, pData)
		return false
	}
	return true
}

func (r *userMgoRepository)InsertOrUpdate(pData interface{}) bool {
	user := pData.(UserMgo)
	search := bson.M{"user": user.User}
	result := UserMgo{}
	if err := r.Select(search,result);err{
		r.Update(result.Id,result)
		return false
	}else {
		coll := r.source.MongoDb.Database.C(Table)
		err := coll.Insert(pData)
		if err != nil {
			fmt.Printf("InsertSync error: %v \r\n", err.Error())
			return false
		}
		return true
	}
}

func (r *userMgoRepository)Insert(pData interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.Insert(pData)
	if err != nil {
		fmt.Printf("InsertSync error: %v \r\n", err.Error())
		return false
	}
	return true
}