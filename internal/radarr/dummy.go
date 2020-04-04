// Package radarr here only exist for testing
package radarr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var dummyMovieResponse string = `
{
  "title": "Frozen II",
  "alternativeTitles": [
    {
      "sourceType": "tmdb",
      "movieId": 217,
      "title": "Frozen 2",
      "sourceId": 330457,
      "votes": 0,
      "voteCount": 0,
      "language": {
        "id": 1,
        "name": "English"
      },
      "id": 461
    }
  ],
  "secondaryYearSourceId": 0,
  "sortTitle": "frozen ii",
  "sizeOnDisk": 4099483594,
  "status": "released",
  "overview": "Elsa, Anna, Kristoff and Olaf head far into the forest to learn the truth about an ancient mystery of their kingdom.",
  "inCinemas": "2019-11-19T23:00:00Z",
  "physicalRelease": "2020-02-11T00:00:00Z",
  "images": [
    {
      "coverType": "poster",
      "url": "/MediaCover/217/poster.jpg?lastWrite=637214530603577317"
    },
    {
      "coverType": "fanart",
      "url": "/MediaCover/217/fanart.jpg?lastWrite=637202450497927734"
    }
  ],
  "website": "https://movies.disney.com/frozen-2",
  "downloaded": true,
  "year": 2019,
  "hasFile": true,
  "youTubeTrailerId": "Zi4LMpSDccc",
  "studio": "Walt Disney Animation Studios",
  "path": "/movies/Frozen II (2019)",
  "profileId": 3,
  "monitored": true,
  "minimumAvailability": "inCinemas",
  "isAvailable": true,
  "folderName": "/movies/Frozen II (2019)",
  "runtime": 104,
  "lastInfoSync": "2020-04-03T19:39:49.6265379Z",
  "cleanTitle": "frozenii",
  "imdbId": "tt4520988",
  "tmdbId": 330457,
  "titleSlug": "frozen-ii-330457",
  "genres": ["Animation", "Family", "Adventure"],
  "tags": [1],
  "added": "2020-03-15T15:39:15.8796553Z",
  "ratings": {
    "votes": 3481,
    "value": 7.1
  },
  "movieFile": {
    "movieId": 0,
    "relativePath": "Frozen.2.2019.MULTi.1080p.WEB.x264.EXTREME.mkv",
    "size": 4099483594,
    "dateAdded": "2020-03-15T16:18:06.9156804Z",
    "sceneName": "Frozen.2.2019.MULTi.1080p.WEB.x264.EXTREME",
    "quality": {
      "quality": {
        "id": 3,
        "name": "WEBDL-1080p",
        "source": "webdl",
        "resolution": 1080,
        "modifier": "none"
      },
      "revision": {
        "version": 1,
        "real": 0,
        "isRepack": false
      }
    },
    "edition": "",
    "mediaInfo": {
      "containerFormat": "Matroska",
      "videoFormat": "AVC",
      "videoCodecID": "V_MPEG4/ISO/AVC",
      "videoProfile": "High@L4",
      "videoCodecLibrary": "",
      "videoBitrate": 4526004,
      "videoBitDepth": 8,
      "videoMultiViewCount": 0,
      "videoColourPrimaries": "BT.709",
      "videoTransferCharacteristics": "BT.709",
      "width": 1920,
      "height": 804,
      "audioFormat": "AC-3",
      "audioCodecID": "A_AC3",
      "audioCodecLibrary": "",
      "audioAdditionalFeatures": "",
      "audioBitrate": 384000,
      "runTime": "01:43:12.3530000",
      "audioStreamCount": 2,
      "audioChannels": 6,
      "audioChannelPositions": "3/2/0.1",
      "audioChannelPositionsText": "Front: L C R, Side: L R, LFE",
      "audioProfile": "",
      "videoFps": 23.976,
      "audioLanguages": "French / English",
      "subtitles": "French",
      "scanType": "Progressive",
      "schemaRevision": 5
    },
    "id": 197
  },
  "qualityProfileId": 3,
  "id": 217
}`

var dummySystemStatusResponse string = `
{
  "version": "3.0.0.2741",
  "buildTime": "2020-03-23T16:23:16Z",
  "isDebug": false,
  "isProduction": true,
  "isAdmin": false,
  "isUserInteractive": false,
  "startupPath": "/opt/radarr",
  "appData": "/config",
  "osName": "ubuntu",
  "osVersion": "20.04",
  "isNetCore": true,
  "isMono": false,
  "isLinux": true,
  "isOsx": false,
  "isWindows": false,
  "branch": "develop",
  "authentication": "forms",
  "sqliteVersion": "3.31.1",
  "migrationVersion": 169,
  "urlBase": "",
  "runtimeVersion": "3.1.2",
  "runtimeName": "netCore"
}`

var dummyMoviesResponse string = fmt.Sprintf("[%s, %s]", dummyMovieResponse, dummyMovieResponse)

// DummyUnauthorizedResponse describe Unauthorized Radarr response
var DummyUnauthorizedResponse string = `{"error": "Unauthorized"}`

// DummyNotFoundResponse describe NoFound Radarr response
var DummyNotFoundResponse string = `{"message": "NotFound"}`

var (
	// DummyHTTPClient mocked http client
	DummyHTTPClient *HTTPClient

	// DummyURL dummy Radarr instance URL
	DummyURL string = "https://radarr.dummy"

	// DummyAPIKey dummy Radarr API keys
	DummyAPIKey string = "dummy-api-key"
)

func init() {
	// Create a mock http client
	DummyHTTPClient = &HTTPClient{}
}

// TestCase describe a test case
type TestCase struct {
	Title    string
	Expected interface{}
	Got      interface{}
}

// HTTPClient implements HTTPClientInterface
type HTTPClient struct{}

// Get mock GET requests
func (c *HTTPClient) Get(targetURL string) (resp *http.Response, err error) {
	// Test valid API key
	t, _ := url.Parse(targetURL)
	params, _ := url.ParseQuery(t.RawQuery)
	key := params.Get("apikey")

	if key != DummyAPIKey {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Status:     http.StatusText(http.StatusUnauthorized),
			Body:       ioutil.NopCloser(bytes.NewBufferString(DummyUnauthorizedResponse)),
		}, nil
	}

	// Mock GET /movie
	if strings.Contains(targetURL, "/movie") {

		switch targetURL {
		case fmt.Sprintf("%s/api%s/%d?apikey=%s", DummyURL, "/movie", 217, DummyAPIKey):
			// Get one movie
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       ioutil.NopCloser(bytes.NewBufferString(dummyMovieResponse)),
			}, nil

		case fmt.Sprintf("%s/api%s?apikey=%s", DummyURL, "/movie", DummyAPIKey):
			// List of movies
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       ioutil.NopCloser(bytes.NewBufferString(dummyMoviesResponse)),
			}, nil

		default:
			// Defaulting to 404
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     http.StatusText(http.StatusNotFound),
				Body:       ioutil.NopCloser(bytes.NewBufferString(DummyNotFoundResponse)),
			}, nil
		}
	}

	// Mock GET /system/status
	if strings.Contains(targetURL, "/system/status") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Body:       ioutil.NopCloser(bytes.NewBufferString(dummySystemStatusResponse)),
		}, nil
	}

	return &http.Response{
		StatusCode: http.StatusNotFound,
		Status:     http.StatusText(http.StatusNotFound),
		Body:       ioutil.NopCloser(bytes.NewBufferString(DummyNotFoundResponse)),
	}, nil
}
