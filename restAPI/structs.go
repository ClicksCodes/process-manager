package restAPI

type containerTemplate struct {
	Packages []string `json:"packages"` // List of NixOS packages to be installed in the container
} // Template for a container

type container struct {
	ID       int               `json:"id"`       // Container ID
	Template containerTemplate `json:"template"` // Container template
} // Container

type containerList struct {
	Containers []container `json:"containers"` // List of containers
} // List of containers

type containerCreate struct {
	Template containerTemplate `json:"template"` // Container template
} // Container creation request

type containerCreateResponse struct {
	ID int `json:"id"` // Container ID
} // Container creation response

type containerDelete struct {
	ID int `json:"id"` // Container ID
} // Container deletion request

type containerDeleteResponse struct {
	Success bool `json:"success"` // True if the container was successfully deleted
} // Container deletion response

type containerStart struct {
	ID int `json:"id"` // Container ID
} // Container start request

type containerStartResponse struct {
	Success bool `json:"success"` // True if the container was successfully started
} // Container start response

type containerStop struct {
	ID int `json:"id"` // Container ID
} // Container stop request

type containerStopResponse struct {
	Success bool `json:"success"` // True if the container was successfully stopped
} // Container stop response

type containerStatus struct {
	ID int `json:"id"` // Container ID
} // Container status request

type containerStatusResponse struct {
	Status string `json:"status"` // Status of the container
} // Container status response
