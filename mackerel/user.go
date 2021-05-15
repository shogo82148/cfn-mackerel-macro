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
	JoinedAt                Timestamp                  `json:"joinedAt"`
}

// UserAuthority is the authority type for user.
type UserAuthority string

const (
	// UserAuthorityOwner is the owner authority type
	UserAuthorityOwner UserAuthority = "owner"

	// UserAuthorityManager is the manager authority type
	UserAuthorityManager UserAuthority = "manager"

	// UserAuthorityCollaborator is the collaborator authority type
	UserAuthorityCollaborator UserAuthority = "collaborator"

	// UserAuthorityViewer is the viewer authority type
	UserAuthorityViewer UserAuthority = "viewer"
)

func (t UserAuthority) String() string {
	return string(t)
}

// UserAuthenticationMethod is a method of authentication
type UserAuthenticationMethod string

const (
	// UserAuthenticationMethodPassword is the password authentication.
	UserAuthenticationMethodPassword UserAuthenticationMethod = "password"

	// UserAuthenticationMethodGitHub is the GitHub authentication.
	UserAuthenticationMethodGitHub UserAuthenticationMethod = "github"

	// UserAuthenticationMethodIDCF is the IDCF authentication.
	UserAuthenticationMethodIDCF UserAuthenticationMethod = "idcf"

	// UserAuthenticationMethodGoogle is the Google authentication.
	UserAuthenticationMethodGoogle UserAuthenticationMethod = "google"

	// UserAuthenticationMethodNifty is the nifty authentication.
	UserAuthenticationMethodNifty UserAuthenticationMethod = "nifty"

	// UserAuthenticationMethodYammer is the Yammer authentication.
	UserAuthenticationMethodYammer UserAuthenticationMethod = "yammer"

	// UserAuthenticationMethodKDDI is the KDDI authentication.
	UserAuthenticationMethodKDDI UserAuthenticationMethod = "kddi"
)

func (t UserAuthenticationMethod) String() string {
	return string(t)
}

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
