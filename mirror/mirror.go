package mirror


//go:generate counterfeiter . StoryApi
type StoryApi interface {
	GetAllStories(projectId int) Story
}

type Story struct {
	id int
	projectId int
	name string
}

type Mirror struct {
	storyApi StoryApi
}

func NewMirror(givenStoryApi StoryApi) *Mirror {
	return &Mirror{
		givenStoryApi,
	}
}

//Big function. Woah.
func (m *Mirror) MirrorBacklog() {
	_ = m.storyApi.GetAllStories(100)
}