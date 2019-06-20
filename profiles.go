package mixpanel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var exportAPIClient = &http.Client{Timeout: time.Minute}

type ExportClient struct {
	Secret string
	Client *http.Client
}

func NewExportClient(apiSecret string) *ExportClient {
	return &ExportClient{
		Secret: apiSecret,
		Client: exportAPIClient,
	}
}

func (c *ExportClient) get(method string, endpoint string, paramMap map[string]string, dest interface{}) error {
	var (
		err error
		req *http.Request
		r   io.Reader
	)

	if endpoint == "" {
		return fmt.Errorf("endpoint missing")
	}

	endpoint = fmt.Sprintf("https://mixpanel.com/api/2.0/%s", endpoint)

	if paramMap == nil {
		paramMap = map[string]string{}
	}

	params := url.Values{}
	for k, v := range paramMap {
		params[k] = []string{v}
	}

	switch method {
	case "GET":
		enc := params.Encode()
		if enc != "" {
			endpoint = endpoint + "?" + enc
		}
	case "POST":
		r = strings.NewReader(params.Encode())
	default:
		return fmt.Errorf("method not supported: %v", method)
	}

	req, err = http.NewRequest(method, endpoint, r)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Secret, "")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, dest)
}

type Profile struct {
	ID         string                 `json:"$distinct_id"`
	Properties map[string]interface{} `json:"$properties"`
}

type ListReponse struct {
	Status     string    `json:"status"`
	Error      string    `json:"error"`
	SessionID  string    `json:"session_id"`
	ComputedAt time.Time `json:"computed_at"`
	Results    []Profile `json:"results"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
}

type ProfileQuery struct {
	LastSeenAfter    time.Time
	LastSeenBefore   time.Time
	OutputProperties []string
}

func buildQueryString(q *ProfileQuery) (bool, string) {
	parts := []string{}

	if !q.LastSeenBefore.IsZero() {
		parts = append(parts, fmt.Sprintf(`properties["$last_seen"] < datetime(%d)`, q.LastSeenBefore.Unix()))
	}

	if !q.LastSeenAfter.IsZero() {
		parts = append(parts, fmt.Sprintf(`properties["$last_seen"] > datetime(%d)`, q.LastSeenAfter.Unix()))
	}

	if len(parts) == 0 {
		return false, ""
	}

	return true, "(" + strings.Join(parts, " && ") + ")"
}

func mapStr(input []string, mapFunc func(input string) string) []string {
	r := make([]string, len(input))
	for idx, elem := range input {
		r[idx] = mapFunc(elem)
	}
	return r
}

func addQuotes(input string) string {
	return `"` + strings.Replace(input, `"`, `\"`, -1) + `"`
}

func (c *ExportClient) ListProfiles(q *ProfileQuery) ([]Profile, error) {
	list := []Profile{}
	sessID := ""
	page := 0
	total := 0

	for {
		props := map[string]string{}

		if ok, qStr := buildQueryString(q); ok {
			props["where"] = qStr
		}

		if q.OutputProperties != nil {
			props["output_properties"] = "[" + strings.Join(mapStr(q.OutputProperties, addQuotes), ", ") + "]"
		}

		if sessID != "" && page > 0 {
			props["session_id"] = sessID
			props["page"] = fmt.Sprintf("%d", page)
		}

		r := ListReponse{}
		if err := c.get("GET", "engage", props, &r); err != nil {
			return nil, err
		}

		if r.Error != "" {
			return nil, fmt.Errorf("server error: %s", r.Error)
		}

		list = append(list, r.Results...)

		if r.Total > 0 {
			total = r.Total
		}

		if len(list) >= total {
			break
		}

		sessID = r.SessionID
		page = r.Page + 1
	}

	return list, nil
}

func (c *Client) DeleteProfile(distinctID string) error {
	return c.makeRequestWithData("POST", "engage", Properties{
		"$distinct_id":  distinctID,
		"$token":        c.Token,
		"$ignore_alias": "true",
		"$delete":       "",
	}, sourceScript)
}
