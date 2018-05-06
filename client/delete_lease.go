package client

import (
	"net/http"

	"github.com/teamhephy/k8s-claimer/htp"
)

// DeleteLeaseReq is the encoding/json compatible struct that represents the DELETE /lease request body
type DeleteLeaseReq struct {
	CloudProvider string `json:"cloud_provider"`
}

// DeleteLease deletes a lease
func DeleteLease(server, authToken, cloudProvider, leaseToken string) error {
	endpt := newEndpoint(htp.Delete, server, "lease/"+cloudProvider+"/"+leaseToken)
	resp, err := endpt.executeReq(getHTTPClient(), nil, authToken)
	if err != nil {
		return errHTTPRequest{endpoint: endpt.String(), err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return APIError{endpoint: endpt.String(), code: resp.StatusCode}
	}
	return nil
}
