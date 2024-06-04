package linkedin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCampaignsResponse struct {
	MetaData MetaData     `json:"metadata"`
	Elements []AdCampaign `json:"elements"`
}

type AdCampaign struct {
	Account                  string              `json:"account"`
	AssociatedEntity         string              `json:"associatedEntity"`
	AudienceExpansionEnabled bool                `json:"audienceExpansionEnabled"`
	CampaignGroup            string              `json:"campaignGroup"`
	ChangeAuditStamps        AdChangeAuditStamps `json:"changeAuditStamps"`
	CostType                 string              `json:"costType"`
	CreativeSelection        string              `json:"creativeSelection"`
	DailyBudget              AdBudget            `json:"dailyBudget"`
	Format                   string              `json:"format"`
	Id                       int64               `json:"id"`
	Locale                   AdLocale            `json:"locale"`
	Name                     string              `json:"name"`
	ObjectiveType            string              `json:"objectiveType"`
	OffsiteDeliveryEnabled   bool                `json:"offsiteDeliveryEnabled"`
	OffsitePreferences       struct {
		IABCategories struct {
			Exclude []string `json:"exclude"`
			Include []string `json:"include"`
		} `json:"iabCategories"`
		PublisherRestrictionFiles struct {
			Exclude []string `json:"exclude"`
		} `json:"publisherRestrictionFiles"`
	} `json:"offsitePreferences"`
	OptimizationTargetType string        `json:"optimizationTargetType"`
	PacingStrategy         string        `json:"pacingStrategy"`
	RunSchedule            AdRunSchedule `json:"runSchedule"`
	ServingStatuses        []string      `json:"servingStatuses"`
	Status                 string        `json:"status"`
	Targeting              struct {
		IncludedTargetingFacets struct {
			Employers        []string   `json:"employers"`
			Locations        []string   `json:"locations"`
			InterfaceLocales []AdLocale `json:"interfaceLocales"`
		} `json:"includedTargetingFacets"`
	} `json:"targeting"`
	TargetingCriteria json.RawMessage `json:"targetingCriteria"`
	Test              bool            `json:"test"`
	TotalBudget       AdBudget        `json:"totalBudget"`
	Type              string          `json:"type"`
	UnitCost          AdBudget        `json:"unitCost"`
	Version           AdVersion       `json:"version"`
}

type AdCampaignStatus string

const (
	AdCampaignStatusActive    AdCampaignStatus = "ACTIVE"
	AdCampaignStatusPaused    AdCampaignStatus = "PAUSED"
	AdCampaignStatusArchived  AdCampaignStatus = "ARCHIVED"
	AdCampaignStatusCompleted AdCampaignStatus = "COMPLETED"
	AdCampaignStatusCanceled  AdCampaignStatus = "CANCELED"
	AdCampaignStatusDraft     AdCampaignStatus = "DRAFT"
)

type AdCampaignType string

const (
	AdCampaignTypeTextAd           AdCampaignType = "TEXT_AD"
	AdCampaignTypeSponsoredUpdates AdCampaignType = "SPONSORED_UPDATES"
	AdCampaignTypeSponsoredInmails AdCampaignType = "SPONSORED_INMAILS"
	AdCampaignTypeDynamic          AdCampaignType = "DYNAMIC"
)

type SearchAdCampaignsConfig struct {
	Account          int64
	CampaignGroup    *[]int64
	AssociatedEntity *[]string
	Id               *[]int64
	Status           *[]AdCampaignStatus
	Type             *[]AdCampaignType
	Name             *[]string
	Test             *bool
	PageToken        *string
	PageSize         *uint
}

func (service *Service) SearchAdCampaigns(config *SearchAdCampaignsConfig) (*[]AdCampaign, *errortools.Error) {
	var values url.Values = url.Values{}
	var pageToken string
	var pageSize uint = countDefault

	values.Set("q", "search")

	if config != nil {
		if config.CampaignGroup != nil {
			for i, campaignGroup := range *config.CampaignGroup {
				values.Set(fmt.Sprintf("search.campaignGroup.values[%v]", i), fmt.Sprintf("urn:li:sponsoredCampaignGroup:%v", campaignGroup))
			}
		}
		if config.AssociatedEntity != nil {
			for i, associatedEntity := range *config.AssociatedEntity {
				values.Set(fmt.Sprintf("search.associatedEntity.values[%v]", i), associatedEntity)
			}
		}
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
		if config.Type != nil {
			for i, _type := range *config.Type {
				values.Set(fmt.Sprintf("search.type.values[%v]", i), string(_type))
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

	adCampaigns := []AdCampaign{}

	for {
		if pageToken != "" {
			values.Set("pageToken", fmt.Sprintf("%v", pageToken))
		}
		values.Set("pageSize", fmt.Sprintf("%v", pageSize))

		adCampaignsResponse := AdCampaignsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlRest(fmt.Sprintf("adAccounts/%v/adCampaigns?%s", config.Account, values.Encode())),
			ResponseModel: &adCampaignsResponse,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(adCampaignsResponse.Elements) == 0 {
			break
		}

		adCampaigns = append(adCampaigns, adCampaignsResponse.Elements...)

		if config != nil {
			if config.PageToken != nil {
				break
			}
		}

		pageToken = adCampaignsResponse.MetaData.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return &adCampaigns, nil
}
