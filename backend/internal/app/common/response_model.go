package response

// MessageResponse represents a simple message response.
//
// swagger:model
type MessageResponse struct {
	Message string `json:"message" example:"logged in successfully"`
}

// ErrorResponse represents a simple error message response.
//
// swagger:model
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// UserCreatedResponse represents a response for user creation.
//
// swagger:model
type UserCreatedResponse struct {
	Message string `json:"message" example:"user registered successfully"`
}

// UserDeletedResponse represents a response for user deletion.
//
// swagger:model
type UserDeletedResponse struct {
	Message string `json:"message" example:"user deleted successfully"`
}

// UserStatusUpdatedResponse represents a response for user status update.
//
// swagger:model
type UserStatusUpdatedResponse struct {
	Message string `json:"message" example:"user status updated successfully"`
}
