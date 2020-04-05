package circleci

import (
	"context"
)

// User defines struct for item of api response of user.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

// userResponse defines struct of user response.
type userResponse User

// GetMyself get current user information.
func (c *Client) GetMyself(ctx context.Context) (User, error) {
	var res userResponse
	if err := c.get(&res, APIPath[c.apiVersion]["me"], withContext(ctx)); err != nil {
		return User(res), err
	}
	return User(res), nil
}
