/*
Copyright 2022 The Kubernetes Authors.

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

package settings

import (
	"os"
	"strings"

	"github.com/onsi/ginkgo/v2"

	"k8s.io/ingress-nginx/test/e2e/framework"
)

const (
	enableOpentelemetry            = "enable-opentelemetry"
	opentelemetryTrustIncomingSpan = "opentelemetry-trust-incoming-span"

	opentelemetryOperationName         = "opentelemetry-operation-name"
	opentelemetryLocationOperationName = "opentelemetry-location-operation-name"
	opentelemetryConfig                = "opentelemetry-config"
	opentelemetryConfigPath            = "/etc/ingress-controller/telemetry/opentelemetry.toml"

	enable = "true"
)

var _ = framework.IngressNginxDescribe("Configure Opentelemetry", func() {
	f := framework.NewDefaultFramework("enable-opentelemetry")

	shouldSkip := false
	skip, ok := os.LookupEnv("SKIP_OPENTELEMETRY_TESTS")
	if ok && skip == enable {
		shouldSkip = true
	}

	ginkgo.BeforeEach(func() {
		f.NewEchoDeployment()
	})

	ginkgo.AfterEach(func() {
	})

	ginkgo.It("should not exists opentelemetry directive", func() {
		if shouldSkip {
			ginkgo.Skip("skipped")
		}
		config := map[string]string{}
		config[enableOpentelemetry] = disable
		f.SetNginxConfigMapData(config)

		f.EnsureIngress(framework.NewSingleIngress(enableOpentelemetry, "/", enableOpentelemetry, f.Namespace, "http-svc", 80, nil))

		f.WaitForNginxConfiguration(
			func(cfg string) bool {
				return !strings.Contains(cfg, "opentelemetry on")
			})
	})

	ginkgo.It("should exists opentelemetry directive when is enabled", func() {
		if shouldSkip {
			ginkgo.Skip("skipped")
		}
		config := map[string]string{}
		config[enableOpentelemetry] = enable
		config[opentelemetryConfig] = opentelemetryConfigPath
		f.SetNginxConfigMapData(config)

		f.EnsureIngress(framework.NewSingleIngress(enableOpentelemetry, "/", enableOpentelemetry, f.Namespace, "http-svc", 80, nil))

		f.WaitForNginxConfiguration(
			func(cfg string) bool {
				return strings.Contains(cfg, "opentelemetry on")
			})
	})

	ginkgo.It("should include opentelemetry_trust_incoming_spans on directive when enabled", func() {
		if shouldSkip {
			ginkgo.Skip("skipped")
		}
		config := map[string]string{}
		config[enableOpentelemetry] = enable
		config[opentelemetryConfig] = opentelemetryConfigPath
		config[opentelemetryTrustIncomingSpan] = enable
		f.SetNginxConfigMapData(config)

		f.EnsureIngress(framework.NewSingleIngress(enableOpentelemetry, "/", enableOpentelemetry, f.Namespace, "http-svc", 80, nil))

		f.WaitForNginxConfiguration(
			func(cfg string) bool {
				return strings.Contains(cfg, "opentelemetry_trust_incoming_spans on")
			})
	})

	ginkgo.It("should not exists opentelemetry_operation_name directive when is empty", func() {
		if shouldSkip {
			ginkgo.Skip("skipped")
		}
		config := map[string]string{}
		config[enableOpentelemetry] = enable
		config[opentelemetryConfig] = opentelemetryConfigPath
		config[opentelemetryOperationName] = ""
		f.SetNginxConfigMapData(config)

		f.EnsureIngress(framework.NewSingleIngress(enableOpentelemetry, "/", enableOpentelemetry, f.Namespace, "http-svc", 80, nil))

		f.WaitForNginxConfiguration(
			func(cfg string) bool {
				return !strings.Contains(cfg, "opentelemetry_operation_name")
			})
	})

	ginkgo.It("should exists opentelemetry_operation_name directive when is configured", func() {
		if shouldSkip {
			ginkgo.Skip("skipped")
		}
		config := map[string]string{}
		config[enableOpentelemetry] = enable
		config[opentelemetryConfig] = opentelemetryConfigPath
		config[opentelemetryOperationName] = "HTTP $request_method $uri"
		f.SetNginxConfigMapData(config)

		f.EnsureIngress(framework.NewSingleIngress(enableOpentelemetry, "/", enableOpentelemetry, f.Namespace, "http-svc", 80, nil))

		f.WaitForNginxConfiguration(
			func(cfg string) bool {
				return strings.Contains(cfg, `opentelemetry_operation_name "HTTP $request_method $uri"`)
			})
	})
})
