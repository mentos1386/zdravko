package handlers

import (
	"context"
	"net/http"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type IndexData struct {
	*components.Base
	Checks       map[string]ChecksAndStatus
	ChecksLength int
	TimeRange    string
	Status       models.CheckStatus
}

type Check struct {
	Name    string
	Group   string
	Status  models.CheckStatus
	History *History
}

type HistoryItem struct {
	Status models.CheckStatus
	Date   time.Time
}

type History struct {
	List   []HistoryItem
	Uptime float64
}

type ChecksAndStatus struct {
	Status models.CheckStatus
	Checks []*Check
}

func getDateString(date time.Time) string {
	return date.UTC().Format("2006-01-02T15:04:05")
}

func getHistory(history []*models.CheckHistory, period time.Duration, buckets int) *History {
	historyMap := map[string]models.CheckStatus{}
	numOfSuccess := 0.0
	numTotal := 0.0

	for i := 0; i < buckets; i++ {
		dateString := getDateString(time.Now().Add(period * time.Duration(-i)).Truncate(period))
		historyMap[dateString] = models.CheckUnknown
	}

	for _, _history := range history {
		dateString := getDateString(_history.CreatedAt.Time.Truncate(period))

		// Skip if not part of the "buckets"
		if _, ok := historyMap[dateString]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.CheckSuccess {
			numOfSuccess++
		}

		// skip if it is already set to failure
		if historyMap[dateString] == models.CheckFailure {
			continue
		}

		// FIXME: This is wrong! As we can have multiple checks in dateString.
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
	checks, err := services.GetChecks(ctx, h.db)
	if err != nil {
		return err
	}

	timeRange := c.QueryParam("time-range")
	if timeRange != "48hours" && timeRange != "90days" && timeRange != "90minutes" {
		timeRange = "90days"
	}

	overallStatus := models.CheckUnknown
	statusByGroup := make(map[string]models.CheckStatus)

	checksWithHistory := make([]*Check, len(checks))
	for i, check := range checks {
		history, err := services.GetCheckHistoryForCheck(ctx, h.db, check.Id)
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

		if statusByGroup[check.Group] == "" {
			statusByGroup[check.Group] = models.CheckUnknown
		}

		status := historyResult.List[len(historyResult.List)-1]
		if status.Status == models.CheckSuccess {
			if overallStatus == models.CheckUnknown {
				overallStatus = status.Status
			}
			if statusByGroup[check.Group] == models.CheckUnknown {
				statusByGroup[check.Group] = status.Status
			}
		}
		if status.Status != models.CheckSuccess && status.Status != models.CheckUnknown {
			overallStatus = status.Status
			statusByGroup[check.Group] = status.Status
		}

		checksWithHistory[i] = &Check{
			Name:    check.Name,
			Group:   check.Group,
			Status:  status.Status,
			History: historyResult,
		}
	}

	checksByGroup := map[string]ChecksAndStatus{}
	for _, check := range checksWithHistory {
		checksByGroup[check.Group] = ChecksAndStatus{
			Status: statusByGroup[check.Group],
			Checks: append(checksByGroup[check.Group].Checks, check),
		}
	}

	c.Response().Header().Set("Cache-Control", "max-age=10")

	return c.Render(http.StatusOK, "index.tmpl", &IndexData{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Status"),
			Navbar:       Pages,
		},
		Checks:    checksByGroup,
		TimeRange: timeRange,
		Status:    overallStatus,
	})
}
