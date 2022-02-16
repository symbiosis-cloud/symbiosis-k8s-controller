/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	capiv1 "k8s.io/api/certificates/v1"
)

// Source(10/2021): https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/certificates/certificate_controller_utils.go

// GetCertApprovalCondition returns the current condition of the CSR (approved, denied)
func GetCertApprovalCondition(status *capiv1.CertificateSigningRequestStatus) (approved, denied bool) {
	for _, c := range status.Conditions {
		if c.Type == capiv1.CertificateApproved {
			approved = true
		}

		if c.Type == capiv1.CertificateDenied {
			denied = true
		}
	}

	return
}

// Source(10/2021): https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/certificates/helpers.go

// ParseCSR extracts the CSR from the bytes and decodes it.
func ParseCSR(pemBytes []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("PEM block type must be CERTIFICATE REQUEST")
	}

	csr, err := x509.ParseCertificateRequest(block.Bytes)

	if err != nil {
		return nil, err
	}

	return csr, nil
}
