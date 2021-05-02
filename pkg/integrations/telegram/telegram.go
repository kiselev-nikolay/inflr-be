package telegram

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/kiselev-nikolay/inflr-be/pkg/integrations"
)

const htmlSubsSelector = ".tgme_page_extra"
const htmlPhotoSelector = ".tgme_page_photo_image"

func GetInfo(address string) (*integrations.Info, error) {
	tmeResponse, err := http.Get("https://t.me/" + address)
	if err != nil {
		return nil, err
	}
	info, err := ReadWebPage(tmeResponse.Body)
	if err != nil {
		return nil, err
	}
	if (integrations.Info{}) == *info {
		return nil, errors.New("info is empty")
	}
	return info, err
}

func ReadWebPage(webPageReader io.Reader) (*integrations.Info, error) {
	doc, err := goquery.NewDocumentFromReader(webPageReader)

	if err != nil {
		return nil, err
	}
	var docReadErr error
	info := &integrations.Info{}
	doc.Find(htmlSubsSelector).Each(func(i int, s *goquery.Selection) {
		rawValue := s.Text()
		value := make([]rune, 0)
		for _, l := range rawValue {
			switch l {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				value = append(value, l)
			}
		}
		subs, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			docReadErr = err
			return
		}
		info.Subs = subs
	})
	if docReadErr != nil {
		return nil, docReadErr
	}
	doc.Find(htmlPhotoSelector).Each(func(i int, s *goquery.Selection) {
		imageSource, exists := s.Attr("src")
		if !exists {
			docReadErr = errors.New("element has no src")
			return
		}
		imageLink, err := url.Parse(imageSource)
		if imageLink.Host == "" {
			imageLink, err = url.Parse("https://t.me/" + imageSource)
			if err != nil {
				docReadErr = err
				return
			}
		}
		if err != nil {
			docReadErr = err
			return
		}
		info.Image = *imageLink
	})
	if docReadErr != nil {
		return nil, docReadErr
	}
	return info, nil
}
