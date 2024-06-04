package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdAccountsResponse struct {
	MetaData MetaData    `json:"metadata"`
	Elements []AdAccount `json:"elements"`
}

type AdAccount struct {
	ChangeAuditStamps              AdChangeAuditStamps `json:"changeAuditStamps"`
	Currency                       string              `json:"currency"`
	Id                             int64               `json:"id"`
	Name                           string              `json:"name"`
	NotifiedOnCampaignOptimization bool                `json:"notifiedOnCampaignOptimization"`
	NotifiedOnCreativeApproval     bool                `json:"notifiedOnCreativeApproval"`
	NotifiedOnCreativeRejection    bool                `json:"notifiedOnCreativeRejection"`
	NotifiedOnEndOfCampaign        bool                `json:"notifiedOnEndOfCampaign"`
	NotifiedOnNewFeaturesEnabled   bool                `json:"notifiedOnNewFeaturesEnabled"`
	Reference                      string              `json:"reference"`
	ServingStatuses                []string            `json:"servingStatuses"`
	Status                         string              `json:"status"`
	Test                           bool                `json:"test"`
	Type                           string              `json:"type"`
	Version                        AdVersion           `json:"version"`
}

type AdAccountStatus string

const (
	AdAccountStatusDraft    AdAccountStatus = "DRAFT"
	AdAccountStatusCanceled AdAccountStatus = "CANCELED"
	AdAccountStatusActive   AdAccountStatus = "ACTIVE"
)

type AdAccountType string

const (
	AdAccountTypeBusiness   AdAccountType = "BUSINESS"
	AdAccountTypeEnterprise AdAccountType = "ENTERPRISE"
)

type SearchAdAccountsConfig struct {
	Status    *[]AdAccountStatus
	Reference *[]string
	Name      *[]string
	Id        *[]int64
	Type      *[]AdAccountType
	Test      *bool
	PageToken *string
	PageSize  *uint
}

func (service *Service) SearchAdAccounts(config *SearchAdAccountsConfig) (*[]AdAccount, *errortools.Error) {
	var values = url.Values{}
	var pageToken string
	var pageSize uint = countDefault

	values.Set("q", "search")

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)

	if config != nil {
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
			}
		}
		if config.Reference != nil {
			for i, reference := range *config.Reference {
				values.Set(fmt.Sprintf("search.reference.values[%v]", i), reference)
			}
		}
		if config.Name != nil {
			for i, name := range *config.Name {
				values.Set(fmt.Sprintf("search.name.values[%v]", i), name)
			}
		}
		if config.Id != nil {
			for i, id := range *config.Id {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Type != nil {
			for i, _type := range *config.Type {
				values.Set(fmt.Sprintf("search.type.values[%v]", i), string(_type))
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

	adAccounts := []AdAccount{}

	for {
		if pageToken != "" {
			values.Set("pageToken", fmt.Sprintf("%v", pageToken))
		}
		values.Set("pageSize", fmt.Sprintf("%v", pageSize))

		adAccountsResponse := AdAccountsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("adAccounts?%s", values.Encode())),
			ResponseModel:     &adAccountsResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(adAccountsResponse.Elements) == 0 {
			break
		}

		adAccounts = append(adAccounts, adAccountsResponse.Elements...)

		if config != nil {
			if config.PageToken != nil {
				break
			}
		}

		pageToken = adAccountsResponse.MetaData.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return &adAccounts, nil
}

func (service *Service) GetAdAccount(accountId int64) (*AdAccount, *errortools.Error) {
	id := []int64{accountId}

	accounts, e := service.SearchAdAccounts(&SearchAdAccountsConfig{
		Id: &id,
	})
	if e != nil {
		return nil, e
	}

	if accounts == nil {
		return nil, errortools.ErrorMessage("Account not found.")
	}

	if len(*accounts) == 0 {
		return nil, errortools.ErrorMessage("Account not found.")
	}

	if len(*accounts) > 1 {
		return nil, errortools.ErrorMessage("Multiple accounts found.")
	}

	return &(*accounts)[0], nil
}
