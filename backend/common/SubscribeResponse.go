package common

type SubscribeResponse struct {
	Successes []string `json:"successes"`
	Errors    []string `json:"errors"`
}
