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
	"errors"

	"github.com/drone/go-scm/scm"
	_ "github.com/drone/go-scm/scm"
	_ "github.com/drone/go-scm/scm/driver/stash"
	"github.com/google/go-cmp/cmp"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

func newUserProject(ctx *clientContext, apiObj *scm.Repository, ref gitprovider.RepositoryRef) *userProject {
	return &userProject{
		clientContext: ctx,
		p:             *apiObj,
		ref:           ref,
		deployKeys: &DeployKeyClient{
			clientContext: ctx,
			ref:           ref,
		},
	}
}

var _ gitprovider.UserRepository = &userProject{}

type userProject struct {
	*clientContext

	p   scm.Repository
	ref gitprovider.RepositoryRef

	deployKeys *DeployKeyClient
}

func (p *userProject) Get() gitprovider.RepositoryInfo {
	return repositoryFromAPI(&p.p)
}

func (p *userProject) Set(info gitprovider.RepositoryInfo) error {
	if err := info.ValidateInfo(); err != nil {
		return err
	}
	repositoryInfoToAPIObj(&info, &p.p)
	return nil
}

func (p *userProject) APIObject() interface{} {
	return &p.p
}

func (p *userProject) Repository() gitprovider.RepositoryRef {
	return p.ref
}

func (p *userProject) DeployKeys() gitprovider.DeployKeyClient {
	return p.deployKeys
}

// The internal API object will be overridden with the received server data.
func (p *userProject) Update(ctx context.Context) error {
	// PATCH /repos/{owner}/{repo}
	apiObj, err := p.c.UpdateProject(ctx, &p.p)
	if err != nil {
		return err
	}
	p.p = *apiObj
	return nil
}

// Reconcile makes sure the desired state in this object (called "req" here) becomes
// the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be updated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
//
// The internal API object will be overridden with the received server data if actionTaken == true.
func (p *userProject) Reconcile(ctx context.Context) (bool, error) {
	apiObj, err := p.c.GetUserProject(ctx, getRepoPath(p.ref))
	if err != nil {
		// Create if not found
		if errors.Is(err, gitprovider.ErrNotFound) {
			// orgName := ""
			// if orgRef, ok := p.ref.(gitprovider.OrgRepositoryRef); ok {
			// 	orgName = orgRef.Organization
			// }
			project, err := p.c.CreateProject(ctx, &p.p)
			if err != nil {
				return true, err
			}
			p.p = *project
			return true, nil
		}

		return false, err
	}

	// Use wrappers here to extract the "spec" part of the object for comparison
	desiredSpec := newStashProjectSpec(&p.p)
	actualSpec := newStashProjectSpec(apiObj)

	// If desired state already is the actual state, do nothing
	if desiredSpec.Equals(actualSpec) {
		return false, nil
	}
	// Otherwise, make the desired state the actual state
	return true, p.Update(ctx)
}

// Delete deletes the current resource irreversibly.
//
// ErrNotFound is returned if the resource doesn't exist anymore.
func (p *userProject) Delete(ctx context.Context) error {
	return p.c.DeleteProject(ctx, getRepoPath(p.ref))
}

func newGroupProject(ctx *clientContext, apiObj *scm.Repository, ref gitprovider.RepositoryRef) *orgRepository {
	return &orgRepository{
		userProject: *newUserProject(ctx, apiObj, ref),
		teamAccess: &TeamAccessClient{
			clientContext: ctx,
			ref:           ref,
		},
	}
}

var _ gitprovider.OrgRepository = &orgRepository{}

type orgRepository struct {
	userProject

	teamAccess *TeamAccessClient
}

func (r *orgRepository) TeamAccess() gitprovider.TeamAccessClient {
	return r.teamAccess
}

// Reconcile makes sure the desired state in this object (called "req" here) becomes
// the actual state in the backing Git provider.
//
// If req doesn't exist under the hood, it is created (actionTaken == true).
// If req doesn't equal the actual state, the resource will be updated (actionTaken == true).
// If req is already the actual state, this is a no-op (actionTaken == false).
//
// The internal API object will be overridden with the received server data if actionTaken == true.
func (r *orgRepository) Reconcile(ctx context.Context) (bool, error) {
	apiObj, err := r.c.GetGroupProject(ctx, r.ref.GetIdentity(), r.ref.GetRepository())
	if err != nil {
		// Create if not found
		if errors.Is(err, gitprovider.ErrNotFound) {
			project, err := r.c.CreateProject(ctx, &r.p)
			if err != nil {
				return true, err
			}
			r.p = *project
			return true, nil
		}

		return false, err
	}

	// Use wrappers here to extract the "spec" part of the object for comparison
	desiredSpec := newStashProjectSpec(&r.p)
	actualSpec := newStashProjectSpec(apiObj)

	// If desired state already is the actual state, do nothing
	if desiredSpec.Equals(actualSpec) {
		return false, nil
	}
	// Otherwise, make the desired state the actual state
	return true, r.Update(ctx)
}

func repositoryFromAPI(apiObj *scm.Repository) gitprovider.RepositoryInfo {
	repo := gitprovider.RepositoryInfo{}
	return repo
}

func repositoryToAPI(repo *gitprovider.RepositoryInfo, ref gitprovider.RepositoryRef) *scm.Repository {
	apiObj := &scm.Repository{
		Name: *gitprovider.StringVar(ref.GetRepository()),
	}
	repositoryInfoToAPIObj(repo, apiObj)
	return apiObj
}

func repositoryInfoToAPIObj(repo *gitprovider.RepositoryInfo, apiObj *scm.Repository) {
}

// This function copies over the fields that are part of create/update requests of a project
// i.e. the desired spec of the repository. This allows us to separate "spec" from "status" fields.
func newStashProjectSpec(project *scm.Repository) *stashProjectSpec {
	return &stashProjectSpec{
		&scm.Repository{
			// Generic
			Name: project.Name,
		},
	}
}

type stashProjectSpec struct {
	*scm.Repository
}

func (s *stashProjectSpec) Equals(other *stashProjectSpec) bool {
	return cmp.Equal(s, other)
}
