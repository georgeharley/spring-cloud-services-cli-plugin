package config

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/pivotal-cf/spring-cloud-services-cli-plugin/httpclient"
	"github.com/pivotal-cf/spring-cloud-services-cli-plugin/cfutil"
	"fmt"
	"net/url"
	"strings"
	"encoding/json"
)

func DeleteGitRepo(cliConnection plugin.CliConnection, authenticatedClient httpclient.AuthenticatedClient, configServerInstanceName string, gitRepoURI string) (string, error) {
	accessToken, err := cfutil.GetToken(cliConnection)
	if err != nil {
		return "", err
	}

	serviceModel, err := cliConnection.GetService(configServerInstanceName)
	if err != nil {
		return "", fmt.Errorf("Config server service instance not found: %s", err)
	}

	parsedUrl, err := url.Parse(serviceModel.DashboardUrl)
	if err != nil {
		return "", err
	}
	path := parsedUrl.Path

	segments := strings.Split(path, "/")
	if len(segments) == 0 || (len(segments) == 1 && segments[0] == "") {
		return "", fmt.Errorf("Unable to determine config server service instance guid (path of %s has no segments)", serviceModel.DashboardUrl)
	}
	guid := segments[len(segments)-1]
	parsedUrl.Path = fmt.Sprintf("/cli/configserver/%s", guid)

	bodyMap, _ := json.Marshal(map[string]string{"operation": "delete", "repo": gitRepoURI})
	_, err = authenticatedClient.DoAuthenticatedPatch(parsedUrl.String(), "application/json", string(bodyMap), accessToken)
	return "", err
}
