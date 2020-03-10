// Package mixpanel is an unofficial Go client library for the services provided by MixPanel
package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiBaseURL = "https://api.mixpanel.com"
	library    = "vizzlo/mixpanel"
)

// actionSource describes where the call is originating from.
// This determines whether or not the location properties on a profile should be updated.
type actionSource int

const (
	// sourceUser flags the action as having originated with the user in question.
	sourceUser actionSource = iota
	// sourceScript means that the action originated in a backend script.
	// The IP should not be tracked in this case.
	sourceScript
)

// Client is a client to talk to the API
type Client struct {
	Token   string
	BaseURL string
	Client  *http.Client
}

// Properties are key=value pairs that decorate an event or a profile.
type Properties map[string]interface{}

// Operation is an action performed on a user profile.
// Typically this is $set or $unset, but others are available.
type Operation struct {
	Name   string
	Values Properties
}

type BatchEvent struct {
	DistinctID string     `json:"-"`
	Event      string     `json:"event"`
	Props      Properties `json:"properties"`
}

// New returns a configured client.
func New(token string) *Client {
	return newWithTransport(token, nil)
}

func newWithTransport(token string, transport http.RoundTripper) *Client {
	return &Client{
		Token:   token,
		BaseURL: apiBaseURL,
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

var MaxBatchSizeErr = errors.New("max batch size is 50")

func (m *Client) TrackBatch(events []BatchEvent) error {
	if len(events) > 50 {
		return MaxBatchSizeErr
	}

	for _, e := range events {
		if e.DistinctID != "" {
			e.Props["distinct_id"] = e.DistinctID
		}
		e.Props["token"] = m.Token
		e.Props["mp_lib"] = library
	}
	return m.makeRequestWithData("POST", "track", events, sourceUser)
}

// Track sends event data with optional metadata.
func (m *Client) Track(distinctID string, event string, props Properties) error {
	if distinctID != "" {
		props["distinct_id"] = distinctID
	}
	props["token"] = m.Token
	props["mp_lib"] = library

	data := map[string]interface{}{"event": event, "properties": props}
	return m.makeRequestWithData("GET", "track", data, sourceUser)
}

// TrackAsScript sends event data on behalf of a specific user, without using the current
// IP address for location determination.
func (m *Client) TrackAsScript(distinctID string, event string, props Properties) error {
	if distinctID == "" {
		return errors.New("No distinct_id specified.")
	}
	props["token"] = m.Token
	props["mp_lib"] = library

	data := map[string]interface{}{"event": event, "properties": props}
	return m.makeRequestWithData("GET", "track", data, sourceScript)
}

// Engage updates profile data.
// This will update the IP and related data on the profile.
// If you don't have the IP address of the user, then use the UpdateProperties method instead,
// otherwise the user's location will be set to wherever the script was run from.
func (m *Client) Engage(distinctID string, props Properties, op *Operation) error {
	return m.engage(distinctID, props, op, sourceUser)
}

// EngageAsScript calls the engage endpoint, but doesn't set IP, city, country, on the profile.
func (m *Client) EngageAsScript(distinctID string, props Properties, op *Operation) error {
	return m.engage(distinctID, props, op, sourceScript)
}

func (m *Client) engage(distinctID string, props Properties, op *Operation, as actionSource) error {
	if distinctID != "" {
		props["$distinct_id"] = distinctID
	}
	props["$token"] = m.Token
	props["mp_lib"] = library
	if op.Name == "$unset" {
		keys := []interface{}{}
		for key := range op.Values {
			keys = append(keys, key)
		}
		props[op.Name] = keys
	} else {
		props[op.Name] = op.Values
	}

	return m.makeRequestWithData("GET", "engage", props, as)
}

// TrackingPixel returns a url that, when clicked, will track the given data and then redirect to provided url.
func (m *Client) TrackingPixel(distinctID, event string, props Properties) (string, error) {
	if distinctID != "" {
		props["$distinct_id"] = distinctID
	}
	props["$token"] = m.Token
	props["mp_lib"] = library

	data := map[string]interface{}{"event": event, "properties": props}
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"data": base64.StdEncoding.EncodeToString(json),
		"img":  "1",
	}
	query := url.Values{}
	for k, v := range params {
		query[k] = []string{v}
	}
	return fmt.Sprintf("%s/%s?%s", m.BaseURL, "track", query.Encode()), nil
}

// RedirectURL returns a url that, when clicked, will track the given data and then redirect to provided url.
func (m *Client) RedirectURL(distinctID, event, uri string, props Properties) (string, error) {
	if distinctID != "" {
		props["$distinct_id"] = distinctID
	}
	props["$token"] = m.Token
	props["mp_lib"] = library

	data := map[string]interface{}{"event": event, "properties": props}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"data":     base64.StdEncoding.EncodeToString(b),
		"redirect": uri,
	}
	query := url.Values{}
	for k, v := range params {
		query[k] = []string{v}
	}
	return fmt.Sprintf("%s/%s?%s", m.BaseURL, "track", query.Encode()), nil
}

func (m *Client) makeRequest(method string, endpoint string, paramMap map[string]string) error {
	var (
		err error
		req *http.Request
		r   io.Reader
	)

	if endpoint == "" {
		return fmt.Errorf("endpoint missing")
	}

	endpoint = fmt.Sprintf("%s/%s", m.BaseURL, endpoint)

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

	if method == "POST" {
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// The API documentation states that success will be reported with either "1" or "1\n".
	if strings.Trim(string(b), "\n") != "1" {
		return fmt.Errorf("request failed - %s", b)
	}
	return nil
}

func (m *Client) makeRequestWithData(method string, endpoint string, data interface{}, as actionSource) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	params := map[string]string{
		"data": base64.StdEncoding.EncodeToString(b),
	}

	if as == sourceScript {
		params["ip"] = "0"
	}

	return m.makeRequest(method, endpoint, params)
}
