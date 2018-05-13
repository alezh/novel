package services

import (
	"github.com/alezh/novel/app/admin/repositories"
)
type (
	UserService interface {
		//GetAll()
		//GetByID(id interface{}) (User, bool)
		//GetByUsernameAndPassword(username, userPassword string)
		//DeleteByID(id int64) bool
		//Update(id int64, user User) (bool, error)
		//UpdatePassword(id int64, newPassword string) (bool, error)
		//UpdateUsername(id int64, newUsername string) (bool, error)
		//Create(userPassword string, user User) (bool, error)
	}
	userService struct {
		repo repositories.UserRepository
	}
)
func NewUserService(source  repositories.UserRepository) UserService {
	return &userService{
		repo: source,
	}
}
//func (s *userService) Create(user,userPassword string) (User, bool) {
//	//s.repo.InsertOrUpdate()
//}