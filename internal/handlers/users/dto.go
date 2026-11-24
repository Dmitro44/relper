package users

import (
	"github.com/Dmitro44/relper/internal/handlers/dto"
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Email    *string `json:"email,omitempty" pattern:"^[^@\\s]+@[^@\\s]+$" doc:"User email"`
		Nickname *string `json:"nickname,omitempty" minLength:"1" maxLength:"50" doc:"Display name"`
	}
}

type GetUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `path:"id"`
}

type UserResponse struct {
	Body dto.UserResponse
}

type SetUserRolesRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Roles []string `json:"roles" minItems:"1" doc:"Role codes to assign (replaces existing)"`
	}
}

type UserRolesResponse struct {
	Body struct {
		Roles []string `json:"roles"`
	}
}
