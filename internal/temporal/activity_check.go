package temporal

import "github.com/mentos1386/zdravko/database/models"

type ActivityCheckParam struct {
	Script string
	Target *models.Target
}

type ActivityCheckResult struct {
	Success bool
	Note    string
}

const ActivityCheckName = "CHECK"
