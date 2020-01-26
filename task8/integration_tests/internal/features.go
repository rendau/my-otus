package internal

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"time"
)

func (t *Tests) iRequestEventListForMonth() error {
	events, err := t.uGetEvents("month")
	if err != nil {
		return err
	}

	t.responseEvents = events

	return nil
}

func (t *Tests) iWillReceiveEventCountsInResponse(cnt int) error {
	if len(t.responseEvents) != cnt {
		return fmt.Errorf("unexpected event count: %d", len(t.responseEvents))
	}
	return nil
}

func (t *Tests) iRequestEventListForWeek() error {
	events, err := t.uGetEvents("week")
	if err != nil {
		return err
	}

	t.responseEvents = events

	return nil
}

func (t *Tests) iRequestEventListForDay() error {
	events, err := t.uGetEvents("day")
	if err != nil {
		return err
	}

	t.responseEvents = events

	return nil
}

func (t *Tests) iCreateEventWithData(data *gherkin.DocString) error {
	event := EventCE{}
	err := json.Unmarshal([]byte(data.Content), &event)
	if err != nil {
		return err
	}

	sc, errCode, err := t.uCreateEvent(event)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc
	t.responseErrCode = errCode

	return nil
}

func (t *Tests) theResponseCodeShouldBe(code int) error {
	if t.responseStatusCode != code {
		return fmt.Errorf("unexpected status code in response: %d", t.responseStatusCode)
	}
	return nil
}

func (t *Tests) iReceiveDataWithErrorCode(errCode string) error {
	if t.responseErrCode != errCode {
		return fmt.Errorf("unexpected error code in response: %s", t.responseErrCode)
	}
	return nil
}

func (t *Tests) iCreateEventForComingDayWithTitle(title string) error {
	event := EventCE{
		Owner:     new(string),
		Title:     &title,
		Text:      new(string),
		StartTime: new(time.Time),
		EndTime:   new(time.Time),
	}
	*event.Owner = "owner"
	*event.Text = "text"
	*event.StartTime = time.Now().Add(time.Minute)
	*event.EndTime = (*event.StartTime).Add(time.Hour)

	sc, errCode, err := t.uCreateEvent(event)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc
	t.responseErrCode = errCode

	return nil
}

func (t *Tests) iCreateEventForComingWeekWithTitle(title string) error {
	event := EventCE{
		Owner:     new(string),
		Title:     &title,
		Text:      new(string),
		StartTime: new(time.Time),
		EndTime:   new(time.Time),
	}
	*event.Owner = "owner"
	*event.Text = "text"
	*event.StartTime = time.Now().Add(48 * time.Hour)
	*event.EndTime = (*event.StartTime).Add(time.Hour)

	sc, errCode, err := t.uCreateEvent(event)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc
	t.responseErrCode = errCode

	return nil
}

func (t *Tests) iCreateEventForComingMonthWithTitle(title string) error {
	event := EventCE{
		Owner:     new(string),
		Title:     &title,
		Text:      new(string),
		StartTime: new(time.Time),
		EndTime:   new(time.Time),
	}
	*event.Owner = "owner"
	*event.Text = "text"
	*event.StartTime = time.Now().Add(10 * 48 * time.Hour)
	*event.EndTime = (*event.StartTime).Add(time.Hour)

	sc, errCode, err := t.uCreateEvent(event)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc
	t.responseErrCode = errCode

	return nil
}

func (t *Tests) theResponseWillContainEventWithTitle(title string) error {
	for _, event := range t.responseEvents {
		if event.Title == title {
			return nil
		}
	}
	return fmt.Errorf("title '%s' does not contain in response", title)
}

func (t *Tests) FeatureContext(s *godog.Suite) {
	s.Step(`^I request event list for month$`, t.iRequestEventListForMonth)
	s.Step(`^I will receive (\d+) event counts in response$`, t.iWillReceiveEventCountsInResponse)
	s.Step(`^I request event list for week$`, t.iRequestEventListForWeek)
	s.Step(`^I request event list for day$`, t.iRequestEventListForDay)
	s.Step(`^I create event with data:$`, t.iCreateEventWithData)
	s.Step(`^The response code should be (\d+)$`, t.theResponseCodeShouldBe)
	s.Step(`^I receive data with error code = "([^"]*)"$`, t.iReceiveDataWithErrorCode)
	s.Step(`^I create event for coming day, with title "([^"]*)"$`, t.iCreateEventForComingDayWithTitle)
	s.Step(`^I create event for coming week, with title "([^"]*)"$`, t.iCreateEventForComingWeekWithTitle)
	s.Step(`^I create event for coming month, with title "([^"]*)"$`, t.iCreateEventForComingMonthWithTitle)
	s.Step(`^The response will contain event with title "([^"]*)"$`, t.theResponseWillContainEventWithTitle)
}
