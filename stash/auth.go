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
	"fmt"
	"net/http"
	"net/url"

	"github.com/drone/go-scm/scm"
	gostash "github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport"
	scmoauth2 "github.com/drone/go-scm/scm/transport/oauth2"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/fluxcd/go-git-providers/gitprovider/cache"
	"golang.org/x/oauth2"
)

const (
	// DefaultDomain specifies the default domain used as the backend.
	DefaultDomain = "gostash.com"
)

// ClientOption is the interface to implement for passing options to NewClient.
// The clientOptions struct is private to force usage of the With... functions.
type ClientOption interface {
	// ApplyToStashClientOptions applies set fields of this object into target.
	ApplyToStashClientOptions(target *clientOptions) error
}

// clientOptions is the struct that tracks data about what options have been set.
type clientOptions struct {
	// clientOptions shares all the common options
	gitprovider.CommonClientOptions

	// AuthTransport is a ChainableRoundTripperFunc adding authentication credentials to the transport chain.
	AuthTransport gitprovider.ChainableRoundTripperFunc

	// EnableConditionalRequests will be set if conditional requests should be used.
	EnableConditionalRequests *bool
}

// ApplyToStashClientOptions implements ClientOption, and applies the set fields of opts
// into target. If both opts and target has the same specific field set, ErrInvalidClientOptions is returned.
func (opts *clientOptions) ApplyToStashClientOptions(target *clientOptions) error {
	// Apply common values, if any
	if err := opts.CommonClientOptions.ApplyToCommonClientOptions(&target.CommonClientOptions); err != nil {
		return err
	}

	if opts.AuthTransport != nil {
		// Make sure the user didn't specify the AuthTransport twice
		if target.AuthTransport != nil {
			return fmt.Errorf("option AuthTransport already configured: %w", gitprovider.ErrInvalidClientOptions)
		}
		target.AuthTransport = opts.AuthTransport
	}

	if opts.EnableConditionalRequests != nil {
		// Make sure the user didn't specify the EnableConditionalRequests twice
		if target.EnableConditionalRequests != nil {
			return fmt.Errorf("option EnableConditionalRequests already configured: %w", gitprovider.ErrInvalidClientOptions)
		}
		target.EnableConditionalRequests = opts.EnableConditionalRequests
	}
	return nil
}

// getTransportChain builds the full chain of transports (from left to right,
// as per gitprovider.BuildClientFromTransportChain) of the form described in NewClient.
func (opts *clientOptions) getTransportChain() (chain []gitprovider.ChainableRoundTripperFunc) {
	if opts.PostChainTransportHook != nil {
		chain = append(chain, opts.PostChainTransportHook)
	}
	if opts.AuthTransport != nil {
		chain = append(chain, opts.AuthTransport)
	}
	if opts.EnableConditionalRequests != nil && *opts.EnableConditionalRequests {
		// TODO: Provide some kind of debug logging if/when the httpcache is used
		// One can see if the request hit the cache using: resp.Header[httpcache.XFromCache]
		chain = append(chain, cache.NewHTTPCacheTransport)
	}
	if opts.PreChainTransportHook != nil {
		chain = append(chain, opts.PreChainTransportHook)
	}
	return
}

// buildCommonOption is a helper for returning a ClientOption out of a common option field.
func buildCommonOption(opt gitprovider.CommonClientOptions) *clientOptions {
	return &clientOptions{CommonClientOptions: opt}
}

// errorOption implements ClientOption, and just wraps an error which is immediately returned.
// This struct can be used through the optionError function, in order to make makeOptions fail
// if there are invalid options given to the With... functions.
type errorOption struct {
	err error
}

// ApplyToStashClientOptions implements ClientOption, but just returns the internal error.
func (e *errorOption) ApplyToStashClientOptions(*clientOptions) error { return e.err }

// optionError is a constructor for errorOption.
func optionError(err error) ClientOption {
	return &errorOption{err}
}

//
// Common options
//

// WithDomain initializes a Client for a custom Stash instance of the given domain.
// Only host and port information should be present in domain. domain must not be an empty string.
func WithDomain(domain string) ClientOption {
	return buildCommonOption(gitprovider.CommonClientOptions{Domain: &domain})
}

// WithDestructiveAPICalls tells the client whether it's allowed to do dangerous and possibly destructive
// actions, like e.g. deleting a repository.
func WithDestructiveAPICalls(destructiveActions bool) ClientOption {
	return buildCommonOption(gitprovider.CommonClientOptions{EnableDestructiveAPICalls: &destructiveActions})
}

