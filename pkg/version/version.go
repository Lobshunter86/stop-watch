package version

import (
	"fmt"
	"runtime"
)

var (
	GitHash   = "unknown"
	GitBranch = "unknown"
	BuildDate = "unknown"
)

func Version() string {
	return fmt.Sprintf("Go Version: %s\nGit Branch: %s\nGitHash: %s\nBuild UTC Time: %s", runtime.Version(), GitBranch, GitHash, BuildDate)
}
