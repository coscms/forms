package defaults

import (
	"embed"

	"github.com/coscms/forms/common"
)

//go:embed templates
var templateFS embed.FS

func init() {
	common.FileSystem.Register(templateFS)
}
