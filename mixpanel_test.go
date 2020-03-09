package mixpanel

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestRunningAsScriptWithoutDistinctID(t *testing.T) {
	mp := New("abc")

	err := mp.TrackAsScript("", "Some Event", nil)
	if err == nil {
		t.Errorf("Error expected here, but could successfully send event as as script without distinct_id")
	}
}

func TestRedirectURL(t *testing.T) {
	mp := New("abc")
	props := Properties{
		"a": "apples",
		"b": "bananas",
		"c": "cherries",
	}

	actualURI, err := mp.RedirectURL("12345", "Clicked through", "http://example.com", props)
	if err != nil {
		t.Fatal(err)
	}

	// data decodes to:
	// "{\"event\":\"Clicked through\",\"properties\":{\"$distinct_id\":\"12345\",\"$token\":\"abc\",\"a\":\"apples\",\"b\":\"bananas\",\"c\":\"cherries\",\"mp_lib\":\"timehop/go-mixpanel\"}}"
	expectedURI := `https://api.mixpanel.com/track?data=eyJldmVudCI6IkNsaWNrZWQgdGhyb3VnaCIsInByb3BlcnRpZXMiOnsiJGRpc3RpbmN0X2lkIjoiMTIzNDUiLCIkdG9rZW4iOiJhYmMiLCJhIjoiYXBwbGVzIiwiYiI6ImJhbmFuYXMiLCJjIjoiY2hlcnJpZXMiLCJtcF9saWIiOiJ2aXp6bG8vbWl4cGFuZWwifX0%3D&redirect=http%3A%2F%2Fexample.com`

	if actualURI != expectedURI {
		t.Errorf("\n got: %s\nwant: %s\n", actualURI, expectedURI)
	}
}

func TestTrackingPixel(t *testing.T) {
	mp := New("abc")
	props := Properties{
		"a": "apples",
		"b": "bananas",
		"c": "cherries",
	}

	actualURI, err := mp.TrackingPixel("12345", "Clicked through", props)
	if err != nil {
		t.Fatal(err)
	}

	// data decodes to:
	// "{\"event\":\"Clicked through\",\"properties\":{\"$distinct_id\":\"12345\",\"$token\":\"abc\",\"a\":\"apples\",\"b\":\"bananas\",\"c\":\"cherries\",\"mp_lib\":\"timehop/go-mixpanel\"}}"
	expectedURI := `https://api.mixpanel.com/track?data=eyJldmVudCI6IkNsaWNrZWQgdGhyb3VnaCIsInByb3BlcnRpZXMiOnsiJGRpc3RpbmN0X2lkIjoiMTIzNDUiLCIkdG9rZW4iOiJhYmMiLCJhIjoiYXBwbGVzIiwiYiI6ImJhbmFuYXMiLCJjIjoiY2hlcnJpZXMiLCJtcF9saWIiOiJ2aXp6bG8vbWl4cGFuZWwifX0%3D&img=1`

	if actualURI != expectedURI {
		t.Errorf("\n got: %s\nwant: %s\n", actualURI, expectedURI)
	}
}

type testTransport struct {
	f func(req *http.Request) *http.Response
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.f(req), nil
}

func TestTrackingBatch(t *testing.T) {
	testTransportFn := func(req *http.Request) *http.Response {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}

		form, err := url.ParseQuery(string(b))
		if err != nil {
			t.Fatal(err)
		}

		data := form["data"]
		if len(data) != 1 {
			t.Fatal("expect data to be present")
		}

		val := data[0]

		b, err = base64.StdEncoding.DecodeString(val)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson := `[{"event":"Clicked through","properties":{"a":"apples","b":"bananas","c":"cherries","distinct_id":"12345","mp_lib":"vizzlo/mixpanel","token":"abc"}}]`

		if string(b) != expectedJson {
			t.Errorf("expected %s to match %s", b, expectedJson)
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("1")),
			Header:     make(http.Header),
		}
	}

	mp := newWithTransport("abc", &testTransport{testTransportFn})
	props := Properties{
		"a": "apples",
		"b": "bananas",
		"c": "cherries",
	}

	err := mp.TrackBatch(
		[]BatchEvent{
			{"12345", "Clicked through", props},
		})

	if err != nil {
		t.Fatal(err)
	}

}
