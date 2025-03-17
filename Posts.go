package linkedin

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Post struct {
	Author                    string              `json:"author,omitempty"`
	AdContext                 *PostAdContext      `json:"adContext,omitempty"`
	Commentary                string              `json:"commentary,omitempty"`
	Container                 string              `json:"container,omitempty"`
	Content                   *PostContent        `json:"content,omitempty"`
	ContentCallToActionLabel  string              `json:"contentCallToActionLabel,omitempty"`
	CreatedAt                 int64               `json:"createdAt,omitempty"`
	Distribution              PostDistribution    `json:"distribution,omitempty"`
	Id                        string              `json:"id,omitempty"`
	IsReshareDisabledByAuthor bool                `json:"isReshareDisabledByAuthor,omitempty"`
	LastModifiedAt            int64               `json:"lastModifiedAt,omitempty"`
	LifecycleState            string              `json:"lifecycleState,omitempty"`
	LifecycleStateInfo        *LifecycleStateInfo `json:"lifecycleStateInfo,omitempty"`
	PublishedAt               int64               `json:"publishedAt,omitempty"`
	ReshareContext            *PostReshareContext `json:"reshareContext,omitempty"`
	Visibility                string              `json:"visibility,omitempty"`
}

type PostAdContext struct {
	DscStatus    string `json:"dscStatus"`
	DscAdType    string `json:"dscAdType"`
	IsDsc        bool   `json:"isDsc"`
	DscAdAccount string `json:"dscAdAccount"`
}

type PostDistribution struct {
	FeedDistribution               string               `json:"feedDistribution"`
	TargetEntities                 []DistributionTarget `json:"targetEntities,omitempty"`
	ThirdPartyDistributionChannels []string             `json:"thirdPartyDistributionChannels,omitempty"`
}

type PostContent struct {
	Media      *PostContentMedia      `json:"media,omitempty"`
	MultiImage *PostContentMultiImage `json:"multiImage,omitempty"`
}

type PostContentMedia struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type PostContentMultiImage struct {
	Images []PostContentMultiImageImage `json:"images"`
}

type PostContentMultiImageImage struct {
	Id      string `json:"id"`
	AltText string `json:"altText"`
}

type LifecycleStateInfo struct {
	IsEditedByAuthor bool `json:"isEditedByAuthor"`
}

type PostReshareContext struct {
	Parent string `json:"parent"`
}

type DistributionTarget struct {
	Degrees          *[]string `json:"degrees"`
	FieldsOfStudy    *[]string `json:"fieldsOfStudy"`
	Industries       *[]string `json:"industries"`
	InterfaceLocales *[]Locale `json:"interfaceLocales"`
	JobFunctions     *[]string `json:"jobFunctions"`
	GeoLocations     *[]string `json:"geoLocations"`
	Seniorities      *[]string `json:"seniorities"`
	StaffCountRanges *[]string `json:"staffCountRanges"`
}

type Locale struct {
	Locale struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"locale"`
}

func (service *Service) CreatePost(post *Post) (string, *errortools.Error) {
	if service == nil {
		return "", errortools.ErrorMessage("Service pointer is nil")
	}
	if post == nil {
		return "", errortools.ErrorMessage("Post pointer is nil")
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlRest("posts"),
		BodyModel: post,
	}
	_, resp, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return "", e
	}

	var postId = resp.Header.Get("X-Linkedin-Id")
	if postId == "" {
		return "", errortools.ErrorMessage("CreatePost did not return PostID in header")
	}

	return postId, nil
}

type PostsByOwnerConfig struct {
	OrganizationId         int64
	Fields                 *string
	CreatedStartDateUnix   *int64
	CreatedEndDateUnix     *int64
	PublishedStartDateUnix *int64
	PublishedEndDateUnix   *int64
}

type PostsByOwnerResponse struct {
	Paging   Paging `json:"paging"`
	Elements []Post `json:"elements"`
}

func (service *Service) PostsByOwner(cfg *PostsByOwnerConfig) (*[]Post, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}
	if cfg == nil {
		return nil, errortools.ErrorMessage("GetPostsByOwnerConfig pointer is nil")
	}

	start := 0
	count := 50

	var posts []Post

	for {
		values := url.Values{}
		values.Set("q", "author")
		values.Set("author", fmt.Sprintf("urn:li:organization:%v", cfg.OrganizationId))
		if cfg.Fields != nil {
			values.Set("fields", *cfg.Fields)
		}
		values.Set("start", strconv.Itoa(start))
		values.Set("count", strconv.Itoa(count))

		postsResponse := PostsByOwnerResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlRest(fmt.Sprintf("posts?%s", values.Encode())),
			ResponseModel: &postsResponse,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		for _, post := range postsResponse.Elements {

			if cfg.CreatedEndDateUnix != nil {
				if post.CreatedAt > *cfg.CreatedEndDateUnix {
					continue
				}
			}

			if cfg.CreatedStartDateUnix != nil {
				if post.CreatedAt < *cfg.CreatedStartDateUnix {
					continue
				}
			}

			if cfg.PublishedEndDateUnix != nil {
				if post.PublishedAt > *cfg.PublishedEndDateUnix {
					continue
				}
			}

			if cfg.PublishedStartDateUnix != nil {
				if post.PublishedAt < *cfg.PublishedStartDateUnix {
					continue
				}
			}

			posts = append(posts, post)
		}

		if !postsResponse.Paging.HasLink("next") {
			break
		}

		start += count
	}

	return &posts, nil
}

type PostsResponse struct {
	Results map[string]Post `json:"results"`
}

func (service *Service) Posts(urns []string) (*[]Post, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var posts []Post

	// deduplicate urns
	var _urnsMap = make(map[string]bool)
	var _urns []string
	for _, urn := range urns {
		_, ok := _urnsMap[urn]
		if ok {
			continue
		}
		_urnsMap[urn] = true
		_urns = append(_urns, url.QueryEscape(urn))
	}

	for {
		var _urnsBatch []string

		if len(_urns) > int(maxUrnsPerCall) {
			_urnsBatch = _urns[:maxUrnsPerCall]
			_urns = _urns[maxUrnsPerCall:]
		} else {
			_urnsBatch = _urns
			_urns = []string{}
		}

		postsResponse := PostsResponse{}

		var header = http.Header{}
		header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
		header.Set("X-RestLi-Method", "BATCH_GET")

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("posts?ids=List(%s)", strings.Join(_urnsBatch, ","))),
			ResponseModel:     &postsResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		for _, post := range postsResponse.Results {
			posts = append(posts, post)
		}

		if uint(len(_urns)) <= maxUrnsPerCall {
			break
		} else {
			_urns = _urns[maxUrnsPerCall:]
		}
	}

	return &posts, nil
}
