//go:build debug

package cmd

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth" // permits to use all authentication providers
)
