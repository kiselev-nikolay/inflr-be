package models

import (
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
)

type Card struct {
	Title, Text, Link string
}

type Landing struct {
	Title string
	Cards []Card
}

type Model struct {
	Send   func(string, *Landing) error
	Find   func(string) (*Landing, error)
	Delete func(string) error
	List   func() ([]*Landing, error)
}

func New(repo repository.Repo) *Model {
	collection := "Landing"
	model := &Model{}
	model.Send = func(k string, v *Landing) error {
		return repo.Send(collection, &irepository.Item{Key: k, Value: *v})
	}
	model.Find = func(k string) (*Landing, error) {
		i := irepository.Item{Key: k}
		err := repo.Find(collection, &i)
		if err != nil {
			return nil, err
		}
		v := i.Value.(Landing)
		return &v, nil
	}
	model.Delete = func(k string) error {
		return repo.Delete(collection, &irepository.Item{Key: k})
	}
	model.List = func() ([]*Landing, error) {
		items, err := repo.List(collection)
		if err != nil {
			return nil, err
		}
		profiles := make([]*Landing, 0)
		for _, item := range items {
			profile := item.Value.(Landing)
			profiles = append(profiles, &profile)
		}
		return profiles, nil
	}
	return model
}
