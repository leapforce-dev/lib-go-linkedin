package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCreativesResponse struct {
	MetaData MetaData     `json:"metadata"`
	Elements []AdCreative `json:"elements"`
}

type AdCreative struct {
	Campaign          string              `json:"campaign"`
	ChangeAuditStamps AdChangeAuditStamps `json:"changeAuditStamps"`
	Id                int64               `json:"id"`
	ProcessingState   string              `json:"processingState"`
	Reference         string              `json:"reference"`
	Review            struct {
		ReviewStatus string `json:"reviewStatus"`
	} `json:"review"`
	ServingStatuses []string `json:"servingStatuses"`
	Status          string   `json:"status"`
	Test            bool     `json:"test"`
	Type            string   `json:"type"`
	Variables       struct {
		ClickUri string                  `json:"clickUri"`
		Data     AdCreativeVariablesData `json:"data"`
	} `json:"variables"`
	Version AdVersion `json:"version"`
}

type AdCreativeVariablesData struct {
	TextAd                  *TextAdCreativeVariables                  `json:"com.linkedin.ads.TextAdCreativeVariables"`
	SponsoredUpdate         *SponsoredUpdateCreativeVariables         `json:"com.linkedin.ads.SponsoredUpdateCreativeVariables"`
	SponsoredInMail         *SponsoredInMailCreativeVariables         `json:"com.linkedin.ads.SponsoredInMailCreativeVariables"`
	SponsoredVideo          *SponsoredVideoCreativeVariables          `json:"com.linkedin.ads.SponsoredVideoCreativeVariables"`
	SponsoredUpdateCarousel *SponsoredUpdateCarouselCreativeVariables `json:"com.linkedin.ads.SponsoredUpdateCarouselCreativeVariables"`
	FollowCompany           *FollowCompanyCreativeVariablesV2         `json:"com.linkedin.ads.FollowCompanyCreativeVariablesV2"`
	Spotlight               *SpotlightCreativeVariablesV2             `json:"com.linkedin.ads.SpotlightCreativeVariablesV2"`
	Jobs                    *JobsCreativeVariablesV2                  `json:"com.linkedin.ads.JobsCreativeVariablesV2"`
}

type TextAdCreativeVariables struct {
	Text  string `json:"text"`
	Title string `json:"title"`
}

type SponsoredUpdateCreativeVariables struct {
	Activity               string `json:"activity"`
	DirectSponsoredContent bool   `json:"directSponsoredContent"`
	Share                  string `json:"share"`
}

type SponsoredInMailCreativeVariables struct {
	Content string `json:"content"`
}

type SponsoredVideoCreativeVariables struct {
	DurationMicro            int64       `json:"durationMicro"`
	MediaAsset               string      `json:"mediaAsset"`
	UserGeneratedContentPost string      `json:"userGeneratedContentPost"`
	VideoAspectRatio         AspectRatio `json:"videoAspectRatio"`
}

type SponsoredUpdateCarouselCreativeVariables struct {
	Activity               string `json:"activity"`
	DirectSponsoredContent bool   `json:"directSponsoredContent"`
	Share                  string `json:"share"`
}

type FollowCompanyCreativeVariablesV2 struct {
	CallToAction string `json:"callToAction"`
	Description  struct {
		FollowCompanyDescription string `json:"com.linkedin.ads.FollowCompanyDescription"`
	} `json:"description"`
	Headline struct {
		FollowCompanyHeadline string `json:"com.linkedin.ads.FollowCompanyHeadline"`
	} `json:"headline"`
	Organization struct {
		Company string `json:"company"`
	} `json:"organization"`
	OrganizationLogo string `json:"organizationLogo"`
	OrganizationName string `json:"organizationName"`
}

type SpotlightCreativeVariablesV2 struct {
	CallToAction           string `json:"callToAction"`
	CustomBackground       string `json:"customBackground"`
	Description            string `json:"description"`
	ForumName              string `json:"forumName"`
	Headline               string `json:"headline"`
	Logo                   string `json:"logo"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
}

type JobsCreativeVariablesV2 struct {
	ButtonLabel struct {
		JobsButtonLabel string `json:"com.linkedin.ads.JobsButtonLabel"`
	} `json:"buttonLabel"`
	CompanyName string `json:"companyName"`
	Headline    struct {
		JobsHeadline string `json:"com.linkedin.ads.JobsHeadline"`
	} `json:"headline"`
	Logo         string `json:"logo"`
	Organization string `json:"organization"`
}

type AdCreativeStatus string

const (
	AdCreativeStatusActive   AdCreativeStatus = "ACTIVE"
	AdCreativeStatusPaused   AdCreativeStatus = "PAUSED"
	AdCreativeStatusDraft    AdCreativeStatus = "DRAFT"
	AdCreativeStatusArchived AdCreativeStatus = "ARCHIVED"
	AdCreativeStatusCanceled AdCreativeStatus = "CANCELED"
)

type AdCreativeType string

const (
	AdCreativeTypeTextAd           AdCreativeType = "TEXT_AD"
	AdCreativeTypeSponsoredUpdates AdCreativeType = "SPONSORED_UPDATES"
	AdCreativeTypeSponsoredInmails AdCreativeType = "SPONSORED_INMAILS"
	AdCreativeTypeDynamic          AdCreativeType = "DYNAMIC"
)

type SearchAdCreativesConfig struct {
	Campaign  *[]int64
	Id        *[]int64
	Reference *[]string
	Status    *[]AdCreativeStatus
	Test      *bool
	PageToken *string
	PageSize  *uint
}

func (service *Service) SearchAdCreatives(config *SearchAdCreativesConfig) (*[]AdCreative, *errortools.Error) {
	var values url.Values = url.Values{}
	var pageToken string
	var pageSize uint = countDefault

	values.Set("q", "search")

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)

	if config != nil {
		if config.Campaign != nil {
			for i, campaign := range *config.Campaign {
				values.Set(fmt.Sprintf("search.campaign.values[%v]", i), fmt.Sprintf("urn:li:sponsoredCampaign:%v", campaign))
			}
		}
		if config.Id != nil {
			for i, id := range *config.Id {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Reference != nil {
			for i, reference := range *config.Reference {
				values.Set(fmt.Sprintf("search.reference.values[%v]", i), reference)
			}
		}
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
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

	adCreatives := []AdCreative{}

	for {
		if pageToken != "" {
			values.Set("pageToken", fmt.Sprintf("%v", pageToken))
		}
		values.Set("pageSize", fmt.Sprintf("%v", pageSize))

		adCreativesResponse := AdCreativesResponse{}

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("adCreatives?%s", values.Encode())),
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

func (service *Service) GetAdCreatives(campaignId int64) (*[]AdCreative, *errortools.Error) {
	campaign := []int64{campaignId}

	creatives, e := service.SearchAdCreatives(&SearchAdCreativesConfig{
		Campaign: &campaign,
	})
	if e != nil {
		return nil, e
	}

	if creatives == nil {
		return nil, nil
	}

	return creatives, nil
}
