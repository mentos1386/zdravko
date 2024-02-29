package api

import "code.tjo.space/mentos1386/zdravko/database/models"

type ApiV1MonitorsHistoryPOSTBody struct {
	Status        models.MonitorStatus `json:"status"`
	Note          string               `json:"note"`
	WorkerGroupId string               `json:"worker_group"`
}
