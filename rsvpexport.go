package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type MetaFields struct {
	Count      int    `json:"count"`
	TotalCount int    `json:"total_count"`
	Next       string `json:"next"`
}

type Meta struct {
	Meta *MetaFields `json:"meta"`
}

type Member struct {
	Name     string `json:"name"`
	MemberId int    `json:"member_id"`
}

type Group struct {
	Id int `json:"id"`
}

type Result struct {
	Response string  `json:"response"`
	Member   *Member `json:"member"`
	Guests   int     `json:"guests"`
	Group    *Group  `json:"group"`
}

type Results struct {
	Results []*Result `json:"results"`
}

type Profile struct {
	Bio string `json:"bio"`
}

func (results *Results) AppendResults(extraResults []*Result) {
	results.Results = append(results.Results, extraResults...)
}

var eventId string
var apiKey string

func init() {
	flag.StringVar(&eventId, "eventid", "", "the event-id of the meetup")
	flag.StringVar(&apiKey, "apikey", "", "your meetup API key")
	flag.Parse()
}

func main() {
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	rsvps := getRSVPS()

	csvWriter.Write([]string{"Member name", "Profile bio", "RSVP status", "Guests"})

	for _, result := range rsvps.Results {
		profile := getProfile(result.Group.Id, result.Member.MemberId)
		csvWriter.Write([]string{result.Member.Name, profile.Bio, result.Response, strconv.Itoa(result.Guests)})
	}

}

func getRSVPS() *Results {
	v := url.Values{}
	v.Set("sign", "true")
	v.Set("event_id", string(eventId))
	v.Set("page", "20")
	v.Set("key", apiKey)
	requestUrl := fmt.Sprintf("https://api.meetup.com/2/rsvps?%s", v.Encode())

	output := new(Results)

	for {
		log.Printf("GET: %s", requestUrl)

		// make request
		resp, err := http.Get(requestUrl)
		if err != nil {
			log.Panic(err)
		} else if resp.StatusCode != 200 {
			log.Fatalf("Request failed: %d - %s\n", resp.StatusCode, resp.Status)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// unmarshal meta part (for next link, in case of new page)
		var meta Meta
		if err := json.Unmarshal(body, &meta); err != nil {
			log.Fatal(err)
		}

		// umarshal results
		var results Results
		if err := json.Unmarshal(body, &results); err != nil {
			log.Fatal(err)
		}

		output.AppendResults(results.Results)

		if meta.Meta.Next != "" {
			requestUrl = meta.Meta.Next
		} else {
			break
		}
	}

	return output
}

func getProfile(groupId, memberId int) *Profile {
	v := url.Values{}
	v.Set("sign", "true")
	v.Set("key", apiKey)
	requestUrl := fmt.Sprintf("https://api.meetup.com/2/profile/%d/%d?%s", groupId, memberId, v.Encode())

	log.Printf("GET: %s", requestUrl)

	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Panic(err)
	} else if resp.StatusCode != 200 {
		log.Fatalf("Request failed: %d - %s\n", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Fatal(err)
	}

	return &profile
}
