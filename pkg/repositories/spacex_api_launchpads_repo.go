package repositories

import (
	"bytes"
	"context"
	"encoding/json"
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
	launchpadsURL      = "https://api.spacexdata.com/v4/launchpads/query"
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
	return parseLaunchpad(string(data))
}

func parseLaunchpad(data string) (types.Launchpad, error) {
	pad := types.Launchpad{
		ID:       gjson.Get(data, "id").String(),
		FullName: gjson.Get(data, "full_name").String(),
		Status:   gjson.Get(data, "status").String(),
	}
	timezone := gjson.Get(data, "timezone").String()
	location, err := time.LoadLocation(timezone)
	pad.Location = location
	return pad, errors.Wrapf(err, `failed to load location: timezone - %s`, timezone)
}

func (r *SpaceXAPILaunchpadsRepo) List(ctx context.Context) ([]types.Launchpad, error) {
	var launchpads []types.Launchpad
	const limit = 10
	var currentOffset int
	for {
		list, total, err := r.queryList(ctx, limit, currentOffset)
		if err != nil {
			return nil, errors.Wrap(err, `failed to query list`)
		}
		launchpads = append(launchpads, list...)
		if len(launchpads) == total {
			return launchpads, nil
		}
		currentOffset += len(list)
	}
}

func (r *SpaceXAPILaunchpadsRepo) queryList(ctx context.Context, limit, offset int) ([]types.Launchpad, int, error) {
	b, err := prepareLaunchpadsPayload(limit, offset)
	if err != nil {
		return nil, 0, errors.Wrapf(err, `failed to prepare paylaod: limit - %d, offset - %d`, limit, offset)
	}
	req, err := http.NewRequest(http.MethodPost, launchpadsURL, b)
	if err != nil {
		return nil, 0, errors.Wrapf(err, `failed to create request: url - %s, payload - %s`, launchpadsURL, b)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := r.cl.Do(req)
	if err != nil {
		return nil, 0, errors.Wrapf(err, `failed to do request: url - %s, payload - %s`, launchpadsURL, b)
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, 0, errors.Wrapf(err, `failed to read response`)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.Errorf(`received non success code: code - %s, response - %s`, resp.Status, data)
	}
	docs := gjson.GetBytes(data, "docs").Array()
	totalDocs := int(gjson.GetBytes(data, "totalDocs").Int())
	var launchpads []types.Launchpad
	for _, el := range docs {
		pad, err := parseLaunchpad(el.Raw)
		if err != nil {
			return nil, 0, errors.Wrapf(err, `failed to parse launchpad: raw - %s`, el.Raw)
		}
		launchpads = append(launchpads, pad)
	}
	return launchpads, totalDocs, nil
}

func prepareLaunchpadsPayload(limit, offset int) (*bytes.Buffer, error) {
	req := queryRequestPayload{
		Query: map[string]interface{}{
			"status": "active",
		},
		Options: map[string]interface{}{
			"select": map[string]int{
				"timezone":  1,
				"full_name": 1,
			},
			"limit":  limit,
			"offset": offset,
			"sort": map[string]int{
				"full_name": 1,
			},
		},
	}
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(req)
	return b, errors.Wrapf(err, `failed marshal request: req - %+v`, req)
}
