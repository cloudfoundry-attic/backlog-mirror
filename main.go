package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

const (
	exitVariablesUnset = 1
	exitMirrorFailure  = 2
)

func main() {
	trackerApiToken := os.Getenv("TRACKER_API_TOKEN")
	if trackerApiToken == "" {
		fmt.Fprint(os.Stderr, "TRACKER_API_TOKEN must be set")
		os.Exit(exitVariablesUnset)
	}
	trackerApiClient := gpt.NewClient(trackerApiToken)

	gptStoryService := trackerApiClient.Stories
	ourClient := mirror.NewGoPivotalTrackerWrapper(gptStoryService, trackerApiClient)
	m := mirror.NewMirror(ourClient)

	origBacklog, errO := strconv.Atoi(os.Getenv("TRACKER_ORIG_BACKLOG"))
	destBacklog, errD := strconv.Atoi(os.Getenv("TRACKER_DEST_BACKLOG"))
	if errO != nil || errD != nil {
		os.Exit(exitVariablesUnset)
	}

	err := m.MirrorBacklog(origBacklog, destBacklog)

	if err != nil {
		fmt.Print(fmt.Errorf("error while mirroring:\n %s", err))
		os.Exit(exitMirrorFailure)
	}
}
