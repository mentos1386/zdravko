package handlers

import (
	"context"
	"net/http"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type IndexData struct {
	*components.Base
	HealthChecks       []*HealthCheck
	HealthchecksLength int
	TimeRange          string
	Status             string
}

type HealthCheck struct {
	Name          string
	Status        string
	HistoryDaily  *History
	HistoryHourly *History
}

type History struct {
	History []string
	Uptime  int
}

func getDay(date time.Time) string {
	return date.Format("2006-01-02")
}

func getHour(date time.Time) string {
	return date.Format("2006-01-02T15:04")
}

func getDailyHistory(history []models.HealthcheckHistory) *History {
	numDays := 90
	historyDailyMap := map[string]string{}
	numOfSuccess := 0
	numTotal := 0

	for i := 0; i < numDays; i++ {
		day := getDay(time.Now().AddDate(0, 0, -i).Truncate(time.Hour * 24))
		historyDailyMap[day] = models.HealthcheckUnknown
	}

	for _, _history := range history {
		day := getDay(_history.CreatedAt.Truncate(time.Hour * 24))

		// skip if day is not in the last 90 days
		if _, ok := historyDailyMap[day]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.HealthcheckSuccess {
			numOfSuccess++
		}

		// skip if day is already set to failure
		if historyDailyMap[day] == models.HealthcheckFailure {
			continue
		}

		historyDailyMap[day] = _history.Status
	}

	historyDaily := make([]string, numDays)
	for i := 0; i < numDays; i++ {
		day := getDay(time.Now().AddDate(0, 0, -numDays+i+1).Truncate(time.Hour * 24))
		historyDaily[i] = historyDailyMap[day]
	}

	uptime := 0
	if numTotal > 0 {
		uptime = 100 * numOfSuccess / numTotal
	}

	return &History{
		History: historyDaily,
		Uptime:  uptime,
	}
}

func getHourlyHistory(history []models.HealthcheckHistory) *History {
	numHours := 48
	historyHourlyMap := map[string]string{}
	numOfSuccess := 0
	numTotal := 0

	for i := 0; i < numHours; i++ {
		hour := getHour(time.Now().Add(time.Hour * time.Duration(-i)).Truncate(time.Hour))
		historyHourlyMap[hour] = models.HealthcheckUnknown
	}

	for _, _history := range history {
		hour := getHour(_history.CreatedAt.Truncate(time.Hour))

		// skip if day is not in the last 90 days
		if _, ok := historyHourlyMap[hour]; !ok {
			continue
		}

		numTotal++
		if _history.Status == models.HealthcheckSuccess {
			numOfSuccess++
		}

		// skip if day is already set to failure
		if historyHourlyMap[hour] == models.HealthcheckFailure {
			continue
		}

		historyHourlyMap[hour] = _history.Status
	}

	historyHourly := make([]string, numHours)
	for i := 0; i < numHours; i++ {
		hour := getHour(time.Now().Add(time.Hour * time.Duration(-numHours+i+1)).Truncate(time.Hour))
		historyHourly[i] = historyHourlyMap[hour]
	}

	uptime := 0
	if numTotal > 0 {
		uptime = 100 * numOfSuccess / numTotal
	}

	return &History{
		History: historyHourly,
		Uptime:  uptime,
	}
}

func (h *BaseHandler) Index(c echo.Context) error {
	ctx := context.Background()
	healthchecks, err := services.GetHealthchecks(ctx, h.query)
	if err != nil {
		return err
	}

	timeRange := c.QueryParam("time-range")
	if timeRange != "48hours" && timeRange != "90days" {
		timeRange = "90days"
	}

	overallStatus := "SUCCESS"

	healthchecksWithHistory := make([]*HealthCheck, len(healthchecks))
	for i, healthcheck := range healthchecks {
		historyDaily := getDailyHistory(healthcheck.History)
		historyHourly := getHourlyHistory(healthcheck.History)

		status := historyDaily.History[89]
		if status != models.HealthcheckSuccess {
			overallStatus = status
		}

		healthchecksWithHistory[i] = &HealthCheck{
			Name:          healthcheck.Name,
			Status:        status,
			HistoryDaily:  historyDaily,
			HistoryHourly: historyHourly,
		}
	}

	return c.Render(http.StatusOK, "index.tmpl", &IndexData{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Status"),
			Navbar:       Pages,
		},
		HealthChecks:       healthchecksWithHistory,
		HealthchecksLength: len(healthchecks),
		TimeRange:          timeRange,
		Status:             overallStatus,
	})
}
