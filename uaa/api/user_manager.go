package api

import (
	"code.cloudfoundry.org/cli/cf/net"
	"encoding/json"
	"fmt"
	apiheaders "github.com/foundcloudry/terraform-provider-uaa/uaa/api/headers"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
)

type UserManager struct {
	log                 *Logger
	api                 *UaaApi
	identityZoneManager *IdentityZoneManager
	clientToken         string
	groupMap            map[string]map[string]string
	defaultGroups       map[string]map[string]int
}

type UAAUser struct {
	Id       string         `json:"id,omitempty"`
	Username string         `json:"userName,omitempty"`
	Password string         `json:"password,omitempty"`
	Origin   string         `json:"origin,omitempty"`
	Name     UAAUserName    `json:"name,omitempty"`
	Emails   []UAAUserEmail `json:"emails,omitempty"`
	Groups   []UAAUserGroup `json:"groups,omitempty"`
	ZoneId   string         `json:"zoneId,omitempty"`
}

type UAAUserResourceList struct {
	Resources []UAAUser `json:"resources"`
}

type UAAUserEmail struct {
	Value string `json:"value"`
}

type UAAUserName struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

type UAAUserGroup struct {
	Value   string `json:"value"`
	Display string `json:"display"`
	Type    string `json:"type"`
}

func newUserManager(config coreconfig.Reader, gateway net.Gateway, identityZoneManager *IdentityZoneManager, logger *Logger) (um *UserManager, err error) {

	api, err := newUaaApi(config, gateway)
	if err != nil {
		return
	}

	um = &UserManager{
		log:                 logger,
		api:                 api,
		identityZoneManager: identityZoneManager,
		groupMap:            make(map[string]map[string]string),
		defaultGroups:       make(map[string]map[string]int),
	}
	return
}

func (um *UserManager) loadGroups(zoneId string) (err error) {

	if _, ok := um.defaultGroups[zoneId]; ok {
		// We've already populated the default groups for this zone; nothing to do
		return
	}

	uaaApi := um.api.WithZoneId(zoneId)

	um.groupMap[zoneId] = make(map[string]string)
	um.defaultGroups[zoneId] = make(map[string]int)

	// Retrieve all groups
	groupList := &UAAGroupResourceList{}
	err = uaaApi.Get("/Groups", groupList)
	if err != nil {
		return
	}
	for _, r := range groupList.Resources {
		um.groupMap[zoneId][r.DisplayName] = r.Id
	}

	// Retrieve the default groups for the identity zone
	identityZone, err := um.identityZoneManager.FindById(zoneId)
	if err != nil {
		return err
	}
	for _, g := range identityZone.Config.UserConfig.DefaultGroups {
		um.defaultGroups[zoneId][g] = 1
	}

	return
}

func (um *UserManager) IsDefaultGroup(zoneId, group string) (ok bool, err error) {

	// Make sure the groups have been loaded for this zone; will noop if so
	if err = um.loadGroups(zoneId); err == nil {
		_, ok = um.defaultGroups[zoneId][group]
	}

	return
}

func (um *UserManager) GetUser(id, zoneId string) (user *UAAUser, err error) {

	uaaApi := um.api.WithZoneId(zoneId)

	user = &UAAUser{}
	err = uaaApi.Get(fmt.Sprintf("/Users/%s", id), user)

	return
}

func (um *UserManager) CreateUser(username, password, origin, givenName, familyName, email, zoneId string) (user *UAAUser, err error) {

	uaaApi := um.api.WithZoneId(zoneId)

	userResource := UAAUser{
		Username: username,
		Password: password,
		Origin:   origin,
		Name: UAAUserName{
			GivenName:  givenName,
			FamilyName: familyName,
		},
	}
	if len(email) > 0 {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{email})
	} else {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{username})
	}

	user = &UAAUser{}
	err = uaaApi.Post("/Users", userResource, user)
	if err != nil {
		switch httpErr := err.(type) {
		case errors.HTTPError:
			if httpErr.StatusCode() == http.StatusConflict {
				err = errors.NewModelAlreadyExistsError("user", username)
			}
		}
	}

	return
}

func (um *UserManager) UpdateUser(id, username, givenName, familyName, email, zoneId string) (user *UAAUser, err error) {

	uaaApi := um.api.WithZoneId(zoneId)

	userResource := UAAUser{
		Username: username,
		Name: UAAUserName{
			GivenName:  givenName,
			FamilyName: familyName,
		},
	}
	if len(email) > 0 {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{email})
	} else {
		userResource.Emails = append(userResource.Emails, UAAUserEmail{username})
	}

	user = &UAAUser{}
	err = uaaApi.
		WithHeaders(map[string]string{
			apiheaders.IfMatch.String(): "*",
		}).
		Put(fmt.Sprintf("/Users/%s", id), userResource, user)

	return
}

func (um *UserManager) DeleteUser(id, zoneId string) error {

	return um.api.
		WithZoneId(zoneId).
		Delete(fmt.Sprintf("/Users/%s", id))
}

func (um *UserManager) ChangePassword(id, oldPassword, newPassword, zoneId string) (err error) {

	uaaApi := um.api.WithZoneId(zoneId)

	body, err := json.Marshal(map[string]string{
		"oldPassword": oldPassword,
		"password":    newPassword,
	})
	if err != nil {
		return
	}

	err = uaaApi.
		WithHeaders(map[string]string{
			apiheaders.Authorization.String(): um.clientToken,
		}).
		Put(fmt.Sprintf("/Users/%s/password", id), body, nil)

	return
}

func (um *UserManager) UpdateRoles(id string, scopesToDelete, scopesToAdd []string, origin, zoneId string) (err error) {

	// Make sure the groups have been loaded for this zone; will noop if so
	if err = um.loadGroups(zoneId); err != nil {
		return
	}

	uaaApi := um.api.WithZoneId(zoneId)

	for _, s := range scopesToDelete {
		roleID := um.groupMap[zoneId][s]
		err = uaaApi.Delete(fmt.Sprintf("/Groups/%s/members/%s", roleID, id))
	}
	for _, s := range scopesToAdd {
		roleID, exists := um.groupMap[zoneId][s]
		if !exists {
			err = fmt.Errorf("Group '%s' was not found", s)
			return
		}

		body := map[string]string{
			"origin": origin,
			"type":   "USER",
			"value":  id,
		}

		response := make(map[string]interface{})
		err = uaaApi.Post(fmt.Sprintf("/Groups/%s/members", roleID), body, &response)
		if err != nil {
			return
		}
	}

	return
}

func (um *UserManager) FindByUsername(username, zoneId string) (user UAAUser, err error) {

	uaaApi := um.api.WithZoneId(zoneId)

	usernameFilter := url.QueryEscape(fmt.Sprintf(`userName Eq "%s"`, username))
	path := fmt.Sprintf("/Users?filter=%s", usernameFilter)

	userResourceList := &UAAUserResourceList{}
	err = uaaApi.Get(path, userResourceList)
	if err != nil {
		return
	}

	if len(userResourceList.Resources) > 0 {
		user = userResourceList.Resources[0]
	} else {
		err = errors.NewModelNotFoundError("User", username)
	}

	return
}
