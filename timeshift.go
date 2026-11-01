package radiko

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/url"
	"path"
	"time"

	"github.com/yyoshiki41/go-radiko/internal/util"
)

// TimeshiftPlaylistInfo holds the information needed for timefree download.
type TimeshiftPlaylistInfo struct {
	BaseURL   string
	StationID string
	Ft        string
	To        string
	Lsid      string
}

// ChunkURL builds a playlist URL for a specific time range chunk.
func (info *TimeshiftPlaylistInfo) ChunkURL(seek, endAt string, l int) string {
	params := url.Values{}
	params.Set("station_id", info.StationID)
	params.Set("start_at", info.Ft)
	params.Set("ft", info.Ft)
	params.Set("seek", seek)
	params.Set("end_at", endAt)
	params.Set("to", endAt)
	params.Set("l", fmt.Sprintf("%d", l))
	params.Set("lsid", info.Lsid)
	params.Set("type", "c")

	return info.BaseURL + "?" + params.Encode()
}

// TimeshiftPlaylistM3U8 returns a playlist URL for ffmpeg-based downloading.
func (c *Client) TimeshiftPlaylistM3U8(ctx context.Context, stationID string, start time.Time) (string, error) {
	info, err := c.GetTimeshiftPlaylistInfo(ctx, stationID, start, false)
	if err != nil {
		return "", err
	}

	return info.ChunkURL(info.Ft, info.To, 15), nil
}

// GetTimeshiftPlaylistInfo returns the information needed for timefree chunk download.
func (c *Client) GetTimeshiftPlaylistInfo(ctx context.Context, stationID string, start time.Time, areafree bool) (*TimeshiftPlaylistInfo, error) {
	prog, err := c.GetProgramByStartTime(ctx, stationID, start)
	if err != nil {
		return nil, err
	}

	baseURL, err := GetTimefreePlaylistBaseURL(stationID, areafree)
	if err != nil {
		return nil, err
	}

	return &TimeshiftPlaylistInfo{
		BaseURL:   baseURL,
		StationID: stationID,
		Ft:        prog.Ft,
		To:        prog.To,
		Lsid:      GenerateLsid(),
	}, nil
}

// GenerateLsid returns a random MD5 hash string used as a session ID.
func GenerateLsid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%f", r.Float64()))))
}

// GetTimeshiftURL returns a timeshift url for web browser.
func GetTimeshiftURL(stationID string, start time.Time) string {
	endpoint := path.Join("#!/ts", stationID, util.Datetime(start))
	return defaultEndpoint + "/" + endpoint
}
