package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

func (t *Tests) uGetEvents(period string) ([]*Event, error) {
	rep, err := http.Get(t.apiUrl + "/events/for_" + period)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	repDataObj := struct {
		Result []*Event `json:"result"`
	}{}
	if err = json.NewDecoder(rep.Body).Decode(&repDataObj); err != nil {
		return nil, err
	}

	return repDataObj.Result, nil
}

func (t *Tests) uCreateEvent(event EventCE) (int, string, error) {
	form := url.Values{}
	if event.Owner != nil {
		form.Add("owner", *event.Owner)
	}
	if event.Title != nil {
		form.Add("title", *event.Title)
	}
	if event.Text != nil {
		form.Add("text", *event.Text)
	}
	if event.StartTime != nil {
		form.Add("start_time", (*event.StartTime).Format(time.RFC3339))
	}
	if event.EndTime != nil {
		form.Add("end_time", (*event.EndTime).Format(time.RFC3339))
	}

	rep, err := http.PostForm(t.apiUrl+"/events", form)
	if err != nil {
		return 0, "", err
	}
	defer rep.Body.Close()

	if rep.StatusCode != 200 {
		return rep.StatusCode, "", nil
	}

	repDataObj := struct {
		Error string `json:"error"`
	}{}
	if err = json.NewDecoder(rep.Body).Decode(&repDataObj); err != nil {
		return 0, "", err
	}

	return rep.StatusCode, repDataObj.Error, nil
}
