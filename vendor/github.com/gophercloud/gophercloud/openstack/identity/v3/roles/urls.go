package roles

import "github.com/gophercloud/gophercloud"

const (
	rolePath = "roles"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func getURL(client *gophercloud.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func listAssignmentsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("role_assignments")
}

func assignURL(client *gophercloud.ServiceClient, targetType, targetID, actorType, actorID, roleID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath, roleID)
}
