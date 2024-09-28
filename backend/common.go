package backend

import (
	"github.com/goccy/go-json"
	"time"
)

type User struct {
	Login string
	Name  string
	Email string
	Role  string
}

type AppInstanceSettings struct {
	JSONData json.RawMessage

	DecryptedSecureJSONData map[string]string

	Updated time.Time
}

type DataSourceInstanceSettings struct {
	ID int64

	UID string

	Type string

	Name string

	URL string

	User string

	Database string

	BasicAuthEnabled bool

	BasicAuthUser string

	JSONData json.RawMessage

	DecryptedSecureJSONData map[string]string

	Updated time.Time
}

type PluginContext struct {
	OrgID int64

	PluginID string

	User *User

	AppInstanceSettings *AppInstanceSettings

	DataSourceInstanceSettings *DataSourceInstanceSettings
}
