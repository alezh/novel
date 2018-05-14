package services

import (
	"github.com/alezh/novel/app/admin/repositories"
)
type (
	UserService interface {
		//GetAll()
		//GetByID(id interface{}) (User, bool)
		GetByUsernameAndPassword(username, userPassword string) (interface{},bool)
		//DeleteByID(id int64) bool
		//Update(id int64, user User) (bool, error)
		//UpdatePassword(id int64, newPassword string) (bool, error)
		//UpdateUsername(id int64, newUsername string) (bool, error)
		CreateAdmin(user,userPassword string) bool
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
func (s *userService) CreateAdmin(user,userPassword string) bool {
	return s.repo.InsertUser(user,userPassword,1)
}

func (s *userService)GetByUsernameAndPassword(username, userPassword string) (interface{},bool) {
	if username == "" || userPassword == "" {
		return  nil,false
	}
	return  s.repo.Verify(username, userPassword)
}