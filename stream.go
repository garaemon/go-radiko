package radiko

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

// URLItem represents a stream url.
type URLItem struct {
	Areafree bool   `xml:"areafree,attr"`
	Item     string `xml:",chardata"`
}

// GetStreamMultiURL returns a slice of the stream url.
func GetStreamMultiURL(stationID string) ([]URLItem, error) {
	endpoint := path.Join(apiV2, "station/stream_multi",
		fmt.Sprintf("%s.xml", stationID))

	resp, err := http.Get("http://radiko.jp/" + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	urlData := streamURLData{}
	if err = xml.Unmarshal(b, &urlData); err != nil {
		return nil, err
	}
	return urlData.Items, err
}

type streamURLData struct {
	XMLName xml.Name  `xml:"url"`
	Items   []URLItem `xml:"item"`
}

// GetLiveURL returns a live url for web browser.
func GetLiveURL(stationID string) string {
	endpoint := path.Join("#!/live", stationID)
	return defaultEndpoint + "/" + endpoint
}

// URLItem represents a stream url.
type SmhURLItem struct {
	Text              string `xml:",chardata"`
	Areafree          bool   `xml:"areafree,attr"`
	MediaURLPath      string `xml:"media_url_path"`
	PlaylistCreateURL string `xml:"playlist_create_url"`
	PlaylistURLPath   string `xml:"playlist_url_path"`
}

type streamSmhURLData struct {
	XMLName xml.Name     `xml:"urls"`
	Text    string       `xml:",chardata"`
	URL     []SmhURLItem `xml:"url"`
}

// GetStreamSmhMultiURL returns a slice of the stream smh url.
func GetStreamSmhMultiURL(stationID string) ([]SmhURLItem, error) {
	endpoint := path.Join(apiV2, "station/stream_smh_multi",
		fmt.Sprintf("%s.xml", stationID))

	resp, err := http.Get("http://radiko.jp/" + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	urlData := streamSmhURLData{}
	if err = xml.Unmarshal(b, &urlData); err != nil {
		return nil, err
	}
	return urlData.URL, err
}

// StreamPcHtml5URLItem represents a stream url from pc_html5 API.
type StreamPcHtml5URLItem struct {
	Areafree          bool   `xml:"areafree,attr"`
	Timefree          bool   `xml:"timefree,attr"`
	PlaylistCreateURL string `xml:"playlist_create_url"`
}

type streamPcHtml5URLData struct {
	XMLName xml.Name               `xml:"urls"`
	URLs    []StreamPcHtml5URLItem `xml:"url"`
}

// GetStreamPcHtml5URL returns a slice of the pc_html5 stream url.
func GetStreamPcHtml5URL(stationID string) ([]StreamPcHtml5URLItem, error) {
	endpoint := path.Join(apiV3, "station/stream/pc_html5",
		fmt.Sprintf("%s.xml", stationID))

	resp, err := http.Get(defaultEndpoint + "/" + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	urlData := streamPcHtml5URLData{}
	if err = xml.Unmarshal(b, &urlData); err != nil {
		return nil, err
	}
	return urlData.URLs, err
}

// GetTimefreePlaylistBaseURL returns the first timefree playlist base URL.
func GetTimefreePlaylistBaseURL(stationID string, areafree bool) (string, error) {
	urls, err := GetStreamPcHtml5URL(stationID)
	if err != nil {
		return "", err
	}

	for _, u := range urls {
		if u.Timefree && u.Areafree == areafree && u.PlaylistCreateURL != "" {
			return u.PlaylistCreateURL, nil
		}
	}
	return "", fmt.Errorf("timefree playlist URL not found for station %s", stationID)
}
