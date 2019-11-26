package transport

import (
	"autocare.org/sandpiper/pkg/model"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*sandpiper.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []sandpiper.User `json:"users"`
		Page  int              `json:"page"`
	}
}
