package messages

type Anon struct {
	ChannelId string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Text string `json:"text"`
	ResponseUrl string `json:"response_url"`
	TriggerId string `json:"trigger_id"`
}
