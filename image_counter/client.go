package image_counter

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	http   *http.Client
	logger *log.Logger
}

func NewClient(logger *log.Logger) *Client {
	cl := &Client{}
	cl.http = &http.Client{Timeout: time.Duration(10) * time.Second}

	cl.logger = logger

	return cl
}

func (c *Client) CountImages(url string) (int, error) {

	request, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return 0, err
	}

	response, err := c.http.Do(request)
	if nil != err {
		return 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return 0, err
	}

	html := string(body)
	imgCnt := strings.Count(html, "<img ")

	return imgCnt, nil
}
