package workos

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/golly/errors"
	"github.com/workos/workos-go/v4/pkg/organizations"
)

func (DefaultClient) CreateOrganization(ctx golly.Context, name string, domains ...string) (string, error) {
	domainData := golly.Map(domains, func(domain string) organizations.OrganizationDomainData {
		return organizations.OrganizationDomainData{
			Domain: domain,
		}
	})

	org, err := organizations.CreateOrganization(ctx.Context(), organizations.CreateOrganizationOpts{
		AllowProfilesOutsideOrganization: len(domainData) == 0,
		Name:                             name,
		DomainData:                       domainData,
	})

	return org.ID, errors.WrapGeneric(err)
}
