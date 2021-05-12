package cfn

import (
	"context"
	"fmt"

	"github.com/shogo82148/cfn-mackerel-macro/mackerel"
)

type fakeMackerelClient struct {
	getOrg                               func(ctx context.Context) (*mackerel.Org, error)
	createHost                           func(ctx context.Context, param *mackerel.CreateHostParam) (string, error)
	updateHost                           func(ctx context.Context, hostID string, param *mackerel.UpdateHostParam) (string, error)
	retireHost                           func(ctx context.Context, id string) error
	getHostMetaData                      func(ctx context.Context, hostID, namespace string, v interface{}) (*mackerel.HostMetaMetaData, error)
	getHostMetaDataNameSpaces            func(ctx context.Context, hostID string) ([]string, error)
	putHostMetaData                      func(ctx context.Context, hostID, namespace string, v interface{}) error
	deleteHostMetaData                   func(ctx context.Context, hostID, namespace string) error
	createMonitor                        func(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error)
	updateMonitor                        func(ctx context.Context, monitorID string, param mackerel.Monitor) (mackerel.Monitor, error)
	deleteMonitor                        func(ctx context.Context, monitorID string) (mackerel.Monitor, error)
	findDashboards                       func(ctx context.Context) ([]*mackerel.Dashboard, error)
	findDashboard                        func(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error)
	createDashboard                      func(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	updateDashboard                      func(ctx context.Context, dashboardID string, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	deleteDashboard                      func(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error)
	createRole                           func(ctx context.Context, serviceName string, param *mackerel.CreateRoleParam) (*mackerel.Role, error)
	deleteRole                           func(ctx context.Context, serviceName, roleName string) (*mackerel.Role, error)
	getRoleMetaData                      func(ctx context.Context, serviceName, roleName, namespace string, v interface{}) (*mackerel.RoleMetaMetaData, error)
	getRoleMetaDataNameSpaces            func(ctx context.Context, serviceName, roleName string) ([]string, error)
	putRoleMetaData                      func(ctx context.Context, serviceName, roleName, namespace string, v interface{}) error
	deleteRoleMetaData                   func(ctx context.Context, serviceName, roleName, namespace string) error
	createService                        func(ctx context.Context, param *mackerel.CreateServiceParam) (*mackerel.Service, error)
	deleteService                        func(ctx context.Context, serviceName string) (*mackerel.Service, error)
	getServiceMetaData                   func(ctx context.Context, serviceName, namespace string, v interface{}) (*mackerel.ServiceMetaMetaData, error)
	getServiceMetaDataNameSpaces         func(ctx context.Context, serviceName string) ([]string, error)
	putServiceMetaData                   func(ctx context.Context, serviceName, namespace string, v interface{}) error
	deleteServiceMetaData                func(ctx context.Context, serviceName, namespace string) error
	findNotificationChannels             func(ctx context.Context) ([]mackerel.NotificationChannel, error)
	createNotificationChannel            func(ctx context.Context, ch mackerel.NotificationChannel) (mackerel.NotificationChannel, error)
	deleteNotificationChannel            func(ctx context.Context, channelID string) (mackerel.NotificationChannel, error)
	findNotificationGroups               func(ctx context.Context) ([]*mackerel.NotificationGroup, error)
	createNotificationGroup              func(ctx context.Context, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
	updateNotificationGroup              func(ctx context.Context, groupID string, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error)
	deleteNotificationGroup              func(ctx context.Context, groupID string) (*mackerel.NotificationGroup, error)
	findUsers                            func(ctx context.Context) ([]*mackerel.User, error)
	deleteUser                           func(ctx context.Context, userID string) (*mackerel.User, error)
	findInvitations                      func(ctx context.Context) ([]*mackerel.Invitation, error)
	createInvitation                     func(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error)
	revokeInvitation                     func(ctx context.Context, email string) error
	findDowntimes                        func(ctx context.Context) ([]*mackerel.Downtime, error)
	createDowntime                       func(ctx context.Context, param *mackerel.Downtime) (*mackerel.Downtime, error)
	updateDowntime                       func(ctx context.Context, downtimeID string, param *mackerel.Downtime) (*mackerel.Downtime, error)
	deleteDowntime                       func(ctx context.Context, downtimeID string) (*mackerel.Downtime, error)
	findAWSIntegrations                  func(ctx context.Context) ([]*mackerel.AWSIntegration, error)
	findAWSIntegration                   func(ctx context.Context, awsIntegrationID string) (*mackerel.AWSIntegration, error)
	createAWSIntegration                 func(ctx context.Context, param *mackerel.AWSIntegration) (*mackerel.AWSIntegration, error)
	updateAWSIntegration                 func(ctx context.Context, awsIntegrationID string, param *mackerel.AWSIntegration) (*mackerel.AWSIntegration, error)
	deleteAWSIntegration                 func(ctx context.Context, awsIntegrationID string) (*mackerel.AWSIntegration, error)
	createAWSIntegrationExternalID       func(ctx context.Context) (string, error)
	findAWSIntegrationsExcludableMetrics func(ctx context.Context) (map[string][]string, error)
}

var _ makerelInterface = (*fakeMackerelClient)(nil)

func (c *fakeMackerelClient) GetOrg(ctx context.Context) (*mackerel.Org, error) {
	return c.getOrg(ctx)
}

func (c *fakeMackerelClient) CreateHost(ctx context.Context, param *mackerel.CreateHostParam) (string, error) {
	return c.createHost(ctx, param)
}

func (c *fakeMackerelClient) UpdateHost(ctx context.Context, hostID string, param *mackerel.UpdateHostParam) (string, error) {
	return c.updateHost(ctx, hostID, param)
}

func (c *fakeMackerelClient) RetireHost(ctx context.Context, id string) error {
	return c.retireHost(ctx, id)
}

func (c *fakeMackerelClient) GetHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) (*mackerel.HostMetaMetaData, error) {
	return c.getHostMetaData(ctx, hostID, namespace, v)
}

func (c *fakeMackerelClient) GetHostMetaDataNameSpaces(ctx context.Context, hostID string) ([]string, error) {
	return c.getHostMetaDataNameSpaces(ctx, hostID)
}

func (c *fakeMackerelClient) PutHostMetaData(ctx context.Context, hostID, namespace string, v interface{}) error {
	return c.putHostMetaData(ctx, hostID, namespace, v)
}

func (c *fakeMackerelClient) DeleteHostMetaData(ctx context.Context, hostID, namespace string) error {
	return c.deleteHostMetaData(ctx, hostID, namespace)
}

func (c *fakeMackerelClient) CreateMonitor(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
	return c.createMonitor(ctx, param)
}

func (c *fakeMackerelClient) UpdateMonitor(ctx context.Context, monitorID string, param mackerel.Monitor) (mackerel.Monitor, error) {
	return c.updateMonitor(ctx, monitorID, param)
}
func (c *fakeMackerelClient) DeleteMonitor(ctx context.Context, monitorID string) (mackerel.Monitor, error) {
	return c.deleteMonitor(ctx, monitorID)
}

func (c *fakeMackerelClient) FindDashboards(ctx context.Context) ([]*mackerel.Dashboard, error) {
	return c.findDashboards(ctx)
}

func (c *fakeMackerelClient) FindDashboard(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error) {
	return c.findDashboard(ctx, dashboardID)
}

func (c *fakeMackerelClient) CreateDashboard(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
	return c.createDashboard(ctx, param)
}

func (c *fakeMackerelClient) UpdateDashboard(ctx context.Context, dashboardID string, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
	return c.updateDashboard(ctx, dashboardID, param)
}

func (c *fakeMackerelClient) DeleteDashboard(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error) {
	return c.deleteDashboard(ctx, dashboardID)
}

func (c *fakeMackerelClient) CreateRole(ctx context.Context, serviceName string, param *mackerel.CreateRoleParam) (*mackerel.Role, error) {
	return c.createRole(ctx, serviceName, param)
}

func (c *fakeMackerelClient) DeleteRole(ctx context.Context, serviceName, roleName string) (*mackerel.Role, error) {
	return c.deleteRole(ctx, serviceName, roleName)
}

func (c *fakeMackerelClient) GetRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) (*mackerel.RoleMetaMetaData, error) {
	return c.getRoleMetaData(ctx, serviceName, roleName, namespace, v)
}

func (c *fakeMackerelClient) GetRoleMetaDataNameSpaces(ctx context.Context, serviceName, roleName string) ([]string, error) {
	return c.getRoleMetaDataNameSpaces(ctx, serviceName, roleName)
}

func (c *fakeMackerelClient) PutRoleMetaData(ctx context.Context, serviceName, roleName, namespace string, v interface{}) error {
	return c.putRoleMetaData(ctx, serviceName, roleName, namespace, v)
}

func (c *fakeMackerelClient) DeleteRoleMetaData(ctx context.Context, serviceName, roleName, namespace string) error {
	return c.deleteRoleMetaData(ctx, serviceName, roleName, namespace)
}

func (c *fakeMackerelClient) CreateService(ctx context.Context, param *mackerel.CreateServiceParam) (*mackerel.Service, error) {
	return c.createService(ctx, param)
}

func (c *fakeMackerelClient) DeleteService(ctx context.Context, serviceName string) (*mackerel.Service, error) {
	return c.deleteService(ctx, serviceName)
}

func (c *fakeMackerelClient) GetServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) (*mackerel.ServiceMetaMetaData, error) {
	return c.getServiceMetaData(ctx, serviceName, namespace, v)
}

