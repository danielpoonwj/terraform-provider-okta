package sdk

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type ListGroupPushResponse struct {
	Mappings []GroupPushMapping `json:"mappings"`
}

type GroupPushMapping struct {
	MappingID         string `json:"mappingId"`
	Status            string `json:"status"`
	SourceUserGroupId string `json:"sourceUserGroupId"`
}

func requestAdminURL(req *http.Request) {
	hostSlice := strings.Split(req.URL.Host, ".")
	hostSlice[0] = fmt.Sprintf("%s-admin", hostSlice[0])
	hostStr := strings.Join(hostSlice, ".")

	req.URL.Host = hostStr
	req.Host = hostStr
}

func (m *ApiSupplement) GetGroupPushMapping(ctx context.Context, appID, groupID string) (*GroupPushMapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/internal/instance/%s/grouppush", appID)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	// API only accessible on admin
	requestAdminURL(req)

	listGroupPushResponse := &ListGroupPushResponse{}
	resp, err := m.RequestExecutor.Do(ctx, req, listGroupPushResponse)

	for _, mapping := range listGroupPushResponse.Mappings {
		if mapping.SourceUserGroupId == groupID {
			return &mapping, resp, err
		}
	}

	return nil, resp, err
}

func (m *ApiSupplement) CreateGroupPushMapping(ctx context.Context, appID, groupID, status string) (*GroupPushMapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/internal/instance/%s/grouppush", appID)
	reqBody := map[string]string{
		"status":      status,
		"userGroupId": groupID,
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// API only accessible on admin
	requestAdminURL(req)

	groupPushMapping := &GroupPushMapping{}
	resp, err := m.RequestExecutor.Do(ctx, req, groupPushMapping)
	return groupPushMapping, resp, err
}

func (m *ApiSupplement) UpdateGroupPushMapping(ctx context.Context, appID, mappingID, status string) (*GroupPushMapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/internal/instance/%s/grouppush/%s", appID, mappingID)
	reqBody := map[string]string{
		"status": status,
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// API only accessible on admin
	requestAdminURL(req)

	groupPushMapping := &GroupPushMapping{}
	resp, err := m.RequestExecutor.Do(ctx, req, groupPushMapping)
	return groupPushMapping, resp, err
}

func (m *ApiSupplement) DeleteGroupPushMapping(ctx context.Context, appID, mappingID string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/internal/instance/%s/grouppush/%s/delete", appID, mappingID)
	reqBody := map[string]bool{
		"deleteAppGroup": true,
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	// API only accessible on admin
	requestAdminURL(req)

	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	return resp, err
}
