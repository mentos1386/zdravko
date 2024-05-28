package handlers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/web/templates/components"
)

type HistoryOutcome string

const (
	HistoryOutcomeHealthy  HistoryOutcome = "HEALTHY"
	HistoryOutcomeDegraded HistoryOutcome = "DEGRADED"
	HistoryOutcomeUnknown  HistoryOutcome = "UNKNOWN"
	HistoryOutcomeDown     HistoryOutcome = "DOWN"
)

type IndexData struct {
	*components.Base
	Targets       map[string]TargetsAndStatus
	TargetsLength int
	TimeRange     string
	Outcome       HistoryOutcome
}

type Target struct {
	Name       string
	Visibility models.TargetVisibility
	Group      string
	Outcome    HistoryOutcome
	History    []*HistoryItem
	Uptime     float64
}

type History struct {
	List   []*HistoryItem
	Uptime float64
}

type HistoryItem struct {
	Outcome      HistoryOutcome
	StatusCounts map[models.TargetStatus]int
	Counts       int
	Date         time.Time
	Checks       []*HistoryItemCheck
	SuccessRate  float64
}

type HistoryItemCheck struct {
	Name            string
	WorkerGroupName string
	Status          models.TargetStatus
	StatusCounts    map[models.TargetStatus]int
	Counts          int
	SuccessRate     float64
}

type TargetsAndStatus struct {
	Outcome HistoryOutcome
	Targets []*Target
}

func getDateString(date time.Time) string {
	return date.UTC().Format("2006-01-02T15:04:05")
}

func TargetStatusToHistoryOutcome(status models.TargetStatus) HistoryOutcome {
	switch status {
	case models.TargetStatusSuccess:
		return HistoryOutcomeHealthy
	case models.TargetStatusFailure:
		return HistoryOutcomeDown
	default:
		return HistoryOutcomeUnknown
	}
}

func getHistory(history []*services.TargetHistory, period time.Duration, buckets int) *History {
	historyMap := map[string]*HistoryItem{}
	numOfSuccess := 0.0
	numTotal := 0.0

	mapKeys := make([]string, buckets)

	for i := 0; i < buckets; i++ {
		date := time.Now().Add(period * time.Duration(-i)).Truncate(period)
		dateString := getDateString(date)
		mapKeys[i] = dateString

		historyMap[dateString] = &HistoryItem{
			Outcome:      HistoryOutcomeUnknown,
			StatusCounts: map[models.TargetStatus]int{},
			Date:         date,
			Checks:       []*HistoryItemCheck{},
			SuccessRate:  0.0,
		}
	}

	for _, _history := range history {
		dateString := getDateString(_history.CreatedAt.Time.Truncate(period))

		entry, ok := historyMap[dateString]
		if !ok {
			continue
		}

		numTotal++
		if _history.Status == models.TargetStatusSuccess {
			numOfSuccess++
		}

		entry.StatusCounts[_history.Status]++
		entry.Counts++
		entry.SuccessRate = 100.0 * float64(entry.StatusCounts[models.TargetStatusSuccess]) / float64(entry.Counts)

		foundCheck := false
		for _, check := range entry.Checks {
			if check.Name == _history.CheckName && check.WorkerGroupName == _history.WorkerGroupName {
				foundCheck = true

				check.StatusCounts[_history.Status]++
				check.Counts++
				check.SuccessRate = 100.0 * float64(check.StatusCounts[models.TargetStatusSuccess]) / float64(check.Counts)

				if check.Status != models.TargetStatusFailure && _history.Status == models.TargetStatusFailure {
					check.Status = models.TargetStatusFailure
				}
			}
		}

		if !foundCheck {
			successRate := 0.0
			if _history.Status == models.TargetStatusSuccess {
				successRate = 100.0
			}
			entry.Checks = append(entry.Checks, &HistoryItemCheck{
				Name:            _history.CheckName,
				WorkerGroupName: _history.WorkerGroupName,
				Status:          _history.Status,
				StatusCounts: map[models.TargetStatus]int{
					_history.Status: 1,
				},
				Counts:      1,
				SuccessRate: successRate,
			})

			sort.Slice(entry.Checks, func(i, j int) bool {
				byName := entry.Checks[i].Name < entry.Checks[j].Name
				byWorkerGroupName := entry.Checks[i].WorkerGroupName < entry.Checks[j].WorkerGroupName
				return byName || (entry.Checks[i].Name == entry.Checks[j].Name && byWorkerGroupName)
			})
		}

		if entry.SuccessRate == 100.0 {
			entry.Outcome = HistoryOutcomeHealthy
		} else if entry.SuccessRate == 0.0 {
			entry.Outcome = HistoryOutcomeDown
		} else {
			entry.Outcome = HistoryOutcomeDegraded
		}

		historyMap[dateString] = entry
	}

	uptime := 0.0
	if numTotal > 0 {
		uptime = 100.0 * numOfSuccess / numTotal
	}

	historyItems := make([]*HistoryItem, 0, len(historyMap))
	for i := buckets - 1; i >= 0; i-- {
		key := mapKeys[i]
		historyItems = append(historyItems, historyMap[key])
	}

	return &History{
		List:   historyItems,
		Uptime: uptime,
	}
}

