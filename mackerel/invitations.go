package mackerel

import (
	"context"
	"errors"
	"net/http"
)

// Invitation is a user invitation.
type Invitation struct {
	Email     string        `json:"email"`
	Authority UserAuthority `json:"authority"`
	ExpiresAt int64         `json:"expiresAt"`
}

// FindInvitations returns a list of user invitations.
func (c *Client) FindInvitations(ctx context.Context) ([]*Invitation, error) {
	var data struct {
		Invitations []*Invitation `json:"invitations"`
	}
	_, err := c.do(ctx, http.MethodGet, "/api/v0/invitations", nil, &data)
	if err != nil {
		return nil, err
	}
	return data.Invitations, nil
}

// CreateInvitation invites a user.
func (c *Client) CreateInvitation(ctx context.Context, email string, authority UserAuthority) (*Invitation, error) {
	param := struct {
		Email     string        `json:"email"`
		Authority UserAuthority `json:"authority"`
	}{
		Email:     email,
		Authority: authority,
	}
	var ret Invitation
	_, err := c.do(ctx, http.MethodPost, "/api/v0/invitations", param, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// RevokeInvitation revokes a user invitation.
func (c *Client) RevokeInvitation(ctx context.Context, email string) error {
	param := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}
	var data struct {
		Success bool `json:"success"`
	}
	_, err := c.do(ctx, http.MethodPost, "/api/v0/invitations", param, &data)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New("unexpected response")
	}
	return nil
}
