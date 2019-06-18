package mirror

import gpt "github.com/salsita/go-pivotaltracker/v5/pivotal"

type StoryApiClient struct{}

func (StoryApiClient) GetAllStories(projectId int) gpt.Story {
	panic("implement me")
}
