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
		Select(interface{},interface{}) bool
		SelectMany(interface{},string,interface{}) bool
		InsertOrUpdate(interface{}) bool
		Insert(interface{}) bool
		//Delete(query Query, limit int) (deleted bool)
	}
	userRepository struct {
		source *storage.DataSource
		mu     sync.RWMutex
	}
)

func NewUserRepository(source *storage.DataSource) UserRepository {
	return &userRepository{source: source}
}

func (r *userRepository)Select(search interface{},pSlice interface{}) bool {
	switch config.Task.OutType {
	case "mgo":
		return r.mgoFind(search.(bson.M),pSlice)
	case "mysql":
	}
	return false
}
func (r *userRepository)SelectMany(search interface{},sortKey string, pSlice interface{}) bool {
	switch config.Task.OutType {
	case "mgo":
		return  r.mgoFindAll(search.(bson.M),sortKey,pSlice)
	case "mysql":
	}
	return false
}

func (r *userRepository)InsertOrUpdate(pSlice interface{}) bool {
	switch config.Task.OutType {
	case "mgo":
		return r.mgoInsertOrUpdate(pSlice)
	case "mysql":
	}
	return false
}

func (r *userRepository)Insert(pSlice interface{}) bool{
	switch config.Task.OutType {
	case "mgo":
		return r.mgoInsert(pSlice)
	case "mysql":
	}
	return false
}

func (r *userRepository)mgoFind(search bson.M, pSlice interface{})  bool {
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

func (r *userRepository)mgoFindAll(search bson.M , sortKey string, pSlice interface{}) bool {
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

func (r *userRepository)mgoUpdate(id,pData interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.UpdateId(id, pData)
	if err != nil {
		//fmt.Printf("UpdateSync error: %v \r\ntable: %s \r\nid: %v \r\ndata: %v \r\n",
		//	err.Error(), table, id, pData)
		return false
	}
	return true
}

func (r *userRepository)mgoInsertOrUpdate(pData interface{}) bool {
	user := pData.(UserMgo)
	search := bson.M{"user": user.User}
	result := UserMgo{}
	if err := r.mgoFind(search,result);err{
		r.mgoUpdate(result.Id,result)
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

func (r *userRepository)mgoInsert(pData interface{}) bool {
	coll := r.source.MongoDb.Database.C(Table)
	err := coll.Insert(pData)
	if err != nil {
		fmt.Printf("InsertSync error: %v \r\n", err.Error())
		return false
	}
	return true
}