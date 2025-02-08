package response

type UserResponse struct {
	ID             string       `json:"id"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	ProfilePicture string       `json:"profile_picture"`
	Email          string       `json:"email"`
	Role           roleResponse `json:"role"`
}
