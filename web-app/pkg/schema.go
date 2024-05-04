package pkg

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpSchema struct {
	CommitedAt time.Time `json:"commited_at"`
}

func (schema HttpSchema) String() string {
	body, err := json.Marshal(schema)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal http schema")

		return ""
	}

	return string(body)
}

type UserRequestLoginSchema struct {
	HttpSchema

	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponseLoginSchema struct {
	HttpSchema

	Message string `json:"message"`
}

type UserChangePasswordRequestSchema struct {
	HttpSchema

	Username    string `json:"username"`
	OldPassword string `json:"password"`
	NewPassword string `json:"new_password"`
}

type UserChangePasswordResponseSchema struct {
	HttpSchema

	Message string `json:"message"`
}

type UserRegisterRequestSchema struct {
	HttpSchema

	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterResponseSchema struct {
	HttpSchema

	Message string `json:"message"`
}

type ExistsUserRequestSchema struct {
	HttpSchema

	Username string `json:"username"`
}

type ExistsUserResponseSchema struct {
	HttpSchema

	Message string `json:"message"`
}
