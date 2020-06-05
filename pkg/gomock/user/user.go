package user

import (
	"github.com/t-xinlin/doc/pkg/gomock/person"
)
type User struct {
	Person person.Male
}

func NewUser(p person.Male) *User {
	return &User{Person: p}
}

func (u *User) GetUserInfo(id int64) error {
	return u.Person.Get(id)
}

func (u *User) PutUserInfo(id int64) error {
	return u.Person.Put(id)
}