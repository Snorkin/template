package git

import (
	"fmt"
	"runtime/debug"
)

var (
	commitTag    string
	commitAuthor string
)

type CommitInfo struct {
	Revision string
	Time     string
	Author   string
}

func (c CommitInfo) String() string {
	return fmt.Sprintf("%s at %s by %s", c.Revision, c.Time, c.Author)
}

func GetCommitInfo() CommitInfo {
	var revision, revisionTime string

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				revision = setting.Value
			}

			if setting.Key == "vcs.time" {
				revisionTime = setting.Value
			}
		}
	}

	if commitTag != "" {
		revision = commitTag
	}

	return CommitInfo{
		Revision: revision,
		Time:     revisionTime,
		Author:   commitAuthor,
	}
}
