package api

import (
	"github.com/foundcloudry/terraform-provider-uaa/uaa/envvars"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/configuration"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/net"
)

type Session struct {
	Log *Logger

	config     coreconfig.Repository
	uaaGateway net.Gateway

	authManager        *AuthManager
	clientManager      *ClientManager
	groupManager       *GroupManager
	identityZoneManger *IdentityZoneManager
	userManager        *UserManager
}

type Config struct {
	LoginEndpoint     string
	AuthEndpoint      string
	ClientID          string
	ClientSecret      string
	CaCert            string
	SkipSslValidation bool
}

func (config *Config) Client() (*Session, error) {
	return NewSession(config)
}

type uaaErrorResponse struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func NewSession(config *Config) (s *Session, err error) {

	s = &Session{}

	envDialTimeout := os.Getenv(envvars.UaaDialTimeout.String())

	debug, _ := strconv.ParseBool(os.Getenv(envvars.UaaDebug.String()))
	s.Log = NewLogger(debug, os.Getenv(envvars.UaaTrace.String()))

	s.config = coreconfig.NewRepositoryFromPersistor(newNoopPersistor(), func(err error) {
		if err != nil {
			s.Log.UI.Failed(err.Error())
			os.Exit(1)
		}
	})
	if i18n.T == nil {
		i18n.T = i18n.Init(s.config)
	}
	s.config.SetSSLDisabled(config.SkipSslValidation)

	s.config.SetAuthenticationEndpoint(endpointAsURL(config.LoginEndpoint))
	s.config.SetUaaEndpoint(endpointAsURL(config.AuthEndpoint))

	s.uaaGateway = net.NewUAAGateway(s.config, s.Log.UI, s.Log.TracePrinter, envDialTimeout)
	s.authManager = NewAuthManager(s.uaaGateway, s.config, net.NewRequestDumper(s.Log.TracePrinter))
	//s.uaaGateway.SetTokenRefresher(s.authManager)

	s.identityZoneManger, err = newIdentityZoneManager(s.config, s.uaaGateway, s.Log)
	if err != nil {
		return nil, err
	}

	s.userManager, err = newUserManager(s.config, s.uaaGateway, s.identityZoneManger, s.Log)
	if err != nil {
		return nil, err
	}

	s.groupManager, err = newGroupManager(s.config, s.uaaGateway, s.Log)
	if err != nil {
		return nil, err
	}

	s.clientManager, err = newClientManager(s.config, s.uaaGateway, s.Log)
	if err != nil {
		return nil, err
	}

	s.userManager.clientToken, err = s.authManager.GetClientToken(config.ClientID, config.ClientSecret, "")

	return
}

func (s *Session) UserManager() *UserManager {
	return s.userManager
}

func (s *Session) ClientManager() *ClientManager {
	return s.clientManager
}

func (s *Session) GroupManager() *GroupManager {
	return s.groupManager
}

func (s *Session) IdentityZoneManager() *IdentityZoneManager {
	return s.identityZoneManger
}

func (s *Session) AuthManager() *AuthManager {
	return s.authManager
}

type noopPersistor struct {
}

func newNoopPersistor() configuration.Persistor {
	return &noopPersistor{}
}

func (p *noopPersistor) Delete() {
}

func (p *noopPersistor) Exists() bool {
	return false
}

func (p *noopPersistor) Load(configuration.DataInterface) error {
	return nil
}

func (p *noopPersistor) Save(configuration.DataInterface) error {
	return nil
}

func endpointAsURL(endpoint string) string {

	endpoint = strings.TrimSuffix(endpoint, "/")
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	return endpoint
}
