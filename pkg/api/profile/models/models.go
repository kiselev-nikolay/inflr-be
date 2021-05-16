package models

import (
	"net/url"

	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/telegram"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/tiktok"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/youtube"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
	"github.com/kiselev-nikolay/inflr-be/pkg/tools/country"
)

const (
	AvailabilityOpen = iota
	AvailabilityBusy = iota
	AvailabilityGone = iota
)

type Bio struct {
	Name         string
	About        string
	Links        []url.URL
	Availability int
	Country      country.Country
}
type Projects struct {
	Youtube  map[string]youtube.YoutubeInfo
	Telegram map[string]telegram.TelegramInfo
	Tiktok   map[string]tiktok.TiktokInfo
}

type Profile struct {
	Login   string
	CostEUR float64
	Bio
	Projects
}

type Model struct {
	Send   func(string, *Profile) error
	Find   func(string) (*Profile, error)
	Delete func(string) error
}

func New(repo repository.Repo) *Model {
	collection := "Profile"
	model := &Model{}
	model.Send = func(k string, v *Profile) error {
		return repo.Send(collection, &irepository.Item{Key: k, Value: *v})
	}
	model.Find = func(k string) (*Profile, error) {
		i := irepository.Item{Key: k}
		err := repo.Find(collection, &i)
		if err != nil {
			return nil, err
		}
		v := i.Value.(Profile)
		return &v, nil
	}
	model.Delete = func(k string) error {
		return repo.Delete(collection, &irepository.Item{Key: k})
	}
	return model
}
