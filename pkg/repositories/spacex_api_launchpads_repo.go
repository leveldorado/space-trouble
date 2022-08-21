package repositories

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const (
	launchpadURLPrefix = "https://api.spacexdata.com/v4/launchpads/"
)

type SpaceXAPILaunchpadsRepo struct {
	cl *http.Client
}

func NewSpaceXAPILaunchpadsRepo(cl *http.Client) *SpaceXAPILaunchpadsRepo {
	return &SpaceXAPILaunchpadsRepo{cl: cl}
}

/*
Get returns launchpad for provided id.

		note: interaction with external API good practice to cache results.
	             in given implementation it's skipped
*/
func (r *SpaceXAPILaunchpadsRepo) Get(ctx context.Context, id string) (types.Launchpad, error) {
	u, err := url.Parse(launchpadURLPrefix + id)
	if err != nil {
		return types.Launchpad{}, err
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return types.Launchpad{}, errors.Wrapf(err, `failed to create request: url - %s`, u)
	}
	req = req.WithContext(ctx)
	resp, err := r.cl.Do(req)
	if err != nil {
		return types.Launchpad{}, errors.Wrapf(err, `failed to do request: url - %s`, u)
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return types.Launchpad{}, errors.Wrap(err, `failed to read body`)
	}
	if resp.StatusCode == http.StatusNotFound {
		return types.Launchpad{}, types.ErrNotFound{}
	}
	if resp.StatusCode != http.StatusOK {
		return types.Launchpad{}, errors.Errorf(`non success response code: code - %s, boody - %s`, resp.Status, data)
	}
	pad := types.Launchpad{
		ID:       id,
		FullName: gjson.GetBytes(data, "full_name").String(),
	}
	timezone := gjson.GetBytes(data, "timezone").String()
	location, err := time.LoadLocation(timezone)
	pad.Location = location
	return pad, errors.Wrapf(err, `failed to load location: timezone - %s`, timezone)
}
