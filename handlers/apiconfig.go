package handlers

import (
	"sync/atomic"

	"github.com/Bigesto/chirpy/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
	PolkaKey       string
}
