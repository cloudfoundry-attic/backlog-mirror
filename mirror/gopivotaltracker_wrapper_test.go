package mirror_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	. "github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror/mirrorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

var _ = Describe("GetFilteredStories", func() {

	var testStories []*gpt.Story

	BeforeEach(func() {
		someLabel := gpt.Label{
			ID:        1,
			ProjectID: 5,
			Name:      "someLabel",
		}
		story1 := gpt.Story{
			ID:        123,
			ProjectID: 5,
			Name:      "fakeStory1",
			Labels:    []*gpt.Label{&someLabel},
		}
		story2 := gpt.Story{
			ID:        456,
			ProjectID: 5,
			Name:      "fakeStory2",
			Labels:    []*gpt.Label{&someLabel},
		}
		testStories = []*gpt.Story{&story1, &story2}
	})

	It("calls client with the projectId and filter", func() {
		storyService := new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		client := new(mirrorfakes.FakeTrackerApiClient)

		wrapper := NewGoPivotalTrackerWrapper(storyService, client)

		_, _ = wrapper.GetFilteredStories(0, "someFilter")

		Expect(storyService.ListCallCount()).To(Equal(1))
		projectIdArg, filterArg := storyService.ListArgsForCall(0)
		Expect(projectIdArg).To(Equal(0))
		Expect(filterArg).To(Equal("someFilter"))
	})

	It("returns filtered stories", func() {
		storyService := new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		client := new(mirrorfakes.FakeTrackerApiClient)
		storyService.ListReturns(testStories, nil)
		wrapper := NewGoPivotalTrackerWrapper(storyService, client)

		stories, err := wrapper.GetFilteredStories(0, "someFilter")

		Expect(err).NotTo(HaveOccurred())
		Expect(stories).To(Equal(testStories))
	})

	It("returns an error when List does", func() {
		storyService := new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		client := new(mirrorfakes.FakeTrackerApiClient)
		storyService.ListReturns(nil, errors.New("some list error"))
		wrapper := NewGoPivotalTrackerWrapper(storyService, client)

		_, err := wrapper.GetFilteredStories(0, "someFilter")

		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("AddStoryToProject", func() {
	It("Calls Create with the correct arguments", func() {
		storyService := new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		client := new(mirrorfakes.FakeTrackerApiClient)
		projectId := 4

		wrapper := NewGoPivotalTrackerWrapper(storyService, client)
		storyRequest := gpt.StoryRequest{}

		_ = wrapper.AddStoryToProject(projectId, &storyRequest)

		Expect(storyService.CreateCallCount()).To(Equal(1))
		projectIdArg, storyRequestArg := storyService.CreateArgsForCall(0)
		Expect(projectIdArg).To(Equal(projectId))
		Expect(*storyRequestArg).To(Equal(storyRequest))
	})

	It("returns an error when storyServiceCreate returns an error", func() {
		storyService := new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		client := new(mirrorfakes.FakeTrackerApiClient)
		storyService.CreateReturns(nil, nil, &url.Error{"POST", "http://example.com", errors.New("")})
		projectId := 4
		wrapper := NewGoPivotalTrackerWrapper(storyService, client)
		storyRequest := gpt.StoryRequest{}

		err := wrapper.AddStoryToProject(projectId, &storyRequest)

		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("DeleteStory", func() {
	var (
		storyService *mirrorfakes.FakeGoPivotalTrackerStoryService
		apiClient    *mirrorfakes.FakeTrackerApiClient
		wrapper      *GoPivotalTrackerWrapper
	)

	BeforeEach(func() {
		storyService = new(mirrorfakes.FakeGoPivotalTrackerStoryService)
		apiClient = new(mirrorfakes.FakeTrackerApiClient)
		wrapper = NewGoPivotalTrackerWrapper(storyService, apiClient)

		apiClient.NewRequestReturns(&http.Request{}, nil)
	})

	It("calls delete correctly, without Content-Type header", func() {
		mockRequestHeader := map[string][]string{
			http.CanonicalHeaderKey("Content-Type"):   {"application/json"},
			http.CanonicalHeaderKey("X-Trackertoken"): {"somevalue123"},
		}
		mockDeleteRequest := http.Request{
			Method: "DELETE",
			Header: mockRequestHeader,
			Host:   "example.com",
		}

		apiClient.NewRequestReturns(&mockDeleteRequest, nil)
		expectedDeleteRequest := mockDeleteRequest
		expectedDeleteRequest.Header = map[string][]string{
			http.CanonicalHeaderKey("X-Trackertoken"): {"somevalue123"},
		}

		projectId := 1
		storyId := 2
		err := wrapper.DeleteStory(projectId, storyId)
		Expect(err).NotTo(HaveOccurred())

		Expect(apiClient.NewRequestCallCount()).To(Equal(1), "expected to make call to API client NewRequest")
		meth, path, body := apiClient.NewRequestArgsForCall(0)
		Expect(meth).To(Equal(http.MethodDelete))
		Expect(path).To(Equal(fmt.Sprintf("projects/%d/stories/%d", projectId, storyId)))
		Expect(body).To(BeEmpty())

		Expect(apiClient.DoCallCount()).To(Equal(1), "expected to make call to API client Do")
		actualRequest, responseSideEffectObject := apiClient.DoArgsForCall(0)
		Expect(actualRequest).To(BeEquivalentTo(&expectedDeleteRequest))
		Expect(responseSideEffectObject).To(BeNil())
	})

	It("returns error when API client newRequest does", func() {
		errorMessage := "i failed to build you a request"
		apiClient.NewRequestReturns(nil, errors.New(errorMessage))

		err := wrapper.DeleteStory(1, 2)
		Expect(err).To(MatchError(fmt.Sprintf("failure while building an http request:\n %s", errorMessage)))
	})

	It("returns an error when API client Do does", func() {
		errorMessage := "i failed to do your request"
		apiClient.DoReturns(&http.Response{}, errors.New(errorMessage))

		err := wrapper.DeleteStory(1, 2)
		Expect(err).To(MatchError(fmt.Sprintf("failure while performing http request:\n %s", errorMessage)))
	})
})
