package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.com/pkg/errors"
)

const (
	launchesURL = "https://api.spacexdata.com/v5/launches/query"
)

type SpaceXAPILaunchesRepo struct {
	cl *http.Client
}

func NewSpaceXAPILaunchesRepo(cl *http.Client) *SpaceXAPILaunchesRepo {
	return &SpaceXAPILaunchesRepo{cl: cl}
}

type queryRequestPayload struct {
	Query   map[string]interface{} `json:"query"`
	Options map[string]interface{} `json:"options"`
}

/*
CheckLaunches returns true if launches present for given launchpad and local date.

		note: interaction with external API good practice to cache results.
	             in given implementation it's skipped
*/
func (r *SpaceXAPILaunchesRepo) CheckLaunches(ctx context.Context, launchpad string, localDate time.Time) (bool, error) {
	b, err := preparePayload(localDate, launchpad)
	if err != nil {
		return false, errors.Wrapf(err, `failed to prepare payload: launchpad - %s, date - %s`, launchpad, localDate)
	}
	req, err := http.NewRequest(http.MethodPost, launchesURL, b)
	if err != nil {
		return false, errors.Wrapf(err, `failed to create request: url - %s, payload - %s`, launchesURL, b)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	resp, err := r.cl.Do(req)
	if err != nil {
		return false, errors.Wrapf(err, `failed to perform request: url - %s, payload - %s`, launchesURL, b)
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return false, errors.Wrapf(err, `failed to read response`)
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf(`received non success code: code - %s, response - %s`, resp.Status, data)
	}
	docs := gjson.GetBytes(data, "docs").Array()
	return len(docs) > 0, nil
}

func preparePayload(localDate time.Time, launchpad string) (*bytes.Buffer, error) {
	year, month, day := localDate.Date()
	startOfDate := time.Date(year, month, day, 0, 0, 0, 0, localDate.Location())
	startOfNextDate := startOfDate.AddDate(0, 0, 1)
	p := queryRequestPayload{
		Query: map[string]interface{}{
			"date_local": map[string]interface{}{
				"$gte": startOfDate,
				"$lt":  startOfNextDate,
			},
			"launchpad": launchpad,
		},
		Options: map[string]interface{}{
			"select": map[string]int{
				"id": 1,
			},
		},
	}
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(p); err != nil {
		return nil, errors.Wrapf(err, `failec to marshal payload: p - %+v`, p)
	}
	return b, nil
}
