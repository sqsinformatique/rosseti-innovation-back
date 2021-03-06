package models

import "github.com/sqsinformatique/rosseti-innovation-back/types"

// CentrifugoIntrospection contains an access token's session data as specified by Centrifugo documentation, see:
// https://centrifugal.github.io/centrifugo/server/proxy/#connect-proxy
// swagger:model centrifugoOAuth2TokenIntrospection
type CentrifugoIntrospection struct {
	// User is user ID (calculated on app backend based on request
	// cookie header for example). Return it as empty string for
	// accepting unauthenticated request
	//
	// required: true
	User string `json:"user"`
	// ExpireAt (optional integer) is a timestamp when connection
	// must be considered expired. If not set or set to 0 connection
	// won't expire at all
	ExpireAt int64 `json:"expire_at,omitempty"`
	//Info (optional JSON) is a connection info JSON
	Info interface{} `json:"info,omitempty"`
	// B64Info (optional string) is a binary connection info encoded
	// in base64 format, will be decoded to raw bytes on Centrifugo
	// before using in messages
	B64Info string `json:"b64info,omitempty"`
	// Data (optional JSON) is a custom data to send to client in connect
	// command response. Supported since v2.3.1
	Data interface{} `json:"data,omitempty"`
	// B64Data (optional string) is a custom data to send to client in
	// connect command response for binary connections, will be decoded
	// to raw bytes on Centrifugo side before sending to client.
	// Supported since v2.3.1
	B64Data string `json:"b64data,omitempty"`
	// Channels (optional array of strings) - allows to provide a list of
	// server-side channels to subscribe connection to. See more details
	// about server-side subscriptions (https://centrifugal.github.io/centrifugo/server/server_subs/).
	// Supported since v2.4.0
	Channels []string `json:"channels,omitempty"`
}

// CentrifugoIntrospectionResult contains a result of introspection access token
// swagger:model centrifugoOAuth2TokenIntrospectionReuslt
type CentrifugoIntrospectionResult struct {
	Result interface{} `json:"result"`
}

type CentrifugoParams struct {
	Channel string
	Data    interface{}
}

type CentrifugoAPIRequest struct {
	Method string
	Params *CentrifugoParams
}

type Theme struct {
	ID          int            `json:"id" db:"id"`
	Direction   int            `json:"direction" db:"direction"`
	Tags        string         `json:"tags" db:"tags"`
	Title       string         `json:"title" db:"title"`
	AuthorID    int            `json:"author_id" db:"author_id"`
	LikeCounter int            `json:"like_counter" db:"like_counter"`
	Meta        types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *Theme) SQLParamsRequest() []string {
	return []string{
		"direction",
		"title",
		"tags",
		"author_id",
		"like_counter",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type Direction struct {
	ID    int            `json:"id"`
	Title string         `json:"title"`
	Meta  types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *Direction) SQLParamsRequest() []string {
	return []string{
		"id",
		"title",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type DirectionDetailed struct {
	Themes []Theme `json:"themes"`
	Direction
}
