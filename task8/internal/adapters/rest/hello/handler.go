package hello

import (
	"github.com/rendau/my-otus/task8/internal/adapters/rest/util"
	"net/http"
)

func hGet(w http.ResponseWriter, r *http.Request) {
	util.RespondStr(w, 200, "Hi World!")
}
