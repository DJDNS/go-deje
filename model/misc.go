package model

type JSONObject map[string]interface{}
type DownloadURLs []string

type Manageable interface {
	GetKey() string
	GetGroupKey() string

	Eq(Manageable) bool
}
