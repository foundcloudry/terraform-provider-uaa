package api

import (
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
	"fmt"
	apiheaders "github.com/jlpospisil/terraform-provider-uaa/uaa/api/headers"
	"net/http"
	"net/url"
)

type GroupManager struct {
	log *Logger
	api *UaaApi
}

type UAAGroup struct {
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	ZoneId      string `json:"zoneId,omitempty"`
}

type UAAGroupResourceList struct {
	Resources []UAAGroup `json:"resources"`
}

func newGroupManager(config coreconfig.Reader, gateway net.Gateway, logger *Logger) (gm *GroupManager, err error) {

	api, err := newUaaApi(config, gateway)
	if err != nil {
		return
	}

	gm = &GroupManager{
		log: logger,
		api: api,
	}
	return
}

func (manager *GroupManager) CreateGroup(displayName string, description string, zoneId string) (group *UAAGroup, err error) {

	groupResource := UAAGroup{
		DisplayName: displayName,
		Description: description,
		ZoneId:      zoneId,
	}

	err = manager.api.
		WithZoneId(zoneId).
		Post("/Groups", groupResource, &group)
	if err != nil {
		return nil, err
	}

	switch httpErr := err.(type) {
	case errors.HTTPError:
		if httpErr.StatusCode() == http.StatusConflict {
			err = errors.NewModelAlreadyExistsError("group", displayName)
		}
	}

	return
}

func (manager *GroupManager) GetGroup(id, zoneId string) (group *UAAGroup, err error) {

	path := fmt.Sprintf("/Groups/%s", id)
	err = manager.api.
		WithZoneId(zoneId).
		Get(path, &group)

	return
}

func (manager *GroupManager) UpdateGroup(id, displayName, description, zoneId string) (group *UAAGroup, err error) {

	groupResource := UAAGroup{
		DisplayName: displayName,
		Description: description,
	}

	path := fmt.Sprintf("/Groups/%s", id)
	err = manager.api.
		WithZoneId(zoneId).
		WithHeaders(map[string]string{
			apiheaders.IfMatch.String(): "*",
		}).
		Put(path, groupResource, &group)

	return
}

func (manager *GroupManager) DeleteGroup(id, zoneId string) error {

	path := fmt.Sprintf("/Groups/%s", id)
	return manager.api.WithZoneId(zoneId).Delete(path)
}

func (manager *GroupManager) FindByDisplayName(displayName, zoneId string) (group *UAAGroup, err error) {

	displayNameFilter := url.QueryEscape(fmt.Sprintf(`displayName Eq "%s"`, displayName))
	path := fmt.Sprintf("/Groups?filter=%s", displayNameFilter)

	groupResourceList := &UAAGroupResourceList{}
	err = manager.api.
		WithZoneId(zoneId).
		Get(path, groupResourceList)

	if err == nil {
		if len(groupResourceList.Resources) > 0 {
			group = &groupResourceList.Resources[0]
		} else {
			err = errors.NewModelNotFoundError("Group", displayName)
		}
	}
	return
}
