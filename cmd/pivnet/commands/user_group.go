package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf-experimental/go-pivnet"
)

type UserGroupsCommand struct {
	ProductSlug    string `long:"product-slug" description:"Product slug e.g. p-mysql"`
	ReleaseVersion string `long:"release-version" description:"Release version e.g. 0.1.2-rc1"`
}

func (command *UserGroupsCommand) Execute([]string) error {
	client := NewClient()

	if command.ProductSlug == "" && command.ReleaseVersion == "" {
		var err error
		userGroups, err := client.UserGroups.List()
		if err != nil {
			return err
		}

		return printUserGroups(userGroups)
	}

	if command.ProductSlug == "" || command.ReleaseVersion == "" {
		return fmt.Errorf("Both or neither of product slug and release version must be provided")
	}

	releases, err := client.Releases.List(command.ProductSlug)
	if err != nil {
		return err
	}

	var release pivnet.Release
	for _, r := range releases {
		if r.Version == command.ReleaseVersion {
			release = r
			break
		}
	}

	if release.Version != command.ReleaseVersion {
		return fmt.Errorf("release not found")
	}

	userGroups, err := client.UserGroups.ListForRelease(command.ProductSlug, release.ID)
	if err != nil {
		return err
	}

	return printUserGroups(userGroups)
}

func printUserGroups(userGroups []pivnet.UserGroup) error {

	switch Pivnet.Format {
	case printAsTable:
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Description"})

		for _, u := range userGroups {
			table.Append([]string{
				strconv.Itoa(u.ID), u.Name, u.Description,
			})
		}
		table.Render()
		return nil
	case printAsJSON:
		b, err := json.Marshal(userGroups)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(b))
		return nil
	case printAsYAML:
		b, err := yaml.Marshal(userGroups)
		if err != nil {
			return err
		}

		fmt.Printf("---\n%s\n", string(b))
		return nil
	}

	return nil
}
