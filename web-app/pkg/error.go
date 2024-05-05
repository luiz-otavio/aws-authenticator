package pkg

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type UnexpectedError struct {
	Err     error  `json:"error"`
	Message string `json:"message"`
}

func (e UnexpectedError) String() string {
	marshaled, err := json.Marshal(e)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal unexpected error")

		os.Exit(1)
	}

	return string(marshaled)
}

func NewUnexpectedError(err error, message string) UnexpectedError {
	return UnexpectedError{Err: err, Message: message}
}
