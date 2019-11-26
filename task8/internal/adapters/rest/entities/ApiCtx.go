package entities

import (
	"github.com/rendau/my-otus/task8/internal/domain/usecases"
	"github.com/rendau/my-otus/task8/internal/interfaces"
)

// APICtx - context type in http request
type APICtx struct {
	Ucs *usecases.Usecases
	Log interfaces.Logger
}
