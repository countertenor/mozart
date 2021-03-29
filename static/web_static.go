// +build ui

package static

import (
	"embed"
)

const webappBuildType staticType = "webapp/build"

//go:embed webapp/build
const webapp embed.FS

func init() {
	Webapp = webapp
	WebappBuildType = webappBuildType
}
