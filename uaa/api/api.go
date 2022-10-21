package api

import (
	"bytes"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/net"
	"encoding/json"
	"errors"
	"fmt"
	apiheaders "github.com/foundcloudry/terraform-provider-uaa/uaa/api/headers"
	"net/http"
	"strings"
)

type UaaApi struct {
	additionalHeaders map[string]string
	baseUrl           string
	config            coreconfig.Reader
	gateway           net.Gateway
	zoneId            string
}

func newUaaApi(config coreconfig.Reader, gateway net.Gateway) (*UaaApi, error) {
	if config.UaaEndpoint() == "" {
		return nil, errors.New("no UAA endpoint provided when instantiating the UAA API")
	}

	return &UaaApi{
		additionalHeaders: make(map[string]string),
		baseUrl:           config.UaaEndpoint(),
		config:            config,
		gateway:           gateway,
	}, nil
}

func (api *UaaApi) WithHeaders(headers map[string]string) *UaaApi {
	additionalHeaders := make(map[string]string)

	for i, v := range api.additionalHeaders {
		additionalHeaders[i] = v
	}
	for i, v := range headers {
		additionalHeaders[i] = v
	}

	return &UaaApi{
		additionalHeaders: additionalHeaders,
		baseUrl:           api.baseUrl,
		config:            api.config,
		gateway:           api.gateway,
		zoneId:            api.zoneId,
	}
}

func (api *UaaApi) WithZoneId(zoneId string) *UaaApi {
	return &UaaApi{
		additionalHeaders: api.additionalHeaders,
		baseUrl:           api.baseUrl,
		config:            api.config,
		gateway:           api.gateway,
		zoneId:            zoneId,
	}
}

func (api *UaaApi) newRequest(method, path string, body any, responseBody any) error {

	path = strings.Replace(path, "//", "/", -1)
	path = strings.TrimPrefix(path, "/")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := api.gateway.NewRequest(
		method,
		fmt.Sprintf("%s/%s", api.baseUrl, path),
		api.config.AccessToken(),
		bytes.NewReader(jsonBody),
	)

	//buf := new(strings.Builder)
	//_, err = io.Copy(buf, request.HTTPReq.Body)
	//return fmt.Errorf("***** WTF = %s", buf)

	if err != nil {
		return err
	}

	request.HTTPReq.Header.Set(apiheaders.ZoneId.String(), api.zoneId)
	for i, v := range api.additionalHeaders {
		request.HTTPReq.Header.Del(i)
		request.HTTPReq.Header.Set(i, v)
	}

	_, err = api.gateway.PerformRequestForJSONResponse(request, &responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (api *UaaApi) Get(path string, responseBody any) error {
	return api.newRequest(http.MethodGet, path, nil, responseBody)
}

func (api *UaaApi) Post(path string, body any, responseBody any) error {
	return api.newRequest(http.MethodPost, path, body, responseBody)
}

func (api *UaaApi) Patch(path string, body any, responseBody any) error {
	return api.newRequest(http.MethodPatch, path, body, responseBody)
}

func (api *UaaApi) Put(path string, body any, responseBody any) error {
	return api.newRequest(http.MethodPut, path, body, responseBody)
}

func (api *UaaApi) Delete(path string) error {
	return api.newRequest(http.MethodDelete, path, nil, nil)
}
