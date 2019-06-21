package mirror_test

import (
	"errors"
	. "github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror/mirrorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
	"net/url"
)

var _ = Describe("GetFilteredStories", func() {

	var testStories []*gpt.Story

	BeforeEach(func() {
		someLabel := gpt.Label{
			ID: 1,
			ProjectID: 5,
			Name: "someLabel",
		}
		story1 := gpt.Story{
			ID:123,
			ProjectID:5,
			Name:"fakeStory1",
			Labels: []*gpt.Label{&someLabel},

		}
		story2 := gpt.Story{
			ID:456,
			ProjectID:5,
			Name:"fakeStory2",
			Labels: []*gpt.Label{&someLabel},
		}
		testStories = []*gpt.Story{&story1, &story2}
	})

	It("calls client with the projectId and filter", func() {
		mockService := &mirrorfakes.FakeGoPivotalTrackerStoryService{}

		wrapper := NewGoPivotalTrackerWrapper(mockService)

		_,_ = wrapper.GetFilteredStories(0, "someFilter")

		Expect(mockService.ListCallCount()).To(Equal(1))
		projectIdArg, filterArg := mockService.ListArgsForCall(0)
		Expect(projectIdArg).To(Equal(0))
		Expect(filterArg).To(Equal("someFilter"))
	})

	It("returns filtered stories", func() {
		mockService := &mirrorfakes.FakeGoPivotalTrackerStoryService{}
		mockService.ListReturns(testStories, nil)
		wrapper := NewGoPivotalTrackerWrapper(mockService)

		stories, err := wrapper.GetFilteredStories(0, "someFilter")

		Expect(err).NotTo(HaveOccurred())
		Expect(stories).To(Equal(testStories))
	})

	It ("returns an error when List does", func() {
		mockService := &mirrorfakes.FakeGoPivotalTrackerStoryService{}
		mockService.ListReturns(nil, errors.New("some list error"))
		wrapper := NewGoPivotalTrackerWrapper(mockService)

		_, err := wrapper.GetFilteredStories(0, "someFilter")

		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("AddStoryToProject", func() {
	It("Calls Create with the correct arguments", func(){
		mockService := &mirrorfakes.FakeGoPivotalTrackerStoryService{}
		projectId := 4

		wrapper := NewGoPivotalTrackerWrapper(mockService)
		storyRequest := gpt.StoryRequest{}

		_ = wrapper.AddStoryToProject(projectId, &storyRequest)

		Expect(mockService.CreateCallCount()).To(Equal(1))
		projectIdArg, storyRequestArg := mockService.CreateArgsForCall(0)
		Expect(projectIdArg).To(Equal(projectId))
		Expect(*storyRequestArg).To(Equal(storyRequest))
	})

	It("returns an error when storyServiceCreate returns an error", func() {
		mockService := &mirrorfakes.FakeGoPivotalTrackerStoryService{}
		mockService.CreateReturns(nil, nil, &url.Error{"POST", "http://example.com", errors.New("")})
		projectId := 4
		wrapper := NewGoPivotalTrackerWrapper(mockService)
		storyRequest := gpt.StoryRequest{}

		err := wrapper.AddStoryToProject(projectId, &storyRequest)

		Expect(err).To(HaveOccurred())
	})
})
