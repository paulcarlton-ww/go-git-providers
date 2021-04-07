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

	"github.com/fluxcd/go-git-providers/gitprovider"
)

// DeployKeyClient implements the gitprovider.DeployKeyClient interface.
var _ gitprovider.DeployKeyClient = &DeployKeyClient{}

// DeployKeyClient operates on the access deploy key list for a specific repository.
type DeployKeyClient struct {
	*clientContext
	ref gitprovider.RepositoryRef
}

// Get returns the repository at the given path.
//
// ErrNotFound is returned if the resource does not exist.
func (c *DeployKeyClient) Get(ctx context.Context, deployKeyName string) (gitprovider.DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *DeployKeyClient) get(ctx context.Context, deployKeyName string) (*deployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

// List lists all repository deploy keys of the given deploy key type.
//
// List returns all available repository deploy keys for the given type,
// using multiple paginated requests if needed.
func (c *DeployKeyClient) List(ctx context.Context) ([]gitprovider.DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

func (c *DeployKeyClient) list(ctx context.Context) ([]*deployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

// Create creates a deploy key with the given specifications.
//
// ErrAlreadyExists will be returned if the resource already exists.
func (c *DeployKeyClient) Create(ctx context.Context, req gitprovider.DeployKeyInfo) (gitprovider.DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}

// Reconcile makes sure the given desired state (req) becomes the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be deleted and recreated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
func (c *DeployKeyClient) Reconcile(ctx context.Context, req gitprovider.DeployKeyInfo) (gitprovider.DeployKey, bool, error) {
	return nil, false, gitprovider.ErrNoProviderSupport
}

func createDeployKey(ctx context.Context, c stashClient, ref gitprovider.RepositoryRef, req gitprovider.DeployKeyInfo) (*DeployKey, error) {
	return nil, gitprovider.ErrNoProviderSupport
}