// WithPreChainTransportHook registers a ChainableRoundTripperFunc "before" the cache and authentication
// transports in the chain. For more information, see NewClient, and gitprovider.CommonClientOptions.PreChainTransportHook.
func WithPreChainTransportHook(preRoundTripperFunc gitprovider.ChainableRoundTripperFunc) ClientOption {
	// Don't allow an empty value
	if preRoundTripperFunc == nil {
		return optionError(fmt.Errorf("preRoundTripperFunc cannot be nil: %w", gitprovider.ErrInvalidClientOptions))
	}

	return buildCommonOption(gitprovider.CommonClientOptions{PreChainTransportHook: preRoundTripperFunc})
}

// WithPostChainTransportHook registers a ChainableRoundTripperFunc "after" the cache and authentication
// transports in the chain. For more information, see NewClient, and gitprovider.CommonClientOptions.WithPostChainTransportHook.
func WithPostChainTransportHook(postRoundTripperFunc gitprovider.ChainableRoundTripperFunc) ClientOption {
	// Don't allow an empty value
	if postRoundTripperFunc == nil {
		return optionError(fmt.Errorf("postRoundTripperFunc cannot be nil: %w", gitprovider.ErrInvalidClientOptions))
	}

	return buildCommonOption(gitprovider.CommonClientOptions{PostChainTransportHook: postRoundTripperFunc})
}

// WithOAuth2Token initializes a Client which authenticates with Stash through an OAuth2 token.
// oauth2Token must not be an empty string.
func WithOAuth2Token(oauth2Token string) ClientOption {
	// Don't allow an empty value
	if len(oauth2Token) == 0 {
		return optionError(fmt.Errorf("oauth2Token cannot be empty: %w", gitprovider.ErrInvalidClientOptions))
	}

	return &clientOptions{AuthTransport: oauth2Transport(oauth2Token)}
}

func oauth2Transport(oauth2Token string) gitprovider.ChainableRoundTripperFunc {
	return func(in http.RoundTripper) http.RoundTripper {
		// Create a TokenSource of the given access token
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauth2Token})
		// Create a Transport, with "in" as the underlying transport, and the given TokenSource
		return &oauth2.Transport{
			Base:   in,
			Source: oauth2.ReuseTokenSource(nil, ts),
		}
	}
}

// WithConditionalRequests instructs the client to use Conditional Requests to Stash.
// See: https://gitlab.com/gitlab.org/gitlab.foss/-/issues/26926, and
// https://docs.gostash.com/ee/development/polling.html for more info.
func WithConditionalRequests(conditionalRequests bool) ClientOption {
	return &clientOptions{EnableConditionalRequests: &conditionalRequests}
}

// makeOptions assembles a clientOptions struct from ClientOption mutator functions.
func makeOptions(opts ...ClientOption) (*clientOptions, error) {
	o := &clientOptions{}
	for _, opt := range opts {
		if err := opt.ApplyToStashClientOptions(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// NewClient creates a new gostash.Client instance for Stash API endpoints.
func NewClient(token string, tokenType string, optFns ...ClientOption) (gitprovider.Client, error) {
	var st *scm.Client
	var domain, sshDomain string
	var url *url.URL

	// Complete the options struct
	opts, err := makeOptions(optFns...)
	if err != nil {
		return nil, err
	}

	// Create a *http.Client using the transport chain
	httpClient, err := gitprovider.BuildClientFromTransportChain(opts.getTransportChain())
	if err != nil {
		return nil, err
	}

	if opts.Domain == nil || *opts.Domain == DefaultDomain {
		// No domain set or the default localhost used
		domain = DefaultDomain
	} else {
		domain = *opts.Domain
	}

	url, err = url.Parse(domain)
	if err != nil {
		return nil, err
	}

	st, err = gostash.New(fmt.Sprintf("http://%s", domain))
	if err != nil {
		return nil, err
	}

	if tokenType == "oauth2" {
		// provide a custom http.Client with a transport
		// that injects the oauth2 token.
		httpClient.Transport = &scmoauth2.Transport{
			Source: scmoauth2.StaticTokenSource(
				&scm.Token{
					Token: token,
				},
			),
		}
		st.Client = httpClient
	} else {
		// provide a custom http.Client with a transport
		// that injects the private token.
		httpClient.Transport = &transport.PrivateToken{Token: token}
		st.Client = httpClient
	}

	// By default, turn destructive actions off. But allow overrides.
	destructiveActions := false
	if opts.EnableDestructiveAPICalls != nil {
		destructiveActions = *opts.EnableDestructiveAPICalls
	}

	return newClient(st, domain, sshDomain, destructiveActions), nil
}