func (h *BaseHandler) Index(c echo.Context) error {
	ctx := context.Background()
	targets, err := services.GetTargets(ctx, h.db)
	if err != nil {
		return err
	}

	timeRangeQuery := c.QueryParam("time-range")
	if timeRangeQuery != "48hours" && timeRangeQuery != "60days" && timeRangeQuery != "60minutes" {
		timeRangeQuery = "60days"
	}

	var timeBuckets int
	var timeInterval time.Duration
	var timeRange services.TargetHistoryDateRange
	switch timeRangeQuery {
	case "60days":
		timeRange = services.TargetHistoryDateRange60Days
		timeInterval = time.Hour * 24
		timeBuckets = 60
	case "48hours":
		timeRange = services.TargetHistoryDateRange48Hours
		timeInterval = time.Hour
		timeBuckets = 48
	case "60minutes":
		timeRange = services.TargetHistoryDateRange60Minutes
		timeInterval = time.Minute
		timeBuckets = 60
	}

	overallOutcome := HistoryOutcomeUnknown
	outcomeByGroup := make(map[string]HistoryOutcome)

	targetsWithHistory := make([]*Target, len(targets))
	for i, target := range targets {
		history, err := services.GetTargetHistoryForTarget(ctx, h.db, target.Id, timeRange)
		if err != nil {
			return err
		}

		historyResult := getHistory(history, timeInterval, timeBuckets)

		if outcomeByGroup[target.Group] == "" {
			outcomeByGroup[target.Group] = HistoryOutcomeUnknown
		}

		status := historyResult.List[len(historyResult.List)-1]
		if status.Outcome == HistoryOutcomeHealthy {
			if overallOutcome == HistoryOutcomeUnknown {
				overallOutcome = status.Outcome
			}
			if outcomeByGroup[target.Group] == HistoryOutcomeUnknown {
				outcomeByGroup[target.Group] = status.Outcome
			}
		}
		if status.Outcome != HistoryOutcomeHealthy && status.Outcome != HistoryOutcomeUnknown {
			overallOutcome = status.Outcome
			outcomeByGroup[target.Group] = status.Outcome
		}

		targetsWithHistory[i] = &Target{
			Name:       target.Name,
			Visibility: target.Visibility,
			Group:      target.Group,
			Outcome:    status.Outcome,
			History:    historyResult.List,
			Uptime:     historyResult.Uptime,
		}
	}

	targetsByGroup := map[string]TargetsAndStatus{}
	for _, target := range targetsWithHistory {
		targetsByGroup[target.Group] = TargetsAndStatus{
			Outcome: outcomeByGroup[target.Group],
			Targets: append(targetsByGroup[target.Group].Targets, target),
		}
	}

	c.Response().Header().Set("Cache-Control", "max-age=10")

	return c.Render(http.StatusOK, "index.tmpl", &IndexData{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Status"),
			Navbar:       Pages,
		},
		Targets:   targetsByGroup,
		TimeRange: timeRangeQuery,
		Outcome:   overallOutcome,
	})
}
