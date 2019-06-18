package mirror_test

import (
	. "github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gpt "github.com/salsita/go-pivotaltracker/v5/pivotal"
)


func fakeStories() []*gpt.Story {
	story1 := gpt.Story{
		ID:123,
		ProjectID:5,
		Name:"fakeStory1",
	}
	story2 := gpt.Story{
		ID:456,
		ProjectID:5,
		Name:"fakeStory2",
	}

	expectedStories := []*gpt.Story{&story1, &story2}
	return expectedStories
}

type FakeStoryService struct {

}

func (*FakeStoryService) List(int, string) ([]*gpt.Story, error) {
	return fakeStories(), nil
}

var _ = Describe("GetAllStories", func() {

	It("Returns List of stories from api client", func() {

		mockService := &FakeStoryService{}
		wrapper := NewGoPivotalTrackerWrapper(mockService)


		stories := wrapper.GetAllStories(2345567)
		Expect(stories).To(Equal(fakeStories()))

		//assert
		//Expect(thing.callcount).To(Equal(1))
	})
	XIt("Non-mock Returns List of stories from api client", func() {
		//setup

		apiToken := "2b925062ab4acd76cf7cfda319d18158" //TODO: DO NOT COMMIT
		client := gpt.NewClient(apiToken)

		//client.Stories.List()

		//apiClient = gpt.StoryService // mock
		// expectedStoryList = {}
		// when mockStoryService.List() .then(storyList)

		wrapper := NewGoPivotalTrackerWrapper(client.Stories)

		story1 := gpt.Story{
			ID:123,
			ProjectID:5,
			Name:"fakeStory1",
		}
		story2 := gpt.Story{
			ID:456,
			ProjectID:5,
			Name:"fakeStory2",
		}

		expectedStories := []gpt.Story{story1, story2}
		stories := wrapper.GetAllStories(2345567)
		Expect(stories).To(Equal(expectedStories))

		// actualStories = gopivotaltrackerClient.GetAllStories()`


		//assert
		//Expect(thing.callcount).To(Equal(1))

		// expect(actualStories).To(Equal(expectedStoryList))
	})
})