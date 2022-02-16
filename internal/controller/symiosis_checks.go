package controller

import (
	"context"
	"strings"
	"symbiosis-cloud/symbiosis-k8s-controller/internal/symbiosis"

	"crypto/x509"
	"fmt"
	"net"

	certificatesv1 "k8s.io/api/certificates/v1"
)

// ProviderChecks is a function in which the Cloud Provider specifies a series of checks
// to run against the CSRs. The out-of-band identity checks of the CSRs should happen here
func (r *CertificateSigningRequestReconciler) ProviderChecks(ctx context.Context, csr *certificatesv1.CertificateSigningRequest, x509csr *x509.CertificateRequest) (success bool, err string) {
	if len(x509csr.DNSNames) > 1 {
		return false, "The x509 Cert Request contains more than 1 DNS name"
	}

	sanDNSName := x509csr.DNSNames[0]

	hostname := strings.TrimPrefix(csr.Spec.Username, "system:node:")

	if !strings.HasPrefix(sanDNSName, hostname) {
		return false, fmt.Sprintf("The SAN DNS Name %s in the x509 CSR is not prefixed by the node name %s (hostname)", sanDNSName, hostname)
	}

	cluster, _, _ := r.SymbiosisClient.Clusters.GetByID(ctx, r.ClusterID)

	var node *symbiosis.ClusterNode

	for i, n := range cluster.Nodes {
		if n.Name == hostname {
			node = &cluster.Nodes[i]
			break
		}
	}

	if node == nil {
		return false, fmt.Sprintf("No symbiosis node with name %v found", hostname)
	}

	sanIPAddrs := x509csr.IPAddresses
	for _, ip := range sanIPAddrs {
		nodeIP := net.ParseIP(node.PrivateIpv4Address)
		if !nodeIP.Equal(ip) {
			return false, "SAN IP is not equal to private node IP"
		}
	}

	return true, ""
}
