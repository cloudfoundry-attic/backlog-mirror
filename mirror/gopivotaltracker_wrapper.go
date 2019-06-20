package mirror

import (
	"fmt"
	gpt "github.com/salsita/go-pivotaltracker/v5/pivotal"
	"net/http"
)

type GoPivotalTrackerWrapper struct {
	//client *gpt.Client
	storyService GoPivotalTrackerStoryService
}

//go:generate counterfeiter . GoPivotalTrackerStoryService
type GoPivotalTrackerStoryService interface {
	List(int, string) ([]*gpt.Story, error)
	Create(int, *gpt.StoryRequest) (*gpt.Story, *http.Response, error)
}


func NewGoPivotalTrackerWrapper(stories GoPivotalTrackerStoryService) *GoPivotalTrackerWrapper {
	return &GoPivotalTrackerWrapper{stories}
}

func (wrapper *GoPivotalTrackerWrapper) GetFilteredStories(projectId int, filter string) ([]*gpt.Story, error) {
	stories, err := wrapper.storyService.List(projectId, filter)
	if err != nil {
		return nil, fmt.Errorf("tracker API client could not list stories: %s", err)
	}
	return stories, nil
}

func (wrapper *GoPivotalTrackerWrapper) AddStoryToProject(projectId int, request *gpt.StoryRequest) error {
	_, _, err := wrapper.storyService.Create(projectId, request)

	if err != nil {
		return err
	}
	return nil
}
