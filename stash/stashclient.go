/*
Copyright 2020 The Flux CD contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stash

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

// stashClientImpl is a wrapper around *github.Client, which implements higher-level methods,
// operating on the go-github structs. Pagination is implemented for all List* methods, all returned
// objects are validated, and HTTP errors are handled/wrapped using handleHTTPError.
// This interface is also fakeable, in order to unit-test the client.
type stashClient interface {
	// Client returns the underlying *github.Client
	Client() *scm.Client

	// Group methods

	// GetGroup is a wrapper for "GET /groups/{group}".
	// This function HTTP error wrapping, and validates the server result.
	GetGroup(ctx context.Context, groupID interface{}) (*Group, error)
	// ListGroups is a wrapper for "GET /groups".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListGroups(ctx context.Context) ([]*Group, error)
	// ListSubgroups is a wrapper for "GET /groups/{group}/subgroups".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListSubgroups(ctx context.Context, groupName string) ([]Group, error)
	// ListGroupMembers is a wrapper for "GET /groups/{group}/members".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListGroupMembers(ctx context.Context, groupName string) ([]*GroupMember, error)

	// Project methods

	// GetProject is a wrapper for "GET /projects/{project}".
	// This function handles HTTP error wrapping, and validates the server result.
	GetGroupProject(ctx context.Context, groupName string, projectName string) (*scm.Repository, error)
	// ListGroupProjects is a wrapper for "GET /groups/{group}/projects".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListGroupProjects(ctx context.Context, groupName string) ([]*scm.Repository, error)
	// GetProject is a wrapper for "GET rest/api/1.0/projects/{project}".
	// This function handles HTTP error wrapping, and validates the server result.
	GetUserProject(ctx context.Context, projectName string) (*scm.Repository, error)
	// ListUserProjects is a wrapper for "GET /users/{username}/projects".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListUserProjects(ctx context.Context, username string) ([]*scm.Repository, error)
	// ListProjectUsers is a wrapper for "GET /projects/{project}/users".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListProjectUsers(ctx context.Context, projectName string) ([]*User, error)
	// CreateProject is a wrapper for "POST /projects"
	// This function handles HTTP error wrapping, and validates the server result.
	CreateProject(ctx context.Context, req *scm.Repository) (*scm.Repository, error)
	// UpdateProject is a wrapper for "PUT /projects/{project}".
	// This function handles HTTP error wrapping, and validates the server result.
	UpdateProject(ctx context.Context, req *scm.Repository) (*scm.Repository, error)
	// DeleteProject is a wrapper for "DELETE /projects/{project}".
	// This function handles HTTP error wrapping.
	// DANGEROUS COMMAND: In order to use this, you must set destructiveActions to true.
	DeleteProject(ctx context.Context, projectName string) error

	// Deploy key methods

	// ListKeys is a wrapper for "GET /projects/{project}/deploy_keys".
	// This function handles pagination, HTTP error wrapping, and validates the server result.
	ListKeys(ctx context.Context, projectName string) ([]*DeployKey, error)
	// CreateProjectKey is a wrapper for "POST /projects/{project}/deploy_keys".
	// This function handles HTTP error wrapping, and validates the server result.
	CreateKey(ctx context.Context, projectName string, req *DeployKey) (*DeployKey, error)
	// DeleteKey is a wrapper for "DELETE /projects/{project}/deploy_keys/{key_id}".
	// This function handles HTTP error wrapping.
	DeleteKey(ctx context.Context, projectName string, keyID int) error

	// Team related methods

	// ShareGroup is a wrapper for ""
	// This function handles HTTP error wrapping, and validates the server result.
	ShareProject(ctx context.Context, projectName string, groupID, groupAccess int) error
	// UnshareProject is a wrapper for ""
	// This function handles HTTP error wrapping, and validates the server result.
	UnshareProject(ctx context.Context, projectName string, groupID int) error
}

// stashClientImpl is a wrapper around *gostash.Client, which implements higher-level methods,
// operating on the go-gostash.structs. See the stashClient interface for method documentation.
// Pagination is implemented for all List* methods, all returned
// objects are validated, and HTTP errors are handled/wrapped using handleHTTPError.
type stashClientImpl struct {
	c                  *scm.Client
	destructiveActions bool
}

// stashClientImpl implements stashClient.
var _ stashClient = &stashClientImpl{}

func (c *stashClientImpl) Client() *scm.Client {
	return c.c
}

func (c *stashClientImpl) GetGroup(ctx context.Context, groupID interface{}) (*Group, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListGroups(ctx context.Context) ([]*Group, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListSubgroups(ctx context.Context, groupName string) ([]Group, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) GetGroupProject(ctx context.Context, groupName string, projectName string) (*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListGroupProjects(ctx context.Context, groupName string) ([]*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func validateProjectObjects(apiObjs []*scm.Repository) ([]*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListGroupMembers(ctx context.Context, groupName string) ([]*GroupMember, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) GetUserProject(ctx context.Context, projectName string) (*scm.Repository, error) {
	apiObj, resp, err := c.c.Repositories.Find(ctx, projectName)
	fmt.Printf("%+v\n", resp)
	return validateProjectAPIResp(apiObj, err)
}

func validateProjectAPIResp(apiObj *scm.Repository, err error) (*scm.Repository, error) {
	// If the response contained an error, return
	if err != nil {
		return nil, handleHTTPError(err)
	}
	// Make sure apiObj is valid
	if err := validateProjectAPI(apiObj); err != nil {
		return nil, err
	}
	return apiObj, nil
}

func (c *stashClientImpl) ListProjects(ctx context.Context) ([]*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListProjectUsers(ctx context.Context, projectName string) ([]*User, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListUserProjects(ctx context.Context, username string) ([]*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) CreateProject(ctx context.Context, req *scm.Repository) (*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) UpdateProject(ctx context.Context, req *scm.Repository) (*scm.Repository, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) DeleteProject(ctx context.Context, projectName string) error {
	return gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ListKeys(ctx context.Context, projectName string) ([]*DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) CreateKey(ctx context.Context, projectName string, req *DeployKey) (*DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) DeleteKey(ctx context.Context, projectName string, keyID int) error {
	return gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) ShareProject(ctx context.Context, projectName string, groupIDObj, groupAccessObj int) error {
	return gitprovider.ErrNoProviderSupport
}

func (c *stashClientImpl) UnshareProject(ctx context.Context, projectName string, groupID int) error {
	return gitprovider.ErrNoProviderSupport
}
