package api

type ApiV1MonitorsHistoryPOSTBody struct {
	Status      string `json:"status"`
	Note        string `json:"note"`
	WorkerGroup string `json:"worker_group"`
}
