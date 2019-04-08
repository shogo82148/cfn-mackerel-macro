package cfn

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

// Function is a custom resource function for CloudForamtion.
type Function struct {
	APIKey         string
	APIKeyProvider mackerel.APIKeyProvider
	BaseURL        *url.URL

	mu     sync.Mutex
	client makerelInterface
	org    *mackerel.Org
}

type makerelInterface interface {
	// org
	GetOrg(ctx context.Context) (*mackerel.Org, error)

	// host
	CreateHost(ctx context.Context, param *mackerel.CreateHostParam) (string, error)
	UpdateHost(ctx context.Context, hostID string, param *mackerel.UpdateHostParam) (string, error)
	RetireHost(ctx context.Context, id string) error

	// host metadata
	GetHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) (*mackerel.HostMetaMetaData, error)
	GetHostMetaDataNameSpaces(ctx context.Context, hostID string) ([]string, error)
	PutHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) error
	DeleteHostMetaData(ctx context.Context, hostID, namespace string) error

	// monitor
	CreateMonitor(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error)
	UpdateMonitor(ctx context.Context, monitorID string, param mackerel.Monitor) (mackerel.Monitor, error)
	DeleteMonitor(ctx context.Context, monitorID string) (mackerel.Monitor, error)

	// dashboard
	FindDashboards(ctx context.Context) ([]*mackerel.Dashboard, error)
	FindDashboard(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error)
	CreateDashboard(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	UpdateDashboard(ctx context.Context, dashboardID string, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	DeleteDashboard(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error)

	// role
	CreateRole(ctx context.Context, serviceName string, param *mackerel.CreateRoleParam) (*mackerel.Role, error)
	DeleteRole(ctx context.Context, serviceName, roleName string) (*mackerel.Role, error)

	// role metadata
	GetRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) (*mackerel.RoleMetaMetaData, error)
	GetRoleMetaDataNameSpaces(ctx context.Context, serviceName, roleName string) ([]string, error)
	PutRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) error
	DeleteRoleMetaData(ctx context.Context, serviceName, roleName, namespace string) error

	// service
	CreateService(ctx context.Context, param *mackerel.CreateServiceParam) (*mackerel.Service, error)
	DeleteService(ctx context.Context, serviceName string) (*mackerel.Service, error)

	// service metadata
	GetServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) (*mackerel.ServiceMetaMetaData, error)
	GetServiceMetaDataNameSpaces(ctx context.Context, serviceName string) ([]string, error)
	PutServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) error
	DeleteServiceMetaData(ctx context.Context, serviceName, namespace string) error

	// notification channels
	FindNotificationChannels(ctx context.Context) ([]mackerel.NotificationChannel, error)
	CreateNotificationChannel(ctx context.Context, ch mackerel.NotificationChannel) (mackerel.NotificationChannel, error)
	DeleteNotificationChannel(ctx context.Context, channelID string) (mackerel.NotificationChannel, error)

	// notification group
	FindNotificationGroups(ctx context.Context) ([]*mackerel.NotificationGroup, error)
	CreateNotificationGroup(ctx context.Context, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
	UpdateNotificationGroup(ctx context.Context, groupID string, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
	DeleteNotificationGroup(ctx context.Context, groupID string) (*mackerel.NotificationGroup, error)

	// user
	FindUsers(ctx context.Context) ([]*mackerel.User, error)
	DeleteUser(ctx context.Context, userID string) (*mackerel.User, error)

	// invitation
	FindInvitations(ctx context.Context) ([]*mackerel.Invitation, error)
	CreateInvitation(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error)
	RevokeInvitation(ctx context.Context, email string) error
}

type resource interface {
	create(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
	update(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
	delete(ctx context.Context) (physicalResourceID string, data map[string]interface{}, err error)
}

// Handle handles custom resource events of CloudForamtion.
func (f *Function) Handle(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	if strings.HasPrefix(event.PhysicalResourceID, "mkr::error:") {
		// it is dummy resource, just ignore it
		return event.PhysicalResourceID, nil, nil
	}

	typ := strings.TrimPrefix(event.ResourceType, "Custom::")
	var r resource
	switch typ {
	case "Service":
		r = &service{
			Function: f,
			Event:    event,
		}
	case "Role":
		r = &role{
			Function: f,
			Event:    event,
		}
	case "Host":
		r = &host{
			Function: f,
			Event:    event,
		}
	case "Monitor":
		r = &monitor{
			Function: f,
			Event:    event,
		}
	case "Dashboard":
		r = &dashboard{
			Function: f,
			Event:    event,
		}
	case "NotificationChannel":
		r = &notificationChannel{
			Function: f,
			Event:    event,
		}
	case "NotificationGroup":
		r = &notificationGroup{
			Function: f,
			Event:    event,
		}
	case "User":
		r = &user{
			Function: f,
			Event:    event,
		}
	default:
		return "", nil, nil // fmt.Errorf("unknown type: %s", typ)
	}
	switch event.RequestType {
	case cfn.RequestCreate:
		physicalResourceID, data, err = r.create(ctx)
	case cfn.RequestUpdate:
		physicalResourceID, data, err = r.update(ctx)
	case cfn.RequestDelete:
		physicalResourceID, data, err = r.delete(ctx)
	default:
		err = fmt.Errorf("unknown request type: %s", event.RequestType)
	}
	if physicalResourceID == "" {
		// physicalResourceID must not empty.
		// return dummy resource id.
		physicalResourceID = "mkr::error:" + event.RequestID
	}
	return
}

// LambdaWrap returns a CustomResourceLambdaFunction which is something lambda.Start()
// will understand.
func (f *Function) LambdaWrap() cfn.CustomResourceLambdaFunction {
	return cfn.LambdaWrap(f.Handle)
}

func (f *Function) getclient() makerelInterface {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.client == nil {
		f.client = &mackerel.Client{
			APIKey:         f.APIKey,
			APIKeyProvider: f.APIKeyProvider,
			BaseURL:        f.BaseURL,
		}
	}
	return f.client
}

func (f *Function) getorg(ctx context.Context) (*mackerel.Org, error) {
	c := f.getclient()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.org == nil {
		org, err := c.GetOrg(ctx)
		if err != nil {
			return nil, err
		}
		f.org = org
	}
	return f.org, nil
}

func (f *Function) buildID(ctx context.Context, resourceType string, ids ...string) (string, error) {
	org, err := f.getorg(ctx)
	if err != nil {
		return "", err
	}

	ret := append([]string{"mkr", org.Name, resourceType}, ids...)
	return strings.Join(ret, ":"), nil
}

func (f *Function) buildServiceID(ctx context.Context, serviceName string) (string, error) {
	return f.buildID(ctx, "service", serviceName)
}

func (f *Function) buildRoleID(ctx context.Context, serviceName, roleName string) (string, error) {
	return f.buildID(ctx, "role", serviceName, roleName)
}

func (f *Function) buildHostID(ctx context.Context, hostID string) (string, error) {
	return f.buildID(ctx, "host", hostID)
}

func (f *Function) buildMonitorID(ctx context.Context, monitorID string) (string, error) {
	return f.buildID(ctx, "monitor", monitorID)
}

func (f *Function) buildDashboardID(ctx context.Context, dashboardID string) (string, error) {
	return f.buildID(ctx, "dashboard", dashboardID)
}

func (f *Function) buildNotificationChannelID(ctx context.Context, channelID string) (string, error) {
	return f.buildID(ctx, "notification-channel", channelID)
}

func (f *Function) buildNotificationGroupID(ctx context.Context, groupID string) (string, error) {
	return f.buildID(ctx, "notification-group", groupID)
}

func (f *Function) buildUserID(ctx context.Context, email string) (string, error) {
	return f.buildID(ctx, "user", email)
}

// parseID parses ID of Mackerel resources.
func (f *Function) parseID(ctx context.Context, id string, n int) (string, []string, error) {
	org, err := f.getorg(ctx)
	if err != nil {
		return "", nil, err
	}

	ids := strings.Split(id, ":")
	if len(ids) < n+3 {
		return "", nil, fmt.Errorf("invalid mkr id: %s", id)
	}
	if ids[0] != "mkr" {
		return "", nil, fmt.Errorf("invalid mkr id: %s", id)
	}
	if ids[1] != org.Name {
		return "", nil, fmt.Errorf("invalid org name in id: %s", id)
	}
	return ids[2], ids[3:], nil
}

func (f *Function) parseServiceID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "service" {
		return "", fmt.Errorf("invalid type %s, expected service type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseRoleID(ctx context.Context, id string) (string, string, error) {
	typ, parts, err := f.parseID(ctx, id, 2)
	if err != nil {
		return "", "", err
	}
	if typ != "role" {
		return "", "", fmt.Errorf("invalid type %s, expected role type", typ)
	}
	return parts[0], parts[1], nil
}

func (f *Function) parseHostID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "host" {
		return "", fmt.Errorf("invalid type %s, expected host type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseMonitorID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "monitor" {
		return "", fmt.Errorf("invalid type %s, expected monitor type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseDashboardID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "dashboard" {
		return "", fmt.Errorf("invalid type %s, expected monitor type", typ)
	}
	return parts[0], nil
}

func (f *Function) parseNotificationChannelID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "notification-channel" {
		return "", fmt.Errorf("invalid type %s, expected notification channel", typ)
	}
	return parts[0], nil
}

func (f *Function) parseNotificationGroupID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "notification-group" {
		return "", fmt.Errorf("invalid type %s, expected notification group", typ)
	}
	return parts[0], nil
}

func (f *Function) parseUserID(ctx context.Context, id string) (string, error) {
	typ, parts, err := f.parseID(ctx, id, 1)
	if err != nil {
		return "", err
	}
	if typ != "user" {
		return "", fmt.Errorf("invalid type %s, expected user", typ)
	}
	return parts[0], nil
}

type metadata struct {
	StackName string `json:"stack_name"`
	StackID   string `json:"stack_id"`
	LogicalID string `json:"logical_id"`
}

func getmetadata(e cfn.Event) metadata {
	// arn format: arn:aws:cloudformation:${AWS_REGION}:${AWS::ACCOUNT}:stack/${STACK_NAME}/${UUID}
	name := e.StackID
	if idx := strings.LastIndexByte(name, ':'); idx >= 0 {
		name = strings.TrimPrefix(name[idx:], ":stack/")
	}
	if idx := strings.LastIndexByte(name, '/'); idx >= 0 {
		name = name[:idx]
	}
	return metadata{
		StackName: name,
		StackID:   e.StackID,
		LogicalID: e.LogicalResourceID,
	}
}