func (c *fakeMackerelClient) GetServiceMetaDataNameSpaces(ctx context.Context, serviceName string) ([]string, error) {
	return c.getServiceMetaDataNameSpaces(ctx, serviceName)
}

func (c *fakeMackerelClient) PutServiceMetaData(ctx context.Context, serviceName, namespace string, v interface{}) error {
	return c.putServiceMetaData(ctx, serviceName, namespace, v)
}

func (c *fakeMackerelClient) DeleteServiceMetaData(ctx context.Context, serviceName, namespace string) error {
	return c.deleteServiceMetaData(ctx, serviceName, namespace)
}

func (c *fakeMackerelClient) FindNotificationChannels(ctx context.Context) ([]mackerel.NotificationChannel, error) {
	return c.findNotificationChannels(ctx)
}

func (c *fakeMackerelClient) CreateNotificationChannel(ctx context.Context, ch mackerel.NotificationChannel) (mackerel.NotificationChannel, error) {
	return c.createNotificationChannel(ctx, ch)
}

func (c *fakeMackerelClient) DeleteNotificationChannel(ctx context.Context, channelID string) (mackerel.NotificationChannel, error) {
	return c.deleteNotificationChannel(ctx, channelID)
}

func (c *fakeMackerelClient) FindNotificationGroups(ctx context.Context) ([]*mackerel.NotificationGroup, error) {
	return c.findNotificationGroups(ctx)
}

