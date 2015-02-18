// Copyright (C) 2015 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	computeURL = "https://api.cloud.online.net"
	accountURL = "https://account.cloud.online.net"
)

// OnlineLabsClient is a client for the Online Labs Cloud API.
// UserID represents your user identifiant
// Token is to authenticate to the API
// Organization is the ID of the user's organization
type OnlineLabsClient struct {
	UserID       string
	Token        string
	Organization string
	Client       *http.Client
	ComputeURL   string
	AccountURL   string
}

// NewClient creates a new OnlineLabs API client using userId, API token and organization
func NewClient(userid string, token string, organization string) *OnlineLabsClient {
	log.Debugf("Creating client using %s %s %s", userid, token, organization)
	client := &OnlineLabsClient{
		UserID:       userid,
		Token:        token,
		Organization: organization,
		Client:       &http.Client{},
		ComputeURL:   computeURL,
		AccountURL:   accountURL,
	}
	return client
}

// GetUserInformations list informations about your user account
func (c OnlineLabsClient) GetUserInformations(userID string) (UserResponse, error) {
	var data UserResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/users/%s", c.AccountURL, userID),
		&data)
	return data, err
}

// GetUserOrganizations list all organizations associate with your account
func (c OnlineLabsClient) GetUserOrganizations() (OrganizationsResponse, error) {
	var data OrganizationsResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/organizations", c.AccountURL),
		&data)
	return data, err
}

// GetUserTokens list all tokens associate with your account
func (c OnlineLabsClient) GetUserTokens() (TokensResponse, error) {
	var data TokensResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/tokens", c.AccountURL),
		&data)
	return data, err
}

//GetUserToken lList an individual Token
func (c OnlineLabsClient) GetUserToken(tokenID string) (TokenResponse, error) {
	var data TokenResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/tokens/%s", c.AccountURL, tokenID),
		&data)
	return data, err
}

// CreateToken authenticates a user against their email, password,
// and then returns a new Token, which can be used until it expires.
// email is the user email
// password is the user password
// expires is if you want a token wich expires or not
func (c OnlineLabsClient) CreateToken(email string, password string, expires bool) (TokenResponse, error) {
	var data TokenResponse
	json := fmt.Sprintf(`{"email": "%s", "password": "%s", "expires": %t}`,
		email, password, expires)
	err := postAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/tokens", c.AccountURL),
		[]byte(json),
		&data)
	return data, err
}

// DeleteToken delete a specific token
// tokenID ith the token unique identifier
func (c OnlineLabsClient) DeleteToken(tokenID string) error {
	return deleteAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/tokens/%s", c.AccountURL, tokenID),
		nil)
}

// UpdateToken increase Token expiration time of 30 minutes
// tokenID ith the token unique identifier
func (c OnlineLabsClient) UpdateToken(tokenID string) (TokenResponse, error) {
	var data TokenResponse
	json := `{"expires": true}`
	err := patchAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/tokens/%s", c.AccountURL, tokenID),
		[]byte(json),
		&data)
	return data, err
}

// CreateServer creates a new server
// name is the server name
// organization is the organization unique identifier
// image is the image unique identifier
func (c OnlineLabsClient) CreateServer(name string, organization string, image string) (ServerResponse, error) {
	var data ServerResponse
	json := fmt.Sprintf(`{"name": "%s", "organization": "%s", "image": "%s", "tags": ["docker-machine"]}`,
		name, organization, image)
	err := postAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/servers", c.ComputeURL),
		[]byte(json),
		&data)
	return data, err
}

// DeleteServer delete a specific server
// serverID ith the server unique identifier
func (c OnlineLabsClient) DeleteServer(serverID string) error {
	return deleteAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/servers/%s", c.ComputeURL, serverID),
		nil)
}

// GetServer list an individual server
// serverID ith the server unique identifier
func (c OnlineLabsClient) GetServer(serverID string) (ServerResponse, error) {
	var data ServerResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/servers/%s", c.ComputeURL, serverID),
		&data)
	return data, err
}

// GetServers list all servers associate with your account
func (c OnlineLabsClient) GetServers() (ServersResponse, error) {
	var data ServersResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/servers", c.ComputeURL),
		&data)
	return data, err
}

// PerformServerAction execute an action on a server
// serverID ith the server unique identifier
// action is the action to execute
func (c OnlineLabsClient) PerformServerAction(serverID string, action string) (TaskResponse, error) {
	var data TaskResponse
	json := fmt.Sprintf(`{"action": "%s"}`, action)
	err := postAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/servers/%s/action", c.ComputeURL, serverID),
		[]byte(json),
		&data)
	return data, err
}

// GetVolume list an individual volume
// volumeID ith the volume unique identifier
func (c OnlineLabsClient) GetVolume(volumeID string) (VolumeResponse, error) {
	var data VolumeResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/volumes/%s", c.ComputeURL, volumeID),
		&data)
	return data, err
}

// DeleteVolume delete a specific volume
// volumeID ith the volume unique identifier
func (c OnlineLabsClient) DeleteVolume(volumeID string) error {
	return deleteAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/volumes/%s", c.ComputeURL, volumeID),
		nil)
}

// CreateVolume creates a new volume
// name is the volume name
// organization is the organization unique identifier
// volume_type is the volume type
// size is the volume size
func (c OnlineLabsClient) CreateVolume(name string, organization string, volume_type string, size int) (VolumeResponse, error) {
	var data VolumeResponse
	json := fmt.Sprintf(`{"name": "%s", "organization": "%s", "volume_type": "%s", "size": %d}`,
		name, organization, volume_type, size)
	err := postAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/volumes", c.ComputeURL),
		[]byte(json),
		&data)
	return data, err
}

// GetVolumes list all volumes associate with your account
func (c OnlineLabsClient) GetVolumes() (VolumesResponse, error) {
	var data VolumesResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/volumes", c.ComputeURL),
		&data)
	return data, err
}

// GetImages list all images associate with your account
func (c OnlineLabsClient) GetImages() (ImagesResponse, error) {
	var data ImagesResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/images", c.ComputeURL),
		&data)
	return data, err
}

// GetImage list an individual image
// volumeID ith the image unique identifier
func (c OnlineLabsClient) GetImage(volumeID string) (ImageResponse, error) {
	var data ImageResponse
	err := getAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/images/%s", c.ComputeURL, volumeID),
		&data)
	return data, err
}

// DeleteImage delete a specific volume
// volumeID ith the volume unique identifier
func (c OnlineLabsClient) DeleteImage(imageID string) error {
	return deleteAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/images/%s", c.ComputeURL, imageID),
		nil)
}

// UploadPublicKey update user SSH keys
// userId is the user unique identifier
// keyPath is the complete path of the SSH key
func (c OnlineLabsClient) UploadPublicKey(userid string, keyPath string) (UserResponse, error) {
	var data UserResponse
	publicKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return data, err
	}
	json := fmt.Sprintf(`{"ssh_public_keys": [{"key": "%s"}]}`,
		strings.TrimSpace(string(publicKey)))
	err = patchAPIResource(
		c.Client,
		c.Token,
		fmt.Sprintf("%s/users/%s", c.AccountURL, userid),
		[]byte(json),
		&data)
	return data, err
}
