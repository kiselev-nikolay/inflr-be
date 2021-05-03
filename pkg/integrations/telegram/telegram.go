package telegram

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const htmlTitleSelector = ".tgme_page_title"
const htmlTitleDescription = ".tgme_page_description"
const htmlSubsSelector = ".tgme_page_extra"
const htmlImageSelector = ".tgme_page_photo_image"

type TelegramInfo struct {
	Title       string
	Description string
	ImageURL    string
	Subs        uint64
}

func GetInfo(address string) (*TelegramInfo, error) {
	tmeResponse, err := http.Get("https://t.me/" + address)
	if err != nil {
		return nil, err
	}
	info, err := ReadWebPage(tmeResponse.Body)
	if err != nil {
		return nil, err
	}
	if (TelegramInfo{}) == *info {
		return nil, errors.New("info is empty")
	}
	return info, err
}

func ReadWebPage(webPageReader io.Reader) (*TelegramInfo, error) {
	doc, err := goquery.NewDocumentFromReader(webPageReader)
	if err != nil {
		return nil, err
	}

	var docReadErr error
	info := &TelegramInfo{}

	doc.Find(htmlTitleSelector).Each(func(i int, s *goquery.Selection) {
		info.Title = strings.TrimSpace(s.Text())
	})
	if docReadErr != nil {
		return nil, docReadErr
	}

	doc.Find(htmlTitleDescription).Each(func(i int, s *goquery.Selection) {
		texts := make([]string, 0)
		s.Contents().Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				texts = append(texts, text)
			}
		})
		info.Description = strings.Join(texts, " ")
	})
	if docReadErr != nil {
		return nil, docReadErr
	}

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

	doc.Find(htmlImageSelector).Each(func(i int, s *goquery.Selection) {
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
		info.ImageURL = imageLink.String()
	})
	if docReadErr != nil {
		return nil, docReadErr
	}

	return info, nil
}
