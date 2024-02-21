package api

type ApiV1HealthchecksHistoryPOSTBody struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}
