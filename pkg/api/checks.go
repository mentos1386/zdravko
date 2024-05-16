package api

import "code.tjo.space/mentos1386/zdravko/database/models"

type ApiV1ChecksHistoryPOSTBody struct {
	Status        models.CheckStatus `json:"status"`
	Note          string               `json:"note"`
	WorkerGroupId string               `json:"worker_group"`
}
