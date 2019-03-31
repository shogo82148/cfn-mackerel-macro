package cfn

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/dproxy"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type user struct {
	Function *Function
	Event    cfn.Event
}

func (u *user) create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(u.Event.ResourceProperties)

	email := d.String(in.M("Email"))
	authority := d.String(dproxy.Default(in.M("Authority"), "viewer"))
	if err := d.CombineErrors(); err != nil {
		return "", nil, err
	}

	id, err := u.Function.buildUserID(ctx, email)
	if err != nil {
		return "", nil, err
	}

	// try to invite the user.
	c := u.Function.getclient()
	_, err = c.CreateInvitation(ctx, email, mackerel.UserAuthority(authority))
	if merr, ok := err.(mackerel.Error); ok {
		if merr.StatusCode() != http.StatusBadRequest {
			return "", nil, err
		}
		// already invited?
		invited, err := u.alreadyInvited(ctx, email)
		if err != nil {
			return "", nil, err
		}
		if !invited {
			// already in the org?
			uid, err := u.getUserID(ctx, email)
			if err != nil {
				return "", nil, err
			}
			if uid == "" {
				return "", nil, fmt.Errorf("fail to invite %s", email)
			}
		}
	}

	return id, map[string]interface{}{
		"Email": email,
	}, nil
}

func (u *user) alreadyInvited(ctx context.Context, email string) (bool, error) {
	c := u.Function.getclient()
	list, err := c.FindInvitations(ctx)
	if err != nil {
		return false, err
	}
	for _, invite := range list {
		if invite.Email == email {
			return true, nil
		}
	}
	return false, nil
}

// get the user id from email.
// return empty string if the user is not in the org.
func (u *user) getUserID(ctx context.Context, email string) (string, error) {
	c := u.Function.getclient()
	users, err := c.FindUsers(ctx)
	if err != nil {
		return "", err
	}
	for _, user := range users {
		if user.Email == email {
			return user.ID, nil
		}
	}
	return "", nil
}

func (u *user) update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	var d dproxy.Drain
	in := dproxy.New(u.Event.ResourceProperties)
	old := dproxy.New(u.Event.OldResourceProperties)

	email := d.String(in.M("Email"))
	oldEmail := d.String(old.M("Email"))
	if err := d.CombineErrors(); err != nil {
		return u.Event.PhysicalResourceID, nil, err
	}

	if email == oldEmail {
		// no need to update.
		// updating authority is not supported.
		return u.Event.PhysicalResourceID, map[string]interface{}{
			"Email": email,
		}, nil
	}

	// create a new invitation
	return u.create(ctx)
}

func (u *user) delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error) {
	physicalResourceID = u.Event.PhysicalResourceID
	email, err := u.Function.parseUserID(ctx, physicalResourceID)
	if err != nil {
		return
	}
	data = map[string]interface{}{
		"Email": email,
	}

	// revoke invitation
	c := u.Function.getclient()
	err = c.RevokeInvitation(ctx, email)
	if merr, ok := err.(mackerel.Error); ok {
		if merr.StatusCode() != http.StatusNotFound {
			return
		}
		// maybe already accept the invitation
	}

	// delete the user
	uid, err := u.getUserID(ctx, email)
	if err != nil {
		return
	}
	if uid == "" {
		// the user is already deleted.
		return
	}
	_, err = c.DeleteUser(ctx, uid)
	if merr, ok := err.(mackerel.Error); ok {
		if merr.StatusCode() != http.StatusNotFound {
			return
		}
		// the user is already deleted.
	}
	err = nil
	return
}
