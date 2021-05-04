package youtube

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type pageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type statistics struct {
	ViewCount       string `json:"viewCount"`
	SubscriberCount string `json:"subscriberCount"`
	VideoCount      string `json:"videoCount"`
}

type image struct {
	URL string `json:"url"`
}

type thumbnails struct {
	Default image `json:"default"`
	Medium  image `json:"medium"`
	High    image `json:"high"`
}

type snippet struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Published   string     `json:"publishedAt"`
	Thumbnails  thumbnails `json:"thumbnails"`
}

type channel struct {
	Id         string     `json:"id"`
	Statistics statistics `json:"statistics"`
	Snippet    snippet    `json:"snippet"`
}

type response struct {
	PageInfo pageInfo  `json:"pageInfo"`
	Items    []channel `json:"items"`
}

type YoutubeInfo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Register    time.Time `json:"register"`
	ImageURL    string    `json:"imageUrl"`
	Subs        uint64    `json:"subs"`
	Views       uint64    `json:"views"`
	Videos      uint64    `json:"videos"`
}

func GetInfo(ytid string) (*YoutubeInfo, error) {
	yourKey := "AIzaSyBQQ-zTp3e4o0GkJEbnnmH35hTMOSxsW_E"
	q := "key=" + yourKey
	q += "&part=statistics&part=snippet"
	if !(strings.HasPrefix(ytid, "UC") && len(ytid) == 24) {
		return nil, errors.New("invalid ytid")
	}
	q += "&id=" + ytid
	url := "https://www.googleapis.com/youtube/v3/channels?" + q
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	info, err := ReadYTResponse(res.Body)
	if err != nil {
		return nil, err
	}
	if (YoutubeInfo{}) == *info {
		return nil, errors.New("info is empty")
	}
	return info, nil
}

func ReadYTResponse(ytResponse io.Reader) (*YoutubeInfo, error) {
	bodyBytes, err := ioutil.ReadAll(ytResponse)
	if err != nil {
		return nil, err
	}
	data := &response{}
	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		return nil, err
	}
	subs, err := strconv.ParseUint(data.Items[0].Statistics.SubscriberCount, 10, 64)
	if err != nil {
		return nil, err
	}
	views, err := strconv.ParseUint(data.Items[0].Statistics.ViewCount, 10, 64)
	if err != nil {
		return nil, err
	}
	videos, err := strconv.ParseUint(data.Items[0].Statistics.VideoCount, 10, 64)
	if err != nil {
		return nil, err
	}
	register, err := time.Parse(time.RFC3339, data.Items[0].Snippet.Published)
	if err != nil {
		return nil, err
	}
	info := &YoutubeInfo{
		Title:       data.Items[0].Snippet.Title,
		Description: data.Items[0].Snippet.Description,
		Register:    register,
		ImageURL:    data.Items[0].Snippet.Thumbnails.High.URL,
		Subs:        subs,
		Views:       views,
		Videos:      videos,
	}
	return info, err
}

func GetYTIDFromLink(url *url.URL) (string, error) {
	pathSeg := strings.Split(strings.Trim(url.Path, "/"), "/")
	lastPathSeg := pathSeg[len(pathSeg)-1]
	if strings.HasPrefix(lastPathSeg, "UC") && len(lastPathSeg) == 24 {
		return lastPathSeg, nil
	}
	return "", errors.New("not a channel id")
}
