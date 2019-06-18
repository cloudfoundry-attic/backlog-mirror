package mirror

import gpt "github.com/salsita/go-pivotaltracker/v5/pivotal"

//go:generate counterfeiter . StoryApi
type StoryApi interface {
	GetAllStories(projectId int) *[]gpt.Story
}

type Mirror struct {
	storyApi StoryApi
}

func NewMirror(givenStoryApi StoryApi) *Mirror {
	return &Mirror{
		givenStoryApi,
	}
}

func (m *Mirror) MirrorBacklog() {
	_ = m.storyApi.GetAllStories(100)
}