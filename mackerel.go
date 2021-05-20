package main

import "fmt"

type MackerelWebhook struct {
	Orgname string        `json:"orgName"`
	Alert   MackerelAlert `json:"alert"`
}

type MackerelAlert struct {
	Monitorname       string  `json:"monitorName"`
	Criticalthreshold int     `json:"criticalThreshold"`
	Metricvalue       float64 `json:"metricValue"`
	Monitoroperator   string  `json:"monitorOperator"`
	Trigger           string  `json:"trigger"`
	URL               string  `json:"url"`
	Openedat          *int64  `json:"openedAt"`
	Duration          *int64  `json:"duration"`
	Createdat         *int64  `json:"createdAt"`
	Isopen            bool    `json:"isOpen"`
	Metriclabel       string  `json:"metricLabel"`
	ID                string  `json:"id"`
	Closedat          *int64  `json:"closedAt"`
	Status            string  `json:"status"`
}

func (h MackerelWebhook) IncidentTitle() string {
	return fmt.Sprintf("[%s] %s", h.Orgname, h.Alert.Monitorname)
}
