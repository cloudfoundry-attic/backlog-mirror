package mirror_test

import (
	"errors"
	. "github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror/mirrorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

var _ = Describe("Backlog Mirror", func() {

	var privateProjectId int
	var publicProjectId int
	var trackerClient *mirrorfakes.FakeTrackerClient
	var publicStories []*gpt.Story

	BeforeEach(func() {
		privateProjectId = 7
		publicProjectId = 8
		trackerClient = &mirrorfakes.FakeTrackerClient{}
	})

	It("gets public-labeled stories for a project", func() {
		mirror := *NewMirror(trackerClient)

		_ = mirror.MirrorBacklog(privateProjectId, publicProjectId)

		Expect(trackerClient.GetFilteredStoriesCallCount()).To(Equal(1))
		originProjectIdArg, storyFilterArg := trackerClient.GetFilteredStoriesArgsForCall(0)
		Expect(originProjectIdArg).To(Equal(privateProjectId))
		Expect(storyFilterArg).To(Equal("label:public"))
	})

	It("returns error if the tracker client does", func(){
		trackerClient.GetFilteredStoriesReturns(nil, errors.New("fake client error"))
		mirror := *NewMirror(trackerClient)

		err := mirror.MirrorBacklog(privateProjectId, publicProjectId)

		Expect(err).To(HaveOccurred())
	})

	Describe("pushes public-labeled stories to public backlog", func() {

		setup := func() {
			privateProjectId = 7
			publicProjectId = 8
			trackerClient = &mirrorfakes.FakeTrackerClient{}
			story1 := gpt.Story{
				ID:          123,
				ProjectID:   privateProjectId,
				Name:        "fakeStory1",
				Type:        gpt.StoryTypeChore,
				State:       gpt.StoryStateStarted,
				Description: "description1",
			}
			story2 := gpt.Story{
				ID:          456,
				ProjectID:   privateProjectId,
				Name:        "fakeStory2",
				Type:        gpt.StoryTypeFeature,
				State:       gpt.StoryStatePlanned,
				Description: "description2",
			}
			publicStories = []*gpt.Story{&story1, &story2}
			trackerClient.GetFilteredStoriesReturns(publicStories, nil)
		}

		BeforeEach(func() {
			setup()
			mirror := *NewMirror(trackerClient)

			err := mirror.MirrorBacklog(privateProjectId, publicProjectId)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("when adding stories fails", func() {
			BeforeEach(func() {
				setup()
				trackerClient.AddStoryToProjectReturns(errors.New("failed to add"))
			})
			It("returns an error", func() {
				mirror := *NewMirror(trackerClient)
				err := mirror.MirrorBacklog(privateProjectId, publicProjectId)

				Expect(err).To(HaveOccurred())
			})
		})

		It("correct number of stories", func() {
			numPublicStories := len(publicStories)
			Expect(trackerClient.AddStoryToProjectCallCount()).To(Equal(numPublicStories))
		})

		It("mirrored stories have the correct story names", func() {
			_, actualStory1 := trackerClient.AddStoryToProjectArgsForCall(0)
			_, actualStory2 := trackerClient.AddStoryToProjectArgsForCall(1)

			Expect(actualStory1.Name).To(Equal(publicStories[0].Name))
			Expect(actualStory2.Name).To(Equal(publicStories[1].Name))
		})

		It("mirrored stories have the correct project ids", func() {
			actualID1, _ := trackerClient.AddStoryToProjectArgsForCall(0)
			actualID2, _ := trackerClient.AddStoryToProjectArgsForCall(1)

			Expect(actualID1).To(Equal(publicProjectId))
			Expect(actualID2).To(Equal(publicProjectId))
		})

		It("mirrored stories have the correct story Type", func() {
			_, actualStory1 := trackerClient.AddStoryToProjectArgsForCall(0)
			_, actualStory2 := trackerClient.AddStoryToProjectArgsForCall(1)

			Expect(actualStory1.Type).To(Equal(gpt.StoryTypeChore))
			Expect(actualStory2.Type).To(Equal(gpt.StoryTypeFeature))
		})

		It("mirrored stories have the correct story State", func() {
			_, actualStory1 := trackerClient.AddStoryToProjectArgsForCall(0)
			_, actualStory2 := trackerClient.AddStoryToProjectArgsForCall(1)

			Expect(actualStory1.State).To(Equal(gpt.StoryStateStarted))
			Expect(actualStory2.State).To(Equal(gpt.StoryStatePlanned))
		})

		It("mirrored stories have the correct story Description", func() {
			_, actualStory1 := trackerClient.AddStoryToProjectArgsForCall(0)
			_, actualStory2 := trackerClient.AddStoryToProjectArgsForCall(1)

			Expect(actualStory1.Description).To(Equal(publicStories[0].Description))
			Expect(actualStory2.Description).To(Equal(publicStories[1].Description))
		})
	})


	Describe("deletes existing stories in public backlog", func() {

		//var publicLabeledStoriesAlreadyInPublicBacklog []*gpt.Story
		//var storiesInsd

		setup := func() {
			privateProjectId = 7
			publicProjectId = 8
			trackerClient = &mirrorfakes.FakeTrackerClient{}
			story1 := gpt.Story{
				ID:          123,
				ProjectID:   privateProjectId,
				Name:        "fakeStory1",
				Type:        gpt.StoryTypeChore,
				State:       gpt.StoryStateStarted,
				Description: "description1",
			}
			story2 := gpt.Story{
				ID:          456,
				ProjectID:   privateProjectId,
				Name:        "fakeStory2",
				Type:        gpt.StoryTypeFeature,
				State:       gpt.StoryStatePlanned,
				Description: "description2",
			}
			publicStories = []*gpt.Story{&story1, &story2}
			//trackerClient.GetFilteredStoriesReturns(publicStories, nil)

			//trackerClient.GetFilteredStoriesStub = func(i int, s string) ([]*gpt.Story, error){
			//	if int i=
			//}

		}

		BeforeEach(func() {
			setup()
		})

		FIt("deletes existing stories in public backlog", func() {
			privateProjectId := 7
			publicProjectId := 8
			trackerClient := &mirrorfakes.FakeTrackerClient{}

			//publicStories := nil

			trackerClient.GetFilteredStoriesReturns(publicStories, nil)

			//trackerClient.GetFilteredStoriesStub = func() {
			//
			//}

			mirror := *NewMirror(trackerClient)

			_ = mirror.MirrorBacklog(privateProjectId, publicProjectId)

			//correct number of calls
			Expect(trackerClient.DeleteStoryCallCount()).To(BeNumerically(">=", 1))
		})
	})

})
