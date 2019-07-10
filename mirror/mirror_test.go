package mirror_test

import (
	"errors"
	"fmt"
	. "github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror/mirrorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gpt "gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
	"math/rand"
)

var _ = Describe("MirrorBacklog", func() {

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

		var privateCallCount int

		trackerClient.GetFilteredStoriesStub = func(i int, s string) ([]*gpt.Story, error) {
			if i == privateProjectId {
				privateCallCount++
			}
			return nil, nil
		}

		err := mirror.MirrorBacklog(privateProjectId, publicProjectId)
		Expect(err).NotTo(HaveOccurred())

		Expect(privateCallCount).To(Equal(1))
		originProjectIdArg, storyFilterArg := trackerClient.GetFilteredStoriesArgsForCall(0)
		Expect(originProjectIdArg).To(Equal(privateProjectId))
		Expect(storyFilterArg).To(Equal("label:public"))
	})

	It("returns error if the tracker client does", func() {
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
			publicLabelInPrivateBacklog := gpt.Label{
				ID:        999,
				ProjectID: 7,
				Name:      "public",
			}
			story1 := gpt.Story{
				ID:          123,
				ProjectID:   privateProjectId,
				Name:        "fakeStory1",
				Type:        gpt.StoryTypeChore,
				State:       gpt.StoryStateStarted,
				Description: "description1",
				Labels:      []*gpt.Label{&publicLabelInPrivateBacklog},
			}
			story2 := gpt.Story{
				ID:          456,
				ProjectID:   privateProjectId,
				Name:        "fakeStory2",
				Type:        gpt.StoryTypeFeature,
				State:       gpt.StoryStatePlanned,
				Description: "description2",
				Labels:      []*gpt.Label{&publicLabelInPrivateBacklog},
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

		It("mirrored stories have a correct public label", func() {
			_, actualStory1 := trackerClient.AddStoryToProjectArgsForCall(0)
			_, actualStory2 := trackerClient.AddStoryToProjectArgsForCall(1)

			expectedLabel := &gpt.Label{Name: "public"}

			Expect(*actualStory1.Labels).To(ContainElement(expectedLabel))
			Expect(*actualStory2.Labels).To(ContainElement(expectedLabel))
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

		var publicLabeledStoriesAlreadyInPublicBacklog []*gpt.Story

		setup := func() {
			privateProjectId = 7
			publicProjectId = 8
			trackerClient = &mirrorfakes.FakeTrackerClient{}

			story1 := gpt.Story{
				ID:          789,
				ProjectID:   privateProjectId,
				Name:        "fakeStory1",
				Type:        gpt.StoryTypeChore,
				State:       gpt.StoryStateStarted,
				Description: "description1",
			}
			story2 := gpt.Story{
				ID:          234,
				ProjectID:   privateProjectId,
				Name:        "fakeStory2",
				Type:        gpt.StoryTypeFeature,
				State:       gpt.StoryStatePlanned,
				Description: "description2",
			}

			publicLabeledStoriesAlreadyInPublicBacklog = []*gpt.Story{&story1, &story2}

			story3 := gpt.Story{
				ID:          189,
				ProjectID:   privateProjectId,
				Name:        "fakeStory1",
				Type:        gpt.StoryTypeChore,
				State:       gpt.StoryStateStarted,
				Description: "description1",
			}
			story4 := gpt.Story{
				ID:          134,
				ProjectID:   privateProjectId,
				Name:        "fakeStory2",
				Type:        gpt.StoryTypeFeature,
				State:       gpt.StoryStatePlanned,
				Description: "description2",
			}

			publicLabeledStoriesInPrivateBacklog := []*gpt.Story{&story3, &story4}

			trackerClient.GetFilteredStoriesStub = func(i int, s string) ([]*gpt.Story, error) {
				if i == privateProjectId {
					return publicLabeledStoriesInPrivateBacklog, nil
				}
				if i == publicProjectId {
					return publicLabeledStoriesAlreadyInPublicBacklog, nil
				}

				return []*gpt.Story{}, nil
			}

			trackerClient.AddStoryToProjectStub = func(i int, s *gpt.StoryRequest) error {

				publicLabeledStoriesAlreadyInPublicBacklog = append(publicLabeledStoriesAlreadyInPublicBacklog, &gpt.Story{
					ID:          rand.Intn(200),
					ProjectID:   publicProjectId,
					Name:        s.Name,
					Type:        s.Type,
					State:       s.State,
					Description: s.Description,
				})
				return nil
			}
		}

		BeforeEach(func() {
			setup()
		})

		It("deletes existing stories in public backlog", func() {

			mirror := *NewMirror(trackerClient)

			_ = mirror.MirrorBacklog(privateProjectId, publicProjectId)

			//correct number of calls
			Expect(trackerClient.DeleteStoryCallCount()).To(Equal(2), "did not delete both stories")

			deletionProjectId, deletionStoryId := trackerClient.DeleteStoryArgsForCall(0)
			Expect(deletionProjectId).To(Equal(publicProjectId))
			Expect(deletionStoryId).To(Equal(publicLabeledStoriesAlreadyInPublicBacklog[0].ID))

			deletionProjectId, deletionStoryId = trackerClient.DeleteStoryArgsForCall(1)
			Expect(deletionProjectId).To(Equal(publicProjectId))
			Expect(deletionStoryId).To(Equal(publicLabeledStoriesAlreadyInPublicBacklog[1].ID))
		})

		It("gets stories-to-delete from public backlog before adding private-backlog stories", func() {

			mirror := *NewMirror(trackerClient)

			err := mirror.MirrorBacklog(privateProjectId, publicProjectId)
			Expect(err).ToNot(HaveOccurred())

			trackerClientCalls := trackerClient.Invocations()
			deleteCalls := trackerClientCalls["DeleteStory"]
			Expect(len(deleteCalls)).To(Equal(2))
			Expect(deleteCalls).Should(ContainElement(ConsistOf(publicProjectId, 789)))
			Expect(deleteCalls).Should(ContainElement(ConsistOf(publicProjectId, 234)))
		})

		It("returns error if the tracker client errors on retrieval of public stories", func() {
			trackerClient.GetFilteredStoriesStub = func(i int, s string) ([]*gpt.Story, error) {
				if i == publicProjectId {
					return nil, errors.New("fake retrieval error")
				}

				return []*gpt.Story{}, nil
			}
			mirror := *NewMirror(trackerClient)

			err := mirror.MirrorBacklog(privateProjectId, publicProjectId)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("mirror failed with client error: %s", "fake retrieval error")))
		})

		It("returns error if the tracker client errors on story deletion", func() {
			trackerClient.DeleteStoryReturns(errors.New("fake client error"))
			mirror := *NewMirror(trackerClient)

			err := mirror.MirrorBacklog(privateProjectId, publicProjectId)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("mirror failed to delete stories:\n %s", "fake client error")))
		})
	})

})
