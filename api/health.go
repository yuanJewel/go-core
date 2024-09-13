package api

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/common/version"
	"net/http"
)

const (
	HealthcheckStatusUP   = "UP"
	HealthcheckStatusDOWN = "DOWN"
)

type HealthcheckSpec struct {
	Status  string                      `json:"status"`
	Details map[string]*HealthcheckItem `json:"details"`
}

type HealthcheckItem struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}

func (hc *HealthcheckSpec) MergeHealthcheckStatus() {
	if hc.Status == "" {
		hc.Status = HealthcheckStatusUP
	}
	hc.addVersion()
	if hc.Status != HealthcheckStatusUP {
		return
	}
	for _, v := range hc.Details {
		if v.Status != HealthcheckStatusUP {
			hc.Status = HealthcheckStatusDOWN
			return
		}
	}
}

func (hc *HealthcheckSpec) AddItem(name string, item *HealthcheckItem) {
	if hc.Details == nil {
		hc.Details = map[string]*HealthcheckItem{}
	}
	hc.Details[name] = item
}

func (hc *HealthcheckSpec) addVersion() {
	item := &HealthcheckItem{Status: HealthcheckStatusUP, Details: map[string]interface{}{"version": fmt.Sprintf("%s-%s", version.Version, version.BuildDate)}}
	hc.AddItem("Version", item)
}

func healthCheckHandle(errFunc func() map[string]error) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		hc := &HealthcheckSpec{}
		for name, err := range errFunc() {
			item := &HealthcheckItem{Status: HealthcheckStatusUP}
			if err != nil {
				item.Status = HealthcheckStatusDOWN
				item.Details = map[string]interface{}{"exception": err.Error()}
			}
			hc.AddItem(name, item)
		}
		hc.MergeHealthcheckStatus()
		if hc.Status != HealthcheckStatusUP {
			ctx.StatusCode(http.StatusInternalServerError)
		}
		ResponseBody(ctx, ResponseInit(ctx), hc)
	}
}
