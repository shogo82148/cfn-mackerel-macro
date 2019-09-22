package cfn

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

func TestCreateUser(t *testing.T) {
	var invited bool
	u := &user{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createInvitation: func(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error) {
					if email != "john.doe@example.com" {
						t.Errorf("unexpected invite email: want %s, got %s", "john.doe@example.com", email)
					}
					if authority != "viewer" {
						t.Errorf("unexpected authority, want %s, got %s", "viewer", authority)
					}
					invited = true
					return &mackerel.Invitation{
						Email:     email,
						Authority: authority,
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:User",
			LogicalResourceID: "User",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Email": "john.doe@example.com",
			},
		},
	}
	id, param, err := u.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !invited {
		t.Error("not invited")
	}
	if id != "mkr:test-org:user:john.doe@example.com" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:user:john.doe@example.com", id)
	}
	if param["Email"].(string) != "john.doe@example.com" {
		t.Errorf("unexpected email, want %s, got %s", "john.doe@example.com", param["Email"].(string))
	}
}

func TestCreateUser_alreadyInvited(t *testing.T) {
	var invited bool
	u := &user{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createInvitation: func(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error) {
					if email != "john.doe@example.com" {
						t.Errorf("unexpected invite email: want %s, got %s", "john.doe@example.com", email)
					}
					if authority != "viewer" {
						t.Errorf("unexpected authority, want %s, got %s", "viewer", authority)
					}
					invited = true
					return nil, mkrError{
						statusCode: http.StatusBadRequest,
					}
				},
				findInvitations: func(ctx context.Context) ([]*mackerel.Invitation, error) {
					return []*mackerel.Invitation{
						{
							Email:     "john.doe@example.com",
							Authority: "viewer",
						},
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:User",
			LogicalResourceID: "User",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Email": "john.doe@example.com",
			},
		},
	}
	id, param, err := u.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !invited {
		t.Error("not invited")
	}
	if id != "mkr:test-org:user:john.doe@example.com" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:user:john.doe@example.com", id)
	}
	if param["Email"].(string) != "john.doe@example.com" {
		t.Errorf("unexpected email, want %s, got %s", "john.doe@example.com", param["Email"].(string))
	}
}

func TestCreateUser_alreadyExists(t *testing.T) {
	var invited bool
	u := &user{
		Function: &Function{
			org: &mackerel.Org{
				Name: "test-org",
			},
			client: &fakeMackerelClient{
				createInvitation: func(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error) {
					if email != "john.doe@example.com" {
						t.Errorf("unexpected invite email: want %s, got %s", "john.doe@example.com", email)
					}
					if authority != "viewer" {
						t.Errorf("unexpected authority, want %s, got %s", "viewer", authority)
					}
					invited = true
					return nil, mkrError{
						statusCode: http.StatusBadRequest,
					}
				},
				findInvitations: func(ctx context.Context) ([]*mackerel.Invitation, error) {
					return []*mackerel.Invitation{}, nil
				},
				findUsers: func(ctx context.Context) ([]*mackerel.User, error) {
					return []*mackerel.User{
						{
							ID:        "johndoe",
							Email:     "john.doe@example.com",
							Authority: "viewer",
						},
					}, nil
				},
			},
		},
		Event: cfn.Event{
			RequestType:       cfn.RequestCreate,
			RequestID:         "",
			ResponseURL:       "https://cloudformation-custom-resource-response-apnortheast1.s3-ap-northeast-1.amazonaws.com/xxxxx",
			ResourceType:      "Custom:User",
			LogicalResourceID: "User",
			StackID:           "arn:aws:cloudformation:ap-northeast-1:1234567890:stack/foobar/12345678-1234-1234-1234-123456789abc",
			ResourceProperties: map[string]interface{}{
				"Email": "john.doe@example.com",
			},
		},
	}
	id, param, err := u.create(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !invited {
		t.Error("not invited")
	}
	if id != "mkr:test-org:user:john.doe@example.com" {
		t.Errorf("unexpected host id: want %s, got %s", "mkr:test-org:user:john.doe@example.com", id)
	}
	if param["Email"].(string) != "john.doe@example.com" {
		t.Errorf("unexpected email, want %s, got %s", "john.doe@example.com", param["Email"].(string))
	}
}
