package mirror

import (
	"fmt"

	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

//go:generate counterfeiter . TrackerClient
type TrackerClient interface {
	GetFilteredStories(int, string) ([]*gpt.Story, error)
	AddStoryToProject(int, *gpt.StoryRequest) error
	DeleteStory(backlogId int, storyId int) error
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
	publicLabel := &gpt.Label{Name: "public"}
	request := gpt.StoryRequest{
		Name:        story.Name,
		Type:        story.Type,
		State:       story.State,
		Description: story.Description,
		Labels:      &[]*gpt.Label{publicLabel},
		Estimate: 	 story.Estimate,
		OwnerIDs:    &story.OwnerIDs,
	}
	return &request
}

func (m *Mirror) addAllStoriesToBacklog(stories []*gpt.Story, backlogId int) error {
	for _, story := range stories {
		storyRequest := buildStoryRequest(story)
		err := m.trackerClient.AddStoryToProject(backlogId, storyRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mirror) deleteExistingPublicLabelStoriesFromBacklog(publicLabelStoriesInPublicBacklog []*gpt.Story, backlogId int) error {
	for _, story := range publicLabelStoriesInPublicBacklog {
		err := m.trackerClient.DeleteStory(backlogId, story.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mirror) MirrorBacklog(privateProjectId, publicProjectId int) error {
	publicLabelStories, err := m.trackerClient.GetFilteredStories(privateProjectId, "label:public")
	if err != nil {
		return fmt.Errorf("mirror failed with client error: %s", err)
	}

	publicLabelStoriesInPublicBacklog, err := m.trackerClient.GetFilteredStories(publicProjectId, "label:public")
	if err != nil {
		return fmt.Errorf("mirror failed with client error: %s", err)
	}

	err = m.addAllStoriesToBacklog(publicLabelStories, publicProjectId)
	if err != nil {
		return fmt.Errorf("mirror failed with add-story error: %s", err)
	}

	err = m.deleteExistingPublicLabelStoriesFromBacklog(publicLabelStoriesInPublicBacklog, publicProjectId)
	if err != nil {
		return fmt.Errorf("mirror failed to delete stories:\n %s", err)
	}

	return nil
}
