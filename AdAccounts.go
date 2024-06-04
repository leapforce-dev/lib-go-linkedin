package linkedin

import (
	"fmt"
	"net/http"
	"strings"

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
	Id        *[]string
	Type      *[]AdAccountType
	Test      *bool
	PageToken *string
	PageSize  *uint
}

func (service *Service) SearchAdAccounts(config *SearchAdAccountsConfig) (*[]AdAccount, *errortools.Error) {
	var params []string
	var pageToken string
	var pageSize = countDefault

	params = append(params, "q=search")

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)

	if config != nil {
		var search []string

		if config.Status != nil {
			var searchStatus []string

			for _, status := range *config.Status {
				searchStatus = append(searchStatus, string(status))
			}

			search = append(search, fmt.Sprintf("status:(values:List(%s))", strings.Join(searchStatus, ",")))
		}
		if config.Reference != nil {
			var searchReference []string

			for _, reference := range *config.Reference {
				searchReference = append(searchReference, reference)
			}

			search = append(search, fmt.Sprintf("reference:(values:List(%s))", strings.Join(searchReference, ",")))
		}
		if config.Name != nil {
			var searchName []string

			for _, name := range *config.Name {
				searchName = append(searchName, name)
			}

			search = append(search, fmt.Sprintf("name:(values:List(%s))", strings.Join(searchName, ",")))
		}
		if config.Id != nil {
			var searchId []string

			for _, id := range *config.Id {
				searchId = append(searchId, fmt.Sprintf("%s", id))
			}

			search = append(search, fmt.Sprintf("id:(values:List(%s))", strings.Join(searchId, ",")))
		}
		if config.Type != nil {
			var searchType []string

			for _, _type := range *config.Type {
				searchType = append(searchType, fmt.Sprintf("%v", string(_type)))
			}

			search = append(search, fmt.Sprintf("id:(values:List(%s))", strings.Join(searchType, ",")))
		}
		if config.Test != nil {
			search = append(search, fmt.Sprintf("test:(values:List(%s))", fmt.Sprintf("%v", *config.Test)))

		}

		params = append(params, fmt.Sprintf("search=(%s)", strings.Join(search, ",")))

		if config.PageToken != nil {
			pageToken = *config.PageToken
		}
		if config.PageSize != nil {
			pageSize = *config.PageSize
		}
	}

	var adAccounts []AdAccount

	for {
		params_ := params
		if pageToken != "" {
			params_ = append(params_, fmt.Sprintf("pageToken=%s", pageToken))
		}
		params_ = append(params_, fmt.Sprintf("pageSize=%v", pageSize))

		adAccountsResponse := AdAccountsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("adAccounts?%s", strings.Join(params_, "&"))),
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
	var adAccount AdAccount
	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest(fmt.Sprintf("adAccounts/%v", accountId)),
		ResponseModel: &adAccount,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &adAccount, nil
}
