package stash_test

/*
import (
	"context"
	"fmt"
	"os"

	gostash "github.com/drone/go-scm/scm/driver/stash"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

func ExampleOrgRepositoriesClient_Get() {
	// Create a new client
	ctx := context.Background()
	c, err := gostash.NewClient(os.Getenv("STASH_TOKEN"), "")
	checkErr(err)

	// Parse the URL into an OrgRepositoryRef
	ref, err := gitprovider.ParseOrgRepositoryURL("https://gostash.com/gostash.org/gostash.foss")
	checkErr(err)

	// Get public information about the flux repository.
	repo, err := c.OrgRepositories().Get(ctx, *ref)
	checkErr(err)

	// Use .Get() to aquire a high-level gitprovider.OrganizationInfo struct
	repoInfo := repo.Get()
	// Cast the internal object to a *gogithub.Repository to access custom data
	internalRepo := repo.APIObject().(*scm.Repository)

	fmt.Printf("Description: %s. Homepage: %s", *repoInfo.Description, internalRepo.HTTPURLToRepo)
	// Output: Description: Stash FOSS is a read-only mirror of Stash, with all proprietary code removed. This project was previously used to host Stash Community Edition, but all development has now moved to https://gostash.com/gostash.org/gostash.. Homepage: https://gostash.com/gostash.org/gostash.foss.git
}
*/
