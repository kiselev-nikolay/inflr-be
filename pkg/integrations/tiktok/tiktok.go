package tiktok

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TiktokInfo struct {
	Name  string
	Login string
	Likes uint64
	Subs  uint64
}

func GetInfo(login string) (*TiktokInfo, error) {
	res, err := http.Get("https://www.tiktok.com/" + login + "?is_copy_url=1&is_from_webapp=v1&lang=en-EN")
	if err != nil {
		return nil, err
	}
	info, err := ReadWebPage(res.Body)
	if err != nil {
		return nil, err
	}
	if (TiktokInfo{}) == *info {
		return nil, errors.New("info is empty")
	}
	return info, err
}

func ReadWebPage(webPageReader io.Reader) (*TiktokInfo, error) {
	doc, err := goquery.NewDocumentFromReader(webPageReader)
	if err != nil {
		return nil, err
	}

	var docReadErr error
	info := &TiktokInfo{}

	doc.Find("meta[name=description]").Each(func(i int, s *goquery.Selection) {
		content, ok := s.Attr("content")
		if !ok {
			docReadErr = errors.New("no meta[name=description] content found")
			return
		}
		infoFromMeta, err := ReadMetaDescription(content)
		if err != nil {
			docReadErr = err
			return
		}
		*info = *infoFromMeta
	})
	if docReadErr != nil {
		return nil, docReadErr
	}

	return info, nil
}

func ReadMetaDescription(content string) (*TiktokInfo, error) {
	info := &TiktokInfo{}
	parser := regexp.MustCompile(`^(.+?) \((@.+?)\) on TikTok \| (.+?) Likes. (.+?) Fans.`) // https://regex101.com/r/9ZG22k/1
	result := parser.FindStringSubmatch(content)
	info.Name = strings.TrimSpace(result[1])
	info.Login = strings.TrimSpace(result[2])
	likes, err := translateNumber(result[3])
	if err != nil {
		return nil, err
	}
	info.Likes = likes
	subs, err := translateNumber(result[4])
	if err != nil {
		return nil, err
	}
	info.Subs = subs
	return info, nil
}

func translateNumber(number string) (x uint64, err error) {
	number = strings.ToLower(number)
	if strings.HasSuffix(number, "m") {
		f, err := strconv.ParseFloat(strings.TrimSuffix(number, "m"), 64)
		x = uint64((f * 1000000) + 0.5)
		return x, err
	}
	if strings.HasSuffix(number, "k") {
		f, err := strconv.ParseFloat(strings.TrimSuffix(number, "k"), 64)
		x = uint64((f * 1000) + 0.5)
		return x, err
	}
	x, err = strconv.ParseUint(number, 10, 64)
	return
}
