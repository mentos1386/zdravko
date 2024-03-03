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
	Monitors       map[string][]*Monitor
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

type History struct {
	List   []models.MonitorStatus
	Uptime int
}

func getHour(date time.Time) string {
	return date.UTC().Format("2006-01-02T15:04:05")
}

func getHistory(history []*models.MonitorHistory, period time.Duration, buckets int) *History {
	historyMap := map[string]models.MonitorStatus{}
	numOfSuccess := 0
	numTotal := 0

	for i := 0; i < buckets; i++ {
		datetime := getHour(time.Now().Add(period * time.Duration(-i)).Truncate(period))
		historyMap[datetime] = models.MonitorUnknown
	}

	for _, _history := range history {
		hour := getHour(_history.CreatedAt.Time.Truncate(period))

		// Skip if not part of the "buckets"
		if _, ok := historyMap[hour]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.MonitorSuccess {
			numOfSuccess++
		}

		// skip if it is already set to failure
		if historyMap[hour] == models.MonitorFailure {
			continue
		}

		historyMap[hour] = _history.Status
	}

	historyHourly := make([]models.MonitorStatus, buckets)
	for i := 0; i < buckets; i++ {
		datetime := getHour(time.Now().Add(period * time.Duration(-buckets+i+1)).Truncate(period))
		historyHourly[i] = historyMap[datetime]
	}

	uptime := 0
	if numTotal > 0 {
		uptime = 100 * numOfSuccess / numTotal
	}

	return &History{
		List:   historyHourly,
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

	overallStatus := models.MonitorSuccess

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

		status := historyResult.List[len(historyResult.List)-1]
		if status != models.MonitorSuccess {
			overallStatus = status
		}

		monitorsWithHistory[i] = &Monitor{
			Name:    monitor.Name,
			Group:   monitor.Group,
			Status:  status,
			History: historyResult,
		}
	}

	monitorsByGroup := map[string][]*Monitor{}
	for _, monitor := range monitorsWithHistory {
		monitorsByGroup[monitor.Group] = append(monitorsByGroup[monitor.Group], monitor)
	}

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
