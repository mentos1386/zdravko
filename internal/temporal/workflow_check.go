package temporal

type WorkflowCheckParam struct {
	Script         string
	Filter         string
	CheckId        string
	WorkerGroupIds []string
}

type WorkflowCheckResult struct {
	Note string
}

const WorkflowCheckName = "CHECK_WORKFLOW"
