package stash_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	gostash "github.com/drone/go-scm/scm/driver/stash"
	"github.com/fluxcd/go-git-providers/stash"

	"github.com/fluxcd/go-git-providers/gitprovider"
)

// checkErr is used for examples in this repository.
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestExampleOrganizationsClient_Get(t *testing.T) {
	// Create a new client
	ctx := context.Background()
	c, err := stash.NewClient(os.Getenv("STASH_TOKEN"), "")
	checkErr(err)

	// Get public information about the fluxcd organization
	org, err := c.Organizations().Get(ctx, gitprovider.OrganizationRef{
		Domain:       stash.DefaultDomain,
		Organization: "fluxcd-testing-public",
	})
	checkErr(err)

	// Use .Get() to aquire a high-level gitprovider.OrganizationInfo struct
	orgInfo := org.Get()
	// Cast the internal object to a *gogithub.Organization to access custom data
	internalOrg := org.APIObject().(*gostash.Error)

	fmt.Printf("Name: %s. Location: %s.", *orgInfo.Name, internalOrg.Error())
	// Output: Name: Flux project. Location: CNCF incubation.
}
