package mackerel

import (
	"context"
	"fmt"
	"net/http"
)

// User is a user in mackerel.
type User struct {
	ID                      string                     `json:"id"`
	ScreenName              string                     `json:"screenName"`
	Email                   string                     `json:"email"`
	Authority               UserAuthority              `json:"authority"`
	IsInRegistrationProcess bool                       `json:"isInRegistrationProcess"`
	IsMFAEnabled            bool                       `json:"isMFAEnabled"`
	AuthenticationMethods   []UserAuthenticationMethod `json:"authenticationMethods"`
	JoinedAt                int64                      `json:"joinedAt"`
}

// UserAuthority is authority type for user.
type UserAuthority string

// UserAuthenticationMethod is a method of authentication
type UserAuthenticationMethod string

// FindUsers returns a list of users.
func (c *Client) FindUsers(ctx context.Context) ([]*User, error) {
	var users struct {
		Users []*User `json:"users"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/users", nil, &users)
	if err != nil {
		return nil, err
	}
	return users.Users, nil
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(ctx context.Context, userID string) (*User, error) {
	var user User
	_, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v0/users/%s", userID), nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
