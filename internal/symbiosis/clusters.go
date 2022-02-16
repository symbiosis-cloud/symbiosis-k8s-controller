package symbiosis

import (
	"context"
	"fmt"
	"net/http"
)

type ClusterService interface {
	GetByID(context.Context, string) (*Cluster, *http.Response, error)
}

type ClusterServiceOp struct {
	client *Client
}

type ClusterNode struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	PrivateIpv4Address string `json:"privateIpv4Address"`
}

type Cluster struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Nodes []ClusterNode `json:"nodes"`
}

func (svc *ClusterServiceOp) GetByID(ctx context.Context, id string) (*Cluster, *http.Response, error) {
	path := fmt.Sprintf("/rest/v1/cluster/by-id/%v", id)

	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	cluster := &Cluster{}
	resp, err := svc.client.Do(ctx, req, cluster)

	if err != nil {
		return nil, resp, err
	}

	return cluster, resp, nil
}
