package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/web/templates/components"
)

type IndexData struct {
	*components.Base
	Targets       map[string]TargetsAndStatus
	TargetsLength int
	TimeRange     string
	Status        models.TargetStatus
}

type Target struct {
	Name       string
	Visibility models.TargetVisibility
	Group      string
	Status     models.TargetStatus
	History    *History
}

type HistoryItem struct {
	Status models.TargetStatus
	Date   time.Time
}

type History struct {
	List   []HistoryItem
	Uptime float64
}

type TargetsAndStatus struct {
	Status  models.TargetStatus
	Targets []*Target
}

func getDateString(date time.Time) string {
	return date.UTC().Format("2006-01-02T15:04:05")
}

func getHistory(history []*services.TargetHistory, period time.Duration, buckets int) *History {
	historyMap := map[string]models.TargetStatus{}
	numOfSuccess := 0.0
	numTotal := 0.0

	for i := 0; i < buckets; i++ {
		dateString := getDateString(time.Now().Add(period * time.Duration(-i)).Truncate(period))
		historyMap[dateString] = models.TargetStatusUnknown
	}

	for _, _history := range history {
		dateString := getDateString(_history.CreatedAt.Time.Truncate(period))

		// Skip if not part of the "buckets"
		if _, ok := historyMap[dateString]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.TargetStatusSuccess {
			numOfSuccess++
		}

		// skip if it is already set to failure
		if historyMap[dateString] == models.TargetStatusFailure {
			continue
		}

		// FIXME: This is wrong! As we can have multiple targets in dateString.
		// We should look at only the newest one.
		historyMap[dateString] = _history.Status
	}

	historyItems := make([]HistoryItem, buckets)
	for i := 0; i < buckets; i++ {
		date := time.Now().Add(period * time.Duration(-buckets+i+1)).Truncate(period)
		datestring := getDateString(date)
		historyItems[i] = HistoryItem{
			Status: historyMap[datestring],
			Date:   date,
		}
	}

	uptime := 0.0
	if numTotal > 0 {
		uptime = 100.0 * numOfSuccess / numTotal
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

	timeRange := c.QueryParam("time-range")
	if timeRange != "48hours" && timeRange != "90days" && timeRange != "90minutes" {
		timeRange = "90days"
	}

	overallStatus := models.TargetStatusUnknown
	statusByGroup := make(map[string]models.TargetStatus)

	targetsWithHistory := make([]*Target, len(targets))
	for i, target := range targets {
		history, err := services.GetTargetHistoryForTarget(ctx, h.db, target.Id)
		if err != nil {
			return err
		}

		var historyResult *History
		switch timeRange {
		case "48hours":
			historyResult = getHistory(history, time.Hour, 48)
		case "90days":
			historyResult = getHistory(history, time.Hour*24, 90)
		case "90minutes":
			historyResult = getHistory(history, time.Minute, 90)
		}

		if statusByGroup[target.Group] == "" {
			statusByGroup[target.Group] = models.TargetStatusUnknown
		}

		status := historyResult.List[len(historyResult.List)-1]
		if status.Status == models.TargetStatusSuccess {
			if overallStatus == models.TargetStatusUnknown {
				overallStatus = status.Status
			}
			if statusByGroup[target.Group] == models.TargetStatusUnknown {
				statusByGroup[target.Group] = status.Status
			}
		}
		if status.Status != models.TargetStatusSuccess && status.Status != models.TargetStatusUnknown {
			overallStatus = status.Status
			statusByGroup[target.Group] = status.Status
		}

		targetsWithHistory[i] = &Target{
			Name:       target.Name,
			Visibility: target.Visibility,
			Group:      target.Group,
			Status:     status.Status,
			History:    historyResult,
		}
	}

	targetsByGroup := map[string]TargetsAndStatus{}
	for _, target := range targetsWithHistory {
		targetsByGroup[target.Group] = TargetsAndStatus{
			Status:  statusByGroup[target.Group],
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
		TimeRange: timeRange,
		Status:    overallStatus,
	})
}
