package github

import (
	"github.com/ONSdigital/graphql"
	"github.com/pkg/errors"
)

type (

	// Data represents GitHub enterprise licensing information.
	Data struct {
		Enterprise struct {
			BillingInfo struct {
				TotalAvailableLicenses int `json:"totalAvailableLicenses,omitempty"`
				TotalLicenses          int `json:"totalLicenses,omitempty"`
			}
		}
	}
	// }
)

// GetEnterpriseLicensing returns counts of total available licenses and total licenses.
func (c Client) GetEnterpriseLicensing(enterprise string) (data *Data, err error) {
	req := graphql.NewRequest(`
		query GitHubEnterpriseLicensing($slug: String!) {
			enterprise(slug: $slug) {
				billingInfo {
					totalAvailableLicenses
					totalLicenses
				}
			}
		}
	`)

	req.Var("slug", enterprise)

	if err := c.Run(req, &data); err != nil {
		return nil, errors.Wrap(err, "failed to fetch licensing for enterprise")
	}

	return data, nil
}
