package util

import (
	"encoding/json"
	"github.com/rendau/my-otus/task8/internal/adapters/rest/constants"
	"github.com/rendau/my-otus/task8/internal/adapters/rest/entities"
	"log"
	"net/http"
)

// GetAPICtx - retrieves APICtx-context from request
func GetAPICtx(r *http.Request) *entities.APICtx {
	return r.Context().Value(constants.APICtxKey).(*entities.APICtx)
}

// SetContentTypeJSON - sets content-type of response to json
func SetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// RespondStr - sends string
func RespondStr(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	if len(body) > 0 {
		_, _ = w.Write([]byte(body))
	}
}

// RespondJSONObj - sends struct as json
func RespondJSONObj(w http.ResponseWriter, code int, obj interface{}) {
	SetContentTypeJSON(w)
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Println("Fail to encode json for response:", err)
	}
}

// RespondAppJSONObj - sends application reply struct as json
func RespondAppJSONObj(w http.ResponseWriter, code int, result interface{}, rErr interface{}) {
	if rErr != nil {
		RespondJSONObj(w, code, struct {
			Error interface{} `json:"error"`
		}{
			Error: rErr,
		})
	} else {
		RespondJSONObj(w, code, struct {
			Result interface{} `json:"result"`
		}{
			Result: result,
		})
	}
}

// ParseJSONObj - parses json to struct from request
func ParseJSONObj(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		RespondJSONParseError(w)
		return false
	}
	return true
}

// RespondError - sends error
func RespondError(w http.ResponseWriter, code int, err string, detail string) {
	obj := map[string]string{
		"error":     err,
		"error_dsc": detail,
	}
	RespondJSONObj(w, code, obj)
}

// Respond400 - sends 400 error
func Respond400(w http.ResponseWriter, err, detail string) {
	RespondError(w, 400, err, detail)
}

// Respond401 - sends 401 error
func Respond401(w http.ResponseWriter, detail string) {
	RespondError(w, 401, "unauthorized", detail)
}

// Respond403 - sends 403 error
func Respond403(w http.ResponseWriter, detail string) {
	RespondError(w, 403, "permission_denied", detail)
}

// Respond404 - sends 404 error
func Respond404(w http.ResponseWriter, detail string) {
	RespondError(w, 404, "not_found", detail)
}

// Respond500 - sends 500 error
func Respond500(w http.ResponseWriter) {
	RespondError(w, 500, "Internal Server Error", "")
}

// RespondJSONParseError - sends " parse" error
func RespondJSONParseError(w http.ResponseWriter) {
	Respond400(w, "bad_json", "Fail to parse JSON")
}
