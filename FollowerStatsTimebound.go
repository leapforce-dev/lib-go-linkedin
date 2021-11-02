package linkedin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type FollowerStatsTimeboundResponse struct {
	Paging   Paging                   `json:"paging"`
	Elements []FollowerStatsTimebound `json:"elements"`
}

type FollowerStatsTimebound struct {
	TimeRange            TimeRange     `json:"timeRange"`
	FollowerGains        FollowerGains `json:"followerGains"`
	OrganizationalEntity string        `json:"organizationalEntity"`
}

type FollowerGains struct {
	OrganicFollowerGain int64 `json:"organicFollowerGain"`
	PaidFollowerGain    int64 `json:"paidFollowerGain"`
}

func (service *Service) GetFollowerStatsTimebound(organizationID int64, startDateUnix int64, endDateUnix int64) (*[]FollowerStatsTimebound, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	followerStatsResponse := FollowerStatsTimeboundResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("organizationalEntityFollowerStatistics?%s", values.Encode())),
		ResponseModel: &followerStatsResponse,
	}
	_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &followerStatsResponse.Elements, nil
}
