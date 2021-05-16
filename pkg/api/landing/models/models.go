package models

import (
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
)

type Card struct {
	Title, Text, Link string
}

type Landing struct {
	Key   string
	Title string
	Cards []Card
}

type Model struct {
	Send   func(*user.User, string, *Landing) error
	Find   func(*user.User, string) (*Landing, error)
	Delete func(*user.User, string) error
	List   func(*user.User) ([]*Landing, error)
}

func New(repo repository.Repo) *Model {
	collection := "Landing"
	model := &Model{}
	model.Send = func(u *user.User, k string, v *Landing) error {
		return repo.Send(collection+":"+u.Login, &irepository.Item{Key: k, Value: *v})
	}
	model.Find = func(u *user.User, k string) (*Landing, error) {
		i := irepository.Item{Key: k}
		err := repo.Find(collection+":"+u.Login, &i)
		if err != nil {
			return nil, err
		}
		v := i.Value.(Landing)
		return &v, nil
	}
	model.Delete = func(u *user.User, k string) error {
		return repo.Delete(collection+":"+u.Login, &irepository.Item{Key: k})
	}
	model.List = func(u *user.User) ([]*Landing, error) {
		items, err := repo.List(collection + ":" + u.Login)
		if err != nil {
			return nil, err
		}
		landings := make([]*Landing, 0)
		for _, item := range items {
			landing := item.Value.(Landing)
			landings = append(landings, &landing)
		}
		return landings, nil
	}
	return model
}
