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
	Monitors       map[string]MonitorsAndStatus
	MonitorsLength int
	TimeRange      string
	Status         models.MonitorStatus
}

type Monitor struct {
	Name    string
	Group   string
	Status  models.MonitorStatus
	History *History
}

type HistoryItem struct {
	Status models.MonitorStatus
	Date   time.Time
}

type History struct {
	List   []HistoryItem
	Uptime int
}

type MonitorsAndStatus struct {
	Status   models.MonitorStatus
	Monitors []*Monitor
}

func getDateString(date time.Time) string {
	return date.UTC().Format("2006-01-02T15:04:05")
}

func getHistory(history []*models.MonitorHistory, period time.Duration, buckets int) *History {
	historyMap := map[string]models.MonitorStatus{}
	numOfSuccess := 0
	numTotal := 0

	for i := 0; i < buckets; i++ {
		dateString := getDateString(time.Now().Add(period * time.Duration(-i)).Truncate(period))
		historyMap[dateString] = models.MonitorUnknown
	}

	for _, _history := range history {
		dateString := getDateString(_history.CreatedAt.Time.Truncate(period))

		// Skip if not part of the "buckets"
		if _, ok := historyMap[dateString]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.MonitorSuccess {
			numOfSuccess++
		}

		// skip if it is already set to failure
		if historyMap[dateString] == models.MonitorFailure {
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

	uptime := 0
	if numTotal > 0 {
		uptime = 100 * numOfSuccess / numTotal
	}

	return &History{
		List:   historyItems,
		Uptime: uptime,
	}
}

func (h *BaseHandler) Index(c echo.Context) error {
	ctx := context.Background()
	monitors, err := services.GetMonitors(ctx, h.db)
	if err != nil {
		return err
	}

	timeRange := c.QueryParam("time-range")
	if timeRange != "48hours" && timeRange != "90days" && timeRange != "90minutes" {
		timeRange = "90days"
	}

	overallStatus := models.MonitorUnknown
	statusByGroup := make(map[string]models.MonitorStatus)

	monitorsWithHistory := make([]*Monitor, len(monitors))
	for i, monitor := range monitors {
		history, err := services.GetMonitorHistoryForMonitor(ctx, h.db, monitor.Id)
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

		if statusByGroup[monitor.Group] == "" {
			statusByGroup[monitor.Group] = models.MonitorUnknown
		}

		status := historyResult.List[len(historyResult.List)-1]
		if status.Status == models.MonitorSuccess {
			if overallStatus == models.MonitorUnknown {
				overallStatus = status.Status
			}
			if statusByGroup[monitor.Group] == models.MonitorUnknown {
				statusByGroup[monitor.Group] = status.Status
			}
		}
		if status.Status != models.MonitorSuccess && status.Status != models.MonitorUnknown {
			overallStatus = status.Status
			statusByGroup[monitor.Group] = status.Status
		}

		monitorsWithHistory[i] = &Monitor{
			Name:    monitor.Name,
			Group:   monitor.Group,
			Status:  status.Status,
			History: historyResult,
		}
	}

	monitorsByGroup := map[string]MonitorsAndStatus{}
	for _, monitor := range monitorsWithHistory {
		monitorsByGroup[monitor.Group] = MonitorsAndStatus{
			Status:   statusByGroup[monitor.Group],
			Monitors: append(monitorsByGroup[monitor.Group].Monitors, monitor),
		}
	}

	c.Response().Header().Set("Cache-Control", "max-age=10")

	return c.Render(http.StatusOK, "index.tmpl", &IndexData{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Status"),
			Navbar:       Pages,
		},
		Monitors:  monitorsByGroup,
		TimeRange: timeRange,
		Status:    overallStatus,
	})
}
