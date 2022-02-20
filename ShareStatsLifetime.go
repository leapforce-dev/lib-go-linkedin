package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type ShareStatsLifetimeResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []ShareStatsLifetime `json:"elements"`
}

type ShareStatsLifetime struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	Share                *string              `json:"share"`
}

func (service *Service) GetShareStatsLifetime(organizationId int64, shareIds *[]string) (*[]ShareStatsLifetime, *http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationId))

	if shareIds != nil {
		for index, shareId := range *shareIds {
			values.Set(fmt.Sprintf("shares[%v]", index), fmt.Sprintf("urn:li:share:%s", shareId))
		}
	}

	shareStatsResponse := ShareStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &shareStatsResponse,
	}
	_, response, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, response, e
	}

	return &shareStatsResponse.Elements, response, nil
}
