package mixpanel

import "testing"

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
