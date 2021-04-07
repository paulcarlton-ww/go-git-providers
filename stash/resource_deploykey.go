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

	"reflect"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

func newDeployKey(c *DeployKeyClient, key *DeployKey) *deployKey {
	return &deployKey{
		k:       *key,
		c:       c,
		canpush: key.CanPush,
	}
}

var _ gitprovider.DeployKey = &deployKey{}

type deployKey struct {
	k       DeployKey
	c       *DeployKeyClient
	canpush *bool
}

func (dk *deployKey) Get() gitprovider.DeployKeyInfo {
	return deployKeyFromAPI(&dk.k)
}

func (dk *deployKey) Set(info gitprovider.DeployKeyInfo) error {
	if err := info.ValidateInfo(); err != nil {
		return err
	}
	deployKeyInfoToAPIObj(&info, &dk.k)
	return nil
}

func (dk *deployKey) APIObject() interface{} {
	return &dk.k
}

func (dk *deployKey) Repository() gitprovider.RepositoryRef {
	return dk.c.ref
}

// Update will apply the desired state in this object to the server.
// Only set fields will be respected (i.e. PATCH behaviour).
// In order to apply changes to this object, use the .Set({Resource}Info) error
// function, or cast .APIObject() to a pointer to the provider-specific type
// and set custom fields there.
//
// ErrNotFound is returned if the resource does not exist.
//
// The internal API object will be overridden with the received server data.
func (dk *deployKey) Update(ctx context.Context) error {
	// Delete the old key and recreate
	if err := dk.Delete(ctx); err != nil {
		return err
	}
	return dk.createIntoSelf(ctx)
}

// Delete deletes a deploy key from the repository.
//
// ErrNotFound is returned if the resource does not exist.
func (dk *deployKey) Delete(ctx context.Context) error {
	return gitprovider.ErrNoProviderSupport
}

// Reconcile makes sure the desired state in this object (called "req" here) becomes
// the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be updated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
//
// The internal API object will be overridden with the received server data if actionTaken == true.
func (dk *deployKey) Reconcile(ctx context.Context) (bool, error) {
	return false, gitprovider.ErrNoProviderSupport
}

func (dk *deployKey) createIntoSelf(ctx context.Context) error {
	// POST /repos/{owner}/{repo}/keys
	apiObj, err := dk.c.c.CreateKey(ctx, getRepoPath(dk.c.ref), &dk.k)
	if err != nil {
		return err
	}
	dk.k = *apiObj
	return nil
}

func validateDeployKeyAPI(apiObj *DeployKey) error {
	return gitprovider.ErrNoProviderSupport
}

func deployKeyFromAPI(apiObj *DeployKey) gitprovider.DeployKeyInfo {
	return gitprovider.DeployKeyInfo{}
}

func deployKeyToAPI(info *gitprovider.DeployKeyInfo) *DeployKey {
	k := &DeployKey{}
	deployKeyInfoToAPIObj(info, k)
	return k
}

func deployKeyInfoToAPIObj(info *gitprovider.DeployKeyInfo, apiObj *DeployKey) {
	// Required fields, we assume info is validated, and hence these are set
	// optional fields
	derefedBool := false
	if info.ReadOnly != nil {
		if *info.ReadOnly {
			apiObj.CanPush = &derefedBool
		} else {
			derefedBool = true
			apiObj.CanPush = &derefedBool
		}
	}
}

// This function copies over the fields that are part of create request of a deploy
// i.e. the desired spec of the deploy key. This allows us to separate "spec" from "status" fields.
func newStashKeySpec(key *DeployKey) *stashKeySpec {
	return &stashKeySpec{
		&DeployKey{},
	}
}

type stashKeySpec struct {
	*DeployKey
}

func (s *stashKeySpec) Equals(other *stashKeySpec) bool {
	return reflect.DeepEqual(s, other)
}
