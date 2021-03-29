// +build ui

package static

import (
	"embed"
)

var webappBuildType staticType = "webapp/build"

//go:embed webapp/build
var webapp embed.FS

func init() {
	Webapp = webapp
	WebappBuildType = webappBuildType
}
