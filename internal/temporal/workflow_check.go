package temporal

type WorkflowCheckParam struct {
	Script         string
	Filter         string
	CheckId        string
	WorkerGroupIds []string
}

const WorkflowCheckName = "CHECK_WORKFLOW"
