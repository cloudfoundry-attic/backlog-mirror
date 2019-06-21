package mirror

import (
	"fmt"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

//go:generate counterfeiter . TrackerClient
type TrackerClient interface {
	GetFilteredStories(int, string) ([]*gpt.Story, error)
	AddStoryToProject(int, *gpt.StoryRequest) error
}

type Mirror struct {
	trackerClient TrackerClient
}

func NewMirror(givenClient TrackerClient) *Mirror {
	return &Mirror{
		givenClient,
	}
}

func buildStoryRequest(story *gpt.Story) *gpt.StoryRequest {
	request := gpt.StoryRequest{
		Name: story.Name,
		Type: story.Type,
		State: story.State,
		Description: story.Description,
	}
	return &request
}

func (m *Mirror) MirrorBacklog(privateProjectId, publicProjectId int) error{
	publicLabelStories, err := m.trackerClient.GetFilteredStories(privateProjectId, "label:public")
	if err != nil {
		return fmt.Errorf("mirror failed with client error: %s", err)
	}

	for i:=0;i< len(publicLabelStories);i++ {
		story := publicLabelStories[i]
		storyRequest := buildStoryRequest(story)
		err := m.trackerClient.AddStoryToProject(publicProjectId, storyRequest)
		if err != nil {
			return err
		}
	}
	return nil
}
