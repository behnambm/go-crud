package user

import (
	"fmt"
	"github.com/behnambm/assignment/repo/sqlite"
	"log"
)

type Service struct {
	repo *sqlite.Repo
}

func New(repo *sqlite.Repo) Service {
	return Service{repo: repo}
}

func (s Service) GetUserFromEmail(email string) (User, error) {
	userFromDB, err := s.repo.GetUserFromEmail(email)
	if err != nil {
		log.Println("USER SERVICE ERR", err)
		return User{}, fmt.Errorf("couldn't find user")
	}
	user := User{
		ID:       userFromDB.ID,
		Email:    userFromDB.Email,
		Password: userFromDB.Password,
		IsAdmin:  userFromDB.IsAdmin,
	}
	return user, nil
}

func (s Service) GetUserFromID(id int) (User, error) {
	userFromDB, err := s.repo.GetUserFromID(id)
	if err != nil {
		log.Println("USER SERVICE ERR", err)
		return User{}, fmt.Errorf("couldn't find user")
	}
	user := User{
		ID:       userFromDB.ID,
		Email:    userFromDB.Email,
		Password: userFromDB.Password,
		IsAdmin:  userFromDB.IsAdmin,
	}
	return user, nil
}
