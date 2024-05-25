package temporal

type ActivityCheckParam struct {
	Script string
	Target *Target
}

type ActivityCheckResult struct {
	Success bool
	Note    string
}

const ActivityCheckName = "CHECK"
