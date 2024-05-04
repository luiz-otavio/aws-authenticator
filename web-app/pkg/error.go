package pkg

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type UnexpetechedError struct {
	Err     error  `json:"error"`
	Message string `json:"message"`
}

func (e UnexpetechedError) String() string {
	marshaled, err := json.Marshal(e)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal unexpected error")

		os.Exit(1)
	}

	return string(marshaled)
}

func NewUnexpectedError(err error, message string) UnexpetechedError {
	return UnexpetechedError{Err: err, Message: message}
}
