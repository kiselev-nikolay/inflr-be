package profile

import (
	"net/url"

	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/telegram"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/tiktok"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations/youtube"
	"github.com/kiselev-nikolay/inflr-be/pkg/repository"
	irepository "github.com/kiselev-nikolay/inflr-be/pkg/repository/interfaces"
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
	Country      Country
}
type Projects struct {
	Youtube  []youtube.YoutubeInfo
	Telegram []telegram.TelegramInfo
	Tiktok   []tiktok.TiktokInfo
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
	List   func() ([]*Profile, error)
}

func NewModel(repo repository.Repo) *Model {
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
	model.List = func() ([]*Profile, error) {
		items, err := repo.List(collection)
		if err != nil {
			return nil, err
		}
		profiles := make([]*Profile, 0)
		for _, item := range items {
			profile := item.Value.(Profile)
			profiles = append(profiles, &profile)
		}
		return profiles, nil
	}
	return model
}
