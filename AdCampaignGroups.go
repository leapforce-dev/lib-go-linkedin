package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCampaignGroupsResponse struct {
	MetaData MetaData          `json:"metadata"`
	Elements []AdCampaignGroup `json:"elements"`
}

type AdCampaignGroup struct {
	Account              string              `json:"account"`
	AllowedCampaignTypes []AdCampaignType    `json:"allowedCampaignTypes"`
	Backfilled           bool                `json:"backfilled"`
	ChangeAuditStamps    AdChangeAuditStamps `json:"changeAuditStamps"`
	Id                   int64               `json:"id"`
	Name                 string              `json:"name"`
	RunSchedule          AdRunSchedule       `json:"runSchedule"`
	ServingStatuses      []string            `json:"servingStatuses"`
	Status               string              `json:"status"`
	Test                 bool                `json:"test"`
	TotalBudget          AdBudget            `json:"totalBudget"`
}

type AdCampaignGroupStatus string

const (
	AdCampaignGroupStatusActive    AdCampaignGroupStatus = "ACTIVE"
	AdCampaignGroupStatusArchived  AdCampaignGroupStatus = "ARCHIVED"
	AdCampaignGroupStatusCanceled  AdCampaignGroupStatus = "CANCELED"
	AdCampaignGroupStatusDraft     AdCampaignGroupStatus = "DRAFT"
	AdCampaignGroupStatusCompleted AdCampaignGroupStatus = "COMPLETED"
)

type SearchAdCampaignGroupsConfig struct {
	Account   int64
	Id        *[]int64
	Status    *[]AdCampaignGroupStatus
	Name      *[]string
	Test      *bool
	PageToken *string
	PageSize  *uint
}

func (service *Service) SearchAdCampaignGroups(config *SearchAdCampaignGroupsConfig) (*[]AdCampaignGroup, *errortools.Error) {
	var values url.Values = url.Values{}
	var pageToken string
	var pageSize uint = countDefault

	values.Set("q", "search")

	if config != nil {
		if config.Id != nil {
			for i, id := range *config.Id {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
			}
		}
		if config.Name != nil {
			for i, name := range *config.Name {
				values.Set(fmt.Sprintf("search.name.values[%v]", i), name)
			}
		}
		if config.Test != nil {
			values.Set("search.test", fmt.Sprintf("%v", *config.Test))
		}
		if config.PageToken != nil {
			pageToken = *config.PageToken
		}
		if config.PageSize != nil {
			pageSize = *config.PageSize
		}
	}

	adCampaignGroups := []AdCampaignGroup{}

	for {
		if pageToken != "" {
			values.Set("pageToken", fmt.Sprintf("%v", pageToken))
		}
		values.Set("pageSize", fmt.Sprintf("%v", pageSize))

		adCampaignGroupsResponse := AdCampaignGroupsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlRest(fmt.Sprintf("adAccounts/%v/adCampaignGroups?%s", config.Account, values.Encode())),
			ResponseModel: &adCampaignGroupsResponse,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(adCampaignGroupsResponse.Elements) == 0 {
			break
		}

		adCampaignGroups = append(adCampaignGroups, adCampaignGroupsResponse.Elements...)

		if config != nil {
			if config.PageToken != nil {
				break
			}
		}

		pageToken = adCampaignGroupsResponse.MetaData.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return &adCampaignGroups, nil
}
