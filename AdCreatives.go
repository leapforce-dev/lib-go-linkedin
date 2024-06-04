package linkedin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCreativesResponse struct {
	MetaData MetaData     `json:"metadata"`
	Elements []AdCreative `json:"elements"`
}

type AdCreative struct {
	Account            *string            `json:"account,omitempty"`
	Campaign           *string            `json:"campaign,omitempty"`
	Content            *AdCreativeContent `json:"content,omitempty"`
	CreatedAt          *int64             `json:"createdAt,omitempty"`
	CreatedBy          *string            `json:"createdBy,omitempty"`
	Id                 *string            `json:"id,omitempty"`
	InlineContent      *InlineContent     `json:"inlineContent,omitempty"`
	IntendedStatus     *string            `json:"intendedStatus,omitempty"`
	IsServing          *bool              `json:"isServing,omitempty"`
	IsTest             *bool              `json:"isTest,omitempty"`
	LastModifiedAt     *int64             `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string            `json:"lastModifiedBy,omitempty"`
	Review             *AdCreativeReview  `json:"review,omitempty"`
	ServingHoldReasons *[]string          `json:"servingHoldReasons,omitempty"`
}

type AdCreativeContent struct {
	Reference string                      `json:"reference"`
	TextAd    *AdCreativeContentTextAd    `json:"textAd"`
	Jobs      *AdCreativeContentJobs      `json:"jobs"`
	Spotlight *AdCreativeContentSpotlight `json:"spotlight"`
	Follow    *AdCreativeContentFollow    `json:"follow"`
}

type AdCreativeContentTextAd struct {
	Image       string `json:"image"`
	Description string `json:"description"`
	Headline    string `json:"headline"`
	LandingPage string `json:"landingPage"`
}

type AdCreativeContentJobs struct {
	Logo                   string `json:"logo"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
	OrganizationName       string `json:"organizationName"`
	Headline               struct {
		PreApproved string `json:"preApproved"`
	} `json:"headline"`
	ButtonLabel struct {
		PreApproved string `json:"preApproved"`
	} `json:"buttonLabel"`
}

type AdCreativeContentSpotlight struct {
	CallToAction           string `json:"callToAction"`
	Description            string `json:"description"`
	Headline               string `json:"headline"`
	LandingPage            string `json:"landingPage"`
	Logo                   string `json:"logo"`
	OrganizationName       string `json:"organizationName"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
}

type AdCreativeContentFollow struct {
	OrganizationName string `json:"organizationName"`
	Logo             string `json:"logo"`
	Headline         struct {
		PreApproved string `json:"preApproved"`
	} `json:"headline"`
	Description struct {
		PreApproved string `json:"preApproved"`
	} `json:"description"`
	CallToAction           string `json:"callToAction"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
}

type InlineContent struct {
	Post struct {
		AdContext struct {
			DscAdAccount string `json:"dscAdAccount"`
			DscStatus    string `json:"dscStatus"`
		} `json:"adContext"`
		Author                    string `json:"author"`
		Commentary                string `json:"commentary"`
		Visibility                string `json:"visibility"`
		LifecycleState            string `json:"lifecycleState"`
		IsReshareDisabledByAuthor bool   `json:"isReshareDisabledByAuthor"`
		Content                   struct {
			Media struct {
				Title string `json:"title"`
				Id    string `json:"id"`
			} `json:"media"`
		} `json:"content"`
	} `json:"post"`
}

type AdCreativeReview struct {
	Status           string   `json:"status"`
	RejectionReasons []string `json:"rejectionReasons"`
}

type SearchAdCreativesConfig struct {
	Account                                 int64
	Campaigns                               *[]string
	ContentReferences                       *[]string
	Creatives                               *[]string
	IntendedStatuses                        *[]string
	IsTestAccount                           *bool
	IsTotalIncluded                         *bool
	LeadgenCreativeCallToActionDestinations *[]string
	SortOrder                               *string
	PageToken                               *string
	PageSize                                *uint
}

func (service *Service) SearchAdCreatives(config *SearchAdCreativesConfig) (*[]AdCreative, *errortools.Error) {
	var params []string
	var pageToken string
	var pageSize = countDefault

	params = append(params, "q=criteria")

	if config != nil {
		if config.Campaigns != nil {
			if len(*config.Campaigns) > 0 {
				params = append(params, fmt.Sprintf("campaigns=List(%s)", url.QueryEscape(strings.Join(*config.Campaigns, ","))))
			}
		}
		if config.ContentReferences != nil {
			if len(*config.ContentReferences) > 0 {
				params = append(params, fmt.Sprintf("contentReferences=List(%s)", strings.Join(*config.ContentReferences, ",")))
			}
		}
		if config.Creatives != nil {
			if len(*config.Creatives) > 0 {
				params = append(params, fmt.Sprintf("adCreatives=List(%s)", strings.Join(*config.Creatives, ",")))
			}
		}
		if config.IntendedStatuses != nil {
			if len(*config.IntendedStatuses) > 0 {
				params = append(params, fmt.Sprintf("intendedStatuses=(value:List(%s))", strings.Join(*config.IntendedStatuses, ",")))
			}
		}
		if config.IsTestAccount != nil {
			params = append(params, fmt.Sprintf("isTestAccount=%v", *config.IsTestAccount))
		}
		if config.IsTotalIncluded != nil {
			params = append(params, fmt.Sprintf("isTotalIncluded=%v", *config.IsTotalIncluded))
		}
		if config.LeadgenCreativeCallToActionDestinations != nil {
			if len(*config.LeadgenCreativeCallToActionDestinations) > 0 {
				params = append(params, fmt.Sprintf("leadgenCreativeCallToActionDestinations=List(%s)", strings.Join(*config.LeadgenCreativeCallToActionDestinations, ",")))
			}
		}
		if config.SortOrder != nil {
			params = append(params, fmt.Sprintf("sortOrder=%s", *config.SortOrder))
		}
		if config.PageToken != nil {
			pageToken = *config.PageToken
		}
		if config.PageSize != nil {
			pageSize = *config.PageSize
		}
	}

	var adCreatives []AdCreative

	for {
		params_ := params
		if pageToken != "" {
			params_ = append(params_, fmt.Sprintf("pageToken=%s", pageToken))
		}
		params_ = append(params_, fmt.Sprintf("pageSize=%v", pageSize))

		adCreativesResponse := AdCreativesResponse{}

		var header = http.Header{}
		header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
		header.Set("X-RestLi-Method", "FINDER")

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("adAccounts/%v/creatives?%s", config.Account, strings.Join(params_, "&"))),
			ResponseModel:     &adCreativesResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(adCreativesResponse.Elements) == 0 {
			break
		}

		adCreatives = append(adCreatives, adCreativesResponse.Elements...)

		if config != nil {
			if config.PageToken != nil {
				break
			}
		}

		pageToken = adCreativesResponse.MetaData.NextPageToken

		if pageToken == "" {
			break
		}
	}

	return &adCreatives, nil
}
