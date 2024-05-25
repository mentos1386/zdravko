package temporal

type ActivityTargetsFilterParam struct {
	Filter string
}

type ActivityTargetsFilterResult struct {
	Targets []*Target
}

const ActivityTargetsFilterName = "TARGETS_FILTER"
