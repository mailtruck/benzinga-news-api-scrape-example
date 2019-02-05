package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	layout   = "Mon, 02 Jan 2006 15:04:05 -0700"
	pageSize = 50
)

func main() {
	var (
		err  error
		resp *http.Response
	)

	token := os.Getenv("BENZINGA_TOKEN")
	if token == "" {
		panic("missing required environment variable: BENZINGA_TOKEN")
	}

	http := http.Client{
		Timeout: time.Second * 15,
	}

	db := map[int]Story{}
	lastUpdated := time.Now().Unix()

	page := 0
POLL:
	for {
		lastUpdatedIncremented := false

		for i := 0; ; i++ {
			uri := getQueryString(token, lastUpdated, 0)

			fmt.Println(fmt.Sprintf("Requesting URI: %v", uri))
			resp, err = http.Get(uri)
			if err == nil {
				break
			}
			fmt.Println(err)

			if i >= 3 {
				panic(err)
			}
			time.Sleep(time.Duration(i) * time.Second)
		}
		defer resp.Body.Close()
		stories := []Story{}
		err = json.NewDecoder(resp.Body).Decode(&stories)
		if err != nil {
			panic(err)
		}

		for i, story := range stories {
			fmt.Println(story.Title)
			_, ok := db[story.ID]
			if ok {
				fmt.Println(fmt.Sprintf("already have: %v", story.ID))
				if story.Updated == db[story.ID].Updated {
					continue
				}
				fmt.Println("saving new update")
			} else {
				fmt.Println(fmt.Sprintf("found new story: %v", story.ID))

			}

			db[story.ID] = story
			updated, err := time.Parse(layout, story.Updated)
			if err != nil {
				panic(err)
			}
			updatedEpoch := updated.Unix()

			if updatedEpoch-1 > lastUpdated {
				lastUpdated = updatedEpoch - 1
				lastUpdatedIncremented = true
				page = 0

			} else if i == pageSize-1 && !lastUpdatedIncremented {
				fmt.Println("checking for more new content on the next page")
				page++
				continue POLL
			}

		}

		fmt.Println("sleep between requests")
		time.Sleep(time.Second * 10)
	}

}

func getQueryString(token string, lastUpdated int64, page int) string {
	params := url.Values{
		"token":        []string{token},
		"pageSize":     []string{strconv.Itoa(pageSize)},
		"updatedSince": []string{strconv.FormatInt(lastUpdated, 10)},
	}
	query := &url.URL{
		Scheme:   "https",
		Host:     "api.benzinga.com",
		Path:     "api/v2/news",
		RawQuery: params.Encode(),
	}

	return query.String()
}
