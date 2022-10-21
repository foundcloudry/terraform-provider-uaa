package api

import (
	"fmt"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
)

type ClientManager struct {
	log *Logger
	api *UaaApi
}

type UAAClient struct {
	ClientID             string   `json:"client_id,omitempty"`
	ClientSecret         string   `json:"client_secret,omitempty"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types,omitempty"`
	RedirectURI          []string `json:"redirect_uri,omitempty"`
	Scope                []string `json:"scope,omitempty"`
	ResourceIds          []string `json:"resource_ids,omitempty"`
	Authorities          []string `json:"authorities,omitempty"`
	AutoApprove          []string `json:"autoapprove,omitempty"`
	AccessTokenValidity  int      `json:"access_token_validity,omitempty"`
	RefreshTokenValidity int      `json:"refresh_token_validity,omitempty"`
	AllowedProviders     []string `json:"allowedproviders,omitempty"`
	Name                 string   `json:"name,omitempty"`
	TokenSalt            string   `json:"token_salt,omitempty"`
	CreatedWith          string   `json:"createdwith,omitempty"`
	ApprovalsDeleted     bool     `json:"approvals_deleted,omitempty"`
	RequiredUserGroups   []string `json:"required_user_groups,omitempty"`
	LastModified         int64    `json:"lastModified,omitempty"`
}

// UAAClientResourceList -
type UAAClientResourceList struct {
	Resources []UAAClient `json:"resources"`
}

func (c *UAAClient) HasDefaultScope() bool {
	return len(c.Scope) == 1 && c.Scope[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultAuthorites() bool {
	return len(c.Authorities) == 1 && c.Authorities[0] == "uaa.none"
}

func (c *UAAClient) HasDefaultResourceIds() bool {
	return len(c.ResourceIds) == 1 && c.ResourceIds[0] == "none"
}

func newClientManager(config coreconfig.Reader, gateway net.Gateway, logger *Logger) (cm *ClientManager, err error) {

	api, err := newUaaApi(config, gateway)
	if err != nil {
		return
	}

	cm = &ClientManager{
		log: logger,
		api: api,
	}
	return
}

func (manager *ClientManager) GetClient(id string) (client *UAAClient, err error) {

	path := fmt.Sprintf("/oauth/clients/%s", id)
	client = &UAAClient{}
	err = manager.api.Get(path, &client)
	return
}

func (manager *ClientManager) Create(newClient UAAClient) (client UAAClient, err error) {

	err = manager.api.Post("/oauth/clients", newClient, &client)
	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("client", newClient.ClientID)
		}
	}
	return
}

func (manager *ClientManager) UpdateClient(updatedClient *UAAClient) (client UAAClient, err error) {

	path := fmt.Sprintf("/oauth/clients/%s", updatedClient.ClientID)
	if err := manager.api.Put(path, updatedClient, &client); err != nil {
		return client, err
	}

	return
}

func (manager *ClientManager) DeleteClient(id string) error {

	return manager.api.Delete(fmt.Sprintf("/oauth/clients/%s", id))
}

func (manager *ClientManager) ChangeSecret(id, oldSecret, newSecret string) (err error) {

	data := map[string]string{
		"secret": newSecret,
	}

	if len(oldSecret) != 0 {
		data["oldSecret"] = oldSecret
	}

	path := fmt.Sprintf("/oauth/clients/%s/secret", id)
	response := make(map[string]interface{})

	return manager.api.Put(path, data, &response)
}

func (manager *ClientManager) FindByClientID(clientID string) (client UAAClient, err error) {

	filter := url.QueryEscape(fmt.Sprintf(`client_id Eq "%s"`, clientID))
	path := fmt.Sprintf("/oauth/clients?filter=%s", filter)

	clientResourceList := &UAAClientResourceList{}
	err = manager.api.Get(path, clientResourceList)

	if err == nil {
		if len(clientResourceList.Resources) > 0 {
			client = clientResourceList.Resources[0]
		} else {
			err = errors.NewModelNotFoundError("Client", clientID)
		}
	}
	return
}
