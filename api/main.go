package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NdoleStudio/httpsms/docs"
	"github.com/NdoleStudio/httpsms/pkg/di"
)

// Version is injected at runtime
var Version string

// @title       httpSMS API Reference
// @version     1.0
// @description Use your Android phone to send and receive SMS messages via a simple programmable API with end-to-end encryption. This is a free, self-hosted installation: there are no paid plans, message limits or payments. For guides and examples, see the /docs section of your web dashboard.
//
// @contact.name httpSMS self-hosted (GitHub)
// @contact.url  https://github.com/Ciudad-Valles-Open-Source/httpsms
//
// @license.name AGPL-3.0
// @license.url  https://github.com/Ciudad-Valles-Open-Source/httpsms/blob/main/LICENSE
//
// @BasePath /v1
//
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-Key
func main() {
	if len(os.Args) == 1 {
		di.LoadEnv()
	}

	// Without SWAGGER_HOST the spec has no host/schemes, so Swagger UI sends
	// "Try it out" requests to the same origin that serves the documentation.
	if host := strings.TrimSpace(os.Getenv("SWAGGER_HOST")); len(host) > 0 {
		switch {
		case strings.HasPrefix(host, "http://"):
			docs.SwaggerInfo.Schemes = []string{"http"}
		case strings.HasPrefix(host, "https://"):
			docs.SwaggerInfo.Schemes = []string{"https"}
		}
		host = strings.TrimPrefix(host, "http://")
		host = strings.TrimPrefix(host, "https://")
		docs.SwaggerInfo.Host = strings.TrimSuffix(host, "/")
	}
	if len(Version) > 0 {
		docs.SwaggerInfo.Version = Version
	}

	container := di.NewContainer(os.Getenv("GCP_PROJECT_ID"), Version)
	container.Logger().Info(container.App().Listen(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))).Error())
}
