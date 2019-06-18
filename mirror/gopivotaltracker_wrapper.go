package mirror

import(
	gpt "github.com/salsita/go-pivotaltracker/v5/pivotal"
)

type GoPivotalTrackerWrapper struct {
	//client *gpt.Client
	storyService GoPivotalTrackerStories
}

//type GoPivotalTrackerClient interface {
//	//Stories *gpt.StoryService
//}

//go:generate counterfeiter . GoPivotalTrackerStories
type GoPivotalTrackerStories interface {
	List(int, string) ([]*gpt.Story, error)
}


func NewGoPivotalTrackerWrapper(stories GoPivotalTrackerStories) *GoPivotalTrackerWrapper {
	return &GoPivotalTrackerWrapper{stories}
}

func (wrapper *GoPivotalTrackerWrapper) GetAllStories(projectId int) []*gpt.Story {

	stories, _ := wrapper.storyService.List(projectId, "")
	return stories
}
