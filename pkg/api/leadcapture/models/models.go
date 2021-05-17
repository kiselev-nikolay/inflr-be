package models

import (
	"github.com/kiselev-nikolay/inflr-be/pkg/authware/user"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
)

type Lead struct {
	LandingKey string
	Title      string
	Email      string
	Phone      string
	Forms      map[string]string
}

type Model struct {
	Send   func(*user.User, string, *Lead) error
	Find   func(*user.User, string) (*Lead, error)
	Delete func(*user.User, string) error
	List   func(*user.User) ([]*Lead, error)
}

func New(repo repository.Repo) *Model {
	collection := "Lead"
	model := &Model{}
	model.Send = func(u *user.User, k string, v *Lead) error {
		return repo.Send(collection+":"+u.Login, &irepository.Item{Key: k, Value: *v})
	}
	model.Find = func(u *user.User, k string) (*Lead, error) {
		i := irepository.Item{Key: k}
		err := repo.Find(collection+":"+u.Login, &i)
		if err != nil {
			return nil, err
		}
		v := i.Value.(Lead)
		return &v, nil
	}
	model.Delete = func(u *user.User, k string) error {
		return repo.Delete(collection+":"+u.Login, &irepository.Item{Key: k})
	}
	model.List = func(u *user.User) ([]*Lead, error) {
		items, err := repo.List(collection + ":" + u.Login)
		if err != nil {
			return nil, err
		}
		leads := make([]*Lead, 0)
		for _, item := range items {
			lead := item.Value.(Lead)
			leads = append(leads, &lead)
		}
		return leads, nil
	}
	return model
}
