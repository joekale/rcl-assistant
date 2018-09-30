package messages

type Event struct {
	Token string `json:"token"`
	TeamId string `json:"team_id"`
	ApiAppId string `json:"api_app_id"`
	EventInfo []byte `json:"event"`
	EventType string `json:"type"`
	AuthedUsers []string `json:"authed_users"`
	EventId string `json:"event_id"`
	EventTime int64 `json:"event_time"`
}

type Type struct {
	Type string `json:"type"`
}

type Profile struct {
	StatusText string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
	StatusExpiration string `json:"status_expiration"`
	RealName string `json:"real_name"`
	DisplayName string `json:"display_name"`
	RealNameNormalized string `json:"real_name_normalized"`
	DisplayNameNormalized string `json:"display_name_normalized"`
	Team string `json:"team"`
}

type User struct {
	UserId string `json:"type"`
	TeamId string `json:"team_id"`
	Name string `json:"name"`
	Deleted bool `json:"deleted"`
	RealName string `json:"real_name"`
	TZ string `json:"tz"`
	IsAdmin bool `json:"is_admin"`
	IsOwner bool `json:"is_owner"`
	IsPrimaryOwner bool `json:"is_primary_owner"`
	IsBot bool `json:"is_bot"`
	IsStranger bool `json:"is_stranger"`
	Updated int64 `json:"updated"`
	IsAppUser bool `json:"is_app_user"`
	profile Profile
}

type TeamJoin struct {
	event_type Type
	user User
}
