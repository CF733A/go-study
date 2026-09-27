package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"pet-store/internal/entities"
	"pet-store/internal/server/http/response"
	services "pet-store/internal/services"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User object that needs to be created"
// @Success 201 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} response.Response "Invalid input or missing required fields"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 409 {object} response.Response "User already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user [post]
// @Example response
// {
//   "id": 1,
//   "username": "john_doe",
//   "message": "user created successfully"
// }
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("username, email and password are required"))
		return
	}

	user, err := entities.NewUser(
		req.Username,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		req.Phone,
	)
	if err != nil {
		response.JSONer(w, http.StatusInternalServerError, nil, errors.New("failed to create user"))
		return
	}

	if err := h.userService.CreateUser(r.Context(), user); err != nil {
		if err.Error() == "user already exists" {
			response.JSONer(w, http.StatusConflict, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	responseData := CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Message:  "user created successfully",
	}

	response.JSONer(w, http.StatusCreated, responseData, nil)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	if username == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("username is required"))
		return
	}

	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	user.Password = ""
	response.JSONer(w, http.StatusOK, user, nil)
}

// GetUserByUsername godoc
// @Summary Get user by username
// @Description Returns a single user by username
// @Tags user
// @Produce json
// @Param username path string true "The name that needs to be fetched"
// @Success 200 {object} entities.User "Successful operation"
// @Failure 400 {object} response.Response "Invalid username supplied"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/{username} [get]
// @Example response
// {
//   "id": 1,
//   "username": "john_doe",
//   "firstName": "John",
//   "lastName": "Doe",
//   "email": "john@example.com",
//   "phone": "+1234567890",
//   "userStatus": 1
// }
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	if username == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("username is required"))
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	updatedUser := &entities.User{
		Username:   req.Username,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   req.Password,
		Phone:      req.Phone,
		UserStatus: req.Status,
	}

	if err := h.userService.UpdateUser(r.Context(), username, updatedUser); err != nil {
		if err.Error() == "user not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "user updated successfully"}, nil)
}

// DeleteUser godoc
// @Summary Delete user by username
// @Description Delete user by username
// @Tags user
// @Produce json
// @Param username path string true "The username that needs to be deleted"
// @Success 200 {object} response.Response "User deleted successfully"
// @Failure 400 {object} response.Response "Invalid username supplied"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/{username} [delete]
// @Example response
// {
//   "status": true,
//   "message": "user deleted successfully"
// }
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	if username == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("username is required"))
		return
	}

	if err := h.userService.DeleteUser(r.Context(), username); err != nil {
		if err.Error() == "user not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "user deleted successfully"}, nil)
}

// CreateUsersWithArray обрабатывает создание пользователей из массива
// @Summary Creates list of users with given input array
// @Description Creates list of users with given input array
// @Tags user
// @Accept json
// @Produce json
// @Param users body []CreateUserRequest true "Array of user objects"
// @Success 201 {object} CreateUsersResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/createWithArray [post]
func (h *UserHandler) CreateUsersWithArray(w http.ResponseWriter, r *http.Request) {
	h.createUsers(w, r)
}

// CreateUsersWithList обрабатывает создание пользователей из списка
// @Summary Creates list of users with given input list
// @Description Creates list of users with given input list
// @Tags user
// @Accept json
// @Produce json
// @Param users body []CreateUserRequest true "List of user objects"
// @Success 201 {object} CreateUsersResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/createWithList [post]
func (h *UserHandler) CreateUsersWithList(w http.ResponseWriter, r *http.Request) {
	h.createUsers(w, r)
}

// createUsers godoc
// @Summary Create multiple users
// @Description Create multiple users with the provided details
// @Tags user
// @Accept json
// @Produce json
// @Param users body []CreateUserRequest true "Array of user objects that need to be created"
// @Success 201 {object} CreateUsersResponse "Users created successfully"
// @Failure 400 {object} response.Response "Invalid input or missing required fields"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 409 {object} response.Response "Username already exists"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/createWithList [post]
// @Example response
// {
//   "message": "users created successfully",
//   "count": 2,
//   "users": [
//     {
//       "id": 1,
//       "username": "john_doe"
//     },
//     {
//       "id": 2,
//       "username": "jane_smith"
//     }
//   ]
// }
func (h *UserHandler) createUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var userRequests []CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequests); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if len(userRequests) == 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("at least one user is required"))
		return
	}

	users := make([]*entities.User, 0, len(userRequests))
	for _, req := range userRequests {
		if req.Username == "" || req.Email == "" || req.Password == "" {
			response.JSONer(w, http.StatusBadRequest, nil, errors.New("username, email and password are required for all users"))
			return
		}

		user, err := entities.NewUser(
			req.Username,
			req.FirstName,
			req.LastName,
			req.Email,
			req.Password,
			req.Phone,
		)
		if err != nil {
			response.JSONer(w, http.StatusInternalServerError, nil, errors.New("failed to create user entity"))
			return
		}

		users = append(users, user)
	}

	if err := h.userService.CreateUsers(r.Context(), users); err != nil {
		if err.Error() == "username already exists in database" {
			response.JSONer(w, http.StatusConflict, nil, err)
		} else if err.Error() == "duplicate username in request" {
			response.JSONer(w, http.StatusBadRequest, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	createdUsers := make([]UserShort, 0, len(users))
	for _, user := range users {
		createdUsers = append(createdUsers, UserShort{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	responseData := CreateUsersResponse{
		Message: "users created successfully",
		Count:   len(users),
		Users:   createdUsers,
	}

	response.JSONer(w, http.StatusCreated, responseData, nil)
}
