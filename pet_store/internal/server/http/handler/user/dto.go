package user

// CreateUserRequest represents the request structure for creating a user
// @Description Request payload for creating a new user
type CreateUserRequest struct {
	Username  string `json:"username" example:"john_doe"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
	Password  string `json:"password" example:"password123"`
	Phone     string `json:"phone" example:"+1234567890"`
}

// CreateUserResponse represents the response structure after creating a user
// @Description Response structure containing created user details
type CreateUserResponse struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Message  string `json:"message" example:"user created successfully"`
}

// UpdateUserRequest represents the request structure for updating a user
// @Description Request payload for updating user details
type UpdateUserRequest struct {
	Username  string `json:"username" example:"john_doe"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
	Password  string `json:"password" example:"newpassword123"`
	Phone     string `json:"phone" example:"+1234567890"`
	Status    int32  `json:"userStatus" example:"1"`
}

// CreateUsersRequest represents the request structure for creating multiple users
// @Description Request payload for creating multiple users
type CreateUsersRequest struct {
	Users []CreateUserRequest `json:"users"`
}

// CreateUsersResponse represents the response structure after creating multiple users
// @Description Response structure containing created users details
type CreateUsersResponse struct {
	Message string      `json:"message" example:"users created successfully"`
	Count   int         `json:"count" example:"2"`
	Users   []UserShort `json:"users"`
}

// UserShort represents a short user information structure
// @Description Short user information structure
type UserShort struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
}