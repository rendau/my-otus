package events

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rendau/my-otus/task8/internal/adapters/rest/util"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"net/http"
	"strconv"
	"time"
)

func hListForDay(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)

	hListForFilter(w, r, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

func hListForWeek(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(7 * 24 * time.Hour)

	hListForFilter(w, r, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

func hListForMonth(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(30 * 24 * time.Hour)

	hListForFilter(w, r, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

func hListForFilter(w http.ResponseWriter, r *http.Request, filter *entities.EventListFilter) {
	apiCtx := util.GetAPICtx(r)

	events, err := apiCtx.Ucs.Event.List(context.Background(), filter)
	if err != nil {
		util.Respond500(w)
		return
	}

	util.RespondAppJSONObj(w, 200, events, nil)
}

func hCreate(w http.ResponseWriter, r *http.Request) {
	var err error

	apiCtx := util.GetAPICtx(r)

	pOwner := r.FormValue("owner")
	pTitle := r.FormValue("title")
	pText := r.FormValue("text")
	pStartTime := r.FormValue("start_time")
	pEndTime := r.FormValue("end_time")

	var startTime time.Time
	var endTime time.Time

	if pStartTime != "" {
		startTime, err = time.Parse(time.RFC3339, pStartTime)
		if err != nil {
			util.Respond400(w, "bad_start_time", "Wrong format for start-time")
			return
		}
	}

	if pEndTime != "" {
		endTime, err = time.Parse(time.RFC3339, pEndTime)
		if err != nil {
			util.Respond400(w, "bad_end_time", "Wrong format for end-time")
			return
		}
	}

	event := &entities.Event{
		Owner:     pOwner,
		Title:     pTitle,
		Text:      pText,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err = apiCtx.Ucs.Event.Create(context.Background(), event)
	if err != nil {
		switch err.(type) {
		case errors.EventError:
			util.RespondAppJSONObj(w, 200, nil, err.Error())
		default:
			util.Respond500(w)
		}
		return
	}

	util.RespondAppJSONObj(w, 200, struct {
		NewID int64 `json:"new_id"`
	}{
		NewID: event.ID,
	}, nil)
}

func hUpdate(w http.ResponseWriter, r *http.Request) {
	var err error

	apiCtx := util.GetAPICtx(r)

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	pOwner := r.FormValue("owner")
	pTitle := r.FormValue("title")
	pText := r.FormValue("text")
	pStartTime := r.FormValue("start_time")
	pEndTime := r.FormValue("end_time")

	var startTime time.Time
	var endTime time.Time

	if pStartTime != "" {
		startTime, err = time.Parse(time.RFC3339, pStartTime)
		if err != nil {
			util.Respond400(w, "bad_start_time", "Wrong format for start-time")
			return
		}
	}

	if pEndTime != "" {
		endTime, err = time.Parse(time.RFC3339, pEndTime)
		if err != nil {
			util.Respond400(w, "bad_end_time", "Wrong format for end-time")
			return
		}
	}

	event := &entities.Event{
		ID:        id,
		Owner:     pOwner,
		Title:     pTitle,
		Text:      pText,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err = apiCtx.Ucs.Event.Update(context.Background(), event)
	if err != nil {
		switch err.(type) {
		case errors.EventError:
			util.RespondAppJSONObj(w, 200, nil, err.Error())
		default:
			util.Respond500(w)
		}
		return
	}

	w.WriteHeader(200)
}

func hDelete(w http.ResponseWriter, r *http.Request) {
	var err error

	apiCtx := util.GetAPICtx(r)

	args := mux.Vars(r)
	id, _ := strconv.ParseInt(args["id"], 10, 64)

	err = apiCtx.Ucs.Event.Delete(
		context.Background(),
		id,
	)
	if err != nil {
		switch err.(type) {
		case errors.EventError:
			util.RespondAppJSONObj(w, 200, nil, err.Error())
		default:
			util.Respond500(w)
		}
		return
	}

	w.WriteHeader(200)
}
