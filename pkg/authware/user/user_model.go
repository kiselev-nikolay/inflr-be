package user

import (
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
)

type User struct {
	Login          string
	SecretPassword []byte
}

type UserModel struct {
	Send   func(string, *User) error
	Find   func(string) (*User, error)
	Delete func(string) error
	List   func() ([]*User, error)
}

func NewUserModel(repo repository.Repo) *UserModel {
	collection := "User"
	model := &UserModel{}
	model.Send = func(k string, v *User) error {
		return repo.Send(collection, &irepository.Item{Key: k, Value: *v})
	}
	model.Find = func(k string) (*User, error) {
		i := irepository.Item{Key: k}
		err := repo.Find(collection, &i)
		if err != nil {
			return nil, err
		}
		v := i.Value.(User)
		return &v, nil
	}
	model.Delete = func(k string) error {
		return repo.Delete(collection, &irepository.Item{Key: k})
	}
	model.List = func() ([]*User, error) {
		items, err := repo.List(collection)
		if err != nil {
			return nil, err
		}
		users := make([]*User, 0)
		for _, item := range items {
			user := item.Value.(User)
			users = append(users, &user)
		}
		return users, nil
	}
	return model
}