func (c *fakeMackerelClient) CreateNotificationGroup(ctx context.Context, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
	return c.createNotificationGroup(ctx, group)
}

func (c *fakeMackerelClient) UpdateNotificationGroup(ctx context.Context, groupID string, group *mackerel.NotificationGroup) (*mackerel.NotificationGroup, error) {
	return c.updateNotificationGroup(ctx, groupID, group)
}
func (c *fakeMackerelClient) DeleteNotificationGroup(ctx context.Context, groupID string) (*mackerel.NotificationGroup, error) {
	return c.deleteNotificationGroup(ctx, groupID)
}

func (c *fakeMackerelClient) FindUsers(ctx context.Context) ([]*mackerel.User, error) {
	return c.findUsers(ctx)
}

func (c *fakeMackerelClient) DeleteUser(ctx context.Context, userID string) (*mackerel.User, error) {
	return c.deleteUser(ctx, userID)
}

func (c *fakeMackerelClient) FindInvitations(ctx context.Context) ([]*mackerel.Invitation, error) {
	return c.findInvitations(ctx)
}

func (c *fakeMackerelClient) CreateInvitation(ctx context.Context, email string, authority mackerel.UserAuthority) (*mackerel.Invitation, error) {
	return c.createInvitation(ctx, email, authority)
}

func (c *fakeMackerelClient) RevokeInvitation(ctx context.Context, email string) error {
	return c.revokeInvitation(ctx, email)
}

func (c *fakeMackerelClient) FindDowntimes(ctx context.Context) ([]*mackerel.Downtime, error) {
	return c.findDowntimes(ctx)
}

func (c *fakeMackerelClient) CreateDowntime(ctx context.Context, param *mackerel.Downtime) (*mackerel.Downtime, error) {
	return c.createDowntime(ctx, param)
}

func (c *fakeMackerelClient) UpdateDowntime(ctx context.Context, downtimeID string, param *mackerel.Downtime) (*mackerel.Downtime, error) {
	return c.updateDowntime(ctx, downtimeID, param)
}

func (c *fakeMackerelClient) DeleteDowntime(ctx context.Context, downtimeID string) (*mackerel.Downtime, error) {
	return c.deleteDowntime(ctx, downtimeID)
}

func (c *fakeMackerelClient) FindAWSIntegrations(ctx context.Context) ([]*mackerel.AWSIntegration, error) {
	return c.findAWSIntegrations(ctx)
}

func (c *fakeMackerelClient) FindAWSIntegration(ctx context.Context, awsIntegrationID string) (*mackerel.AWSIntegration, error) {
	return c.findAWSIntegration(ctx, awsIntegrationID)
}

func (c *fakeMackerelClient) CreateAWSIntegration(ctx context.Context, param *mackerel.AWSIntegration) (*mackerel.AWSIntegration, error) {
	return c.createAWSIntegration(ctx, param)
}

func (c *fakeMackerelClient) UpdateAWSIntegration(ctx context.Context, awsIntegrationID string, param *mackerel.AWSIntegration) (*mackerel.AWSIntegration, error) {
	return c.updateAWSIntegration(ctx, awsIntegrationID, param)
}

func (c *fakeMackerelClient) DeleteAWSIntegration(ctx context.Context, awsIntegrationID string) (*mackerel.AWSIntegration, error) {
	return c.deleteAWSIntegration(ctx, awsIntegrationID)
}

func (c *fakeMackerelClient) CreateAWSIntegrationExternalID(ctx context.Context) (string, error) {
	return c.createAWSIntegrationExternalID(ctx)
}

func (c *fakeMackerelClient) FindAWSIntegrationsExcludableMetrics(ctx context.Context) (map[string][]string, error) {
	return c.findAWSIntegrationsExcludableMetrics(ctx)
}

type mkrError struct {
	statusCode int
	message    string
}

func (e mkrError) Error() string {
	return fmt.Sprintf("status: %d, %s", e.statusCode, e.message)
}

func (e mkrError) StatusCode() int {
	return e.statusCode
}

func (e mkrError) Message() string {
	return e.message
}
