package temporal

type WorkflowHookParam struct {
	Script string
	Filter string
	HookId string
}

type WorkflowHookResult struct {
	Note string
}

const WorkflowHookName = "HOOK_WORKFLOW"
