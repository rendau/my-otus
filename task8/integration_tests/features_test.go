package main

import "github.com/DATA-DOG/godog"

func iRequestEventListForMonth() error {
	return godog.ErrPending
}

func iWillReceiveEventCountsInResponse(arg1 int) error {
	return godog.ErrPending
}

func iRequestEventListForWeek() error {
	return godog.ErrPending
}

func iRequestEventListForDay() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I request event list for month$`, iRequestEventListForMonth)
	s.Step(`^I will receive (\d+) event counts in response$`, iWillReceiveEventCountsInResponse)
	s.Step(`^I request event list for week$`, iRequestEventListForWeek)
	s.Step(`^I request event list for day$`, iRequestEventListForDay)
}
