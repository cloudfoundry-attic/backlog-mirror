package mirror

import (
	"fmt"
	"net/http"

	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

type GoPivotalTrackerWrapper struct {
	storyService GoPivotalTrackerStoryService
	apiClient    TrackerApiClient
}

//go:generate counterfeiter . GoPivotalTrackerStoryService
type GoPivotalTrackerStoryService interface {
	List(int, string) ([]*gpt.Story, error)
	Create(int, *gpt.StoryRequest) (*gpt.Story, *http.Response, error)
}

//go:generate counterfeiter . TrackerApiClient
type TrackerApiClient interface {
	NewRequest(method, urlPath string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*http.Response, error)
}

func NewGoPivotalTrackerWrapper(stories GoPivotalTrackerStoryService, apiClient TrackerApiClient) *GoPivotalTrackerWrapper {
	return &GoPivotalTrackerWrapper{storyService: stories, apiClient: apiClient}
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

func (wrapper *GoPivotalTrackerWrapper) DeleteStory(projectId int, storyId int) error {
	req, err := wrapper.apiClient.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("projects/%d/stories/%d", projectId, storyId),
		"",
	)

	if err != nil {
		return fmt.Errorf("failure while building an http request:\n %s", err)
	}

	req.Header.Del("Content-Type")
	_, err = wrapper.apiClient.Do(req, nil)
	if err != nil {
		return fmt.Errorf("failure while performing http request:\n %s", err)
	}

	return nil
}
