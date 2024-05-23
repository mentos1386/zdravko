package api

type CheckStatus string

const (
	CheckStatusSuccess CheckStatus = "SUCCESS"
	CheckStatusFailure CheckStatus = "FAILURE"
	CheckStatusUnknown CheckStatus = "UNKNOWN"
)

type ApiV1ChecksHistoryPOSTBody struct {
	Status        CheckStatus `json:"status"`
	Note          string      `json:"note"`
	WorkerGroupId string      `json:"worker_group"`
}
