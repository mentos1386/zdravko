package temporal

type AddTargetHistoryStatus string

const (
	AddTargetHistoryStatusSuccess AddTargetHistoryStatus = "SUCCESS"
	AddTargetHistoryStatusFailure AddTargetHistoryStatus = "FAILURE"
	AddTargetHistoryStatusUnknown AddTargetHistoryStatus = "UNKNOWN"
)

type ActivityAddTargetHistoryParam struct {
	Target *Target
	Status AddTargetHistoryStatus
	Note   string
}

type ActivityAddTargetHistoryResult struct {
}

const ActivityAddTargetHistoryName = "ADD_TARGET_HISTORY"
