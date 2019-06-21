package main

import (
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
	"os"
	"strconv"
)

const(
	exitVariablesUnset = 1
	exitMirrorFailure = 2
)

func main() {
	gptStoryService := gpt.NewClient(os.Getenv("TRACKER_API_TOKEN")).Stories
	client := mirror.NewGoPivotalTrackerWrapper(gptStoryService)
	m := mirror.NewMirror(client)

	origBacklog, errO := strconv.Atoi(os.Getenv("TRACKER_ORIG_BACKLOG"))
	destBacklog, errD := strconv.Atoi(os.Getenv("TRACKER_DEST_BACKLOG"))
	if errO != nil || errD != nil {
		os.Exit(exitVariablesUnset)
	}

	err := m.MirrorBacklog(origBacklog, destBacklog)

	if err != nil {
		os.Exit(exitMirrorFailure)
	}
}
