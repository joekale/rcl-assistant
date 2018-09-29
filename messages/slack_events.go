package messages

type Event struct {
	Token string `json:"token"`
	TeamId string `json:"team_id"`
	ApiAppId string `json:"api_app_id"`
	EventInfo interface{} `json:"event"`
	EventType string `json:"type"`
	AuthedUsers []string `json:"authed_users"`
	EventId string `json:"event_id"`
	EventTime int64 `json:"event_time"`
}
