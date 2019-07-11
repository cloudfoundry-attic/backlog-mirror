package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type story struct {
	Id            int
	Name          string
	Current_state string
}

var _ = Describe("Backlog Mirror Application", func() {

	var createdStoryIDs []int

	runBacklogMirror := func(environment []string) *gexec.Session {
		backlogMirrorCmd := exec.Command(BacklogMirrorExecutable)
		backlogMirrorCmd.Stderr = os.Stderr
		backlogMirrorCmd.Env = environment

		session, err := gexec.Start(backlogMirrorCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		return session
	}

	setupNewPrivateTrackerStoryWithLabel := func(label string) string {

		uuid, err := uuid.NewV4()
		Expect(err).ToNot(HaveOccurred())

		testStoryName := "brand new story " + uuid.String()
		storyRequest := fmt.Sprintf(
			`{
			"name": "%s",
			"current_state": "unstarted",
			"labels": ["%s"]
		}`, testStoryName, label)

		privateBacklogStoriesEndpoint := fmt.Sprintf(StoriesEndpoint, PrivateBacklogId)

		req, err := http.NewRequest("POST", privateBacklogStoriesEndpoint, strings.NewReader(storyRequest))
		Expect(err).ToNot(HaveOccurred())
		req.Header.Add("X-TrackerToken", APIToken)
		req.Header.Add("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		newStoryBody, err := ioutil.ReadAll(resp.Body)

		var newStory story
		err = json.Unmarshal(newStoryBody, &newStory)
		Expect(err).ToNot(HaveOccurred())
		Expect(newStory.Id).NotTo(Equal(0))
		Expect(newStory.Name).To(Equal(testStoryName))
		createdStoryIDs = append(createdStoryIDs, newStory.Id)

		return testStoryName
	}

	fetchStoriesFromBacklog := func(projectID int) []story {

		backlogStoriesEndpoint := fmt.Sprintf(StoriesEndpoint, projectID)
		storiesResponse, err := http.Get(backlogStoriesEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(storiesResponse).NotTo(BeNil())
		storiesResponseBody, err := ioutil.ReadAll(storiesResponse.Body)
		Expect(err).ToNot(HaveOccurred())

		var backlogStories []story
		err = json.Unmarshal(storiesResponseBody, &backlogStories)
		Expect(err).ToNot(HaveOccurred())

		return backlogStories
	}

	Context("with valid api token and project ids", func() {

		var (
			validEnv []string
			newPublicStories []string
			nonPublicStories []string

		)
		BeforeEach(func() {
			createdStoryIDs = []int{}
			validEnv = []string{
				fmt.Sprintf("TRACKER_API_TOKEN=%s", APIToken),
				fmt.Sprintf("TRACKER_ORIG_BACKLOG=%d", PrivateBacklogId),
				fmt.Sprintf("TRACKER_DEST_BACKLOG=%d", PublicBacklogId),
			}

			newPublicStories = []string{}
			for i := 0; i < 5; i++ {
				newPublicStories = append(newPublicStories, setupNewPrivateTrackerStoryWithLabel("public"))
			}

			nonPublicStories = []string{}
			for i := 0; i < 5; i++ {
				nonPublicStories = append(nonPublicStories, setupNewPrivateTrackerStoryWithLabel("notably-not-public"))
			}
		})

		It("Clones all and only public labeled stories from private backlog to public backlog", func() {

			By("running the backlog mirror")
			session := runBacklogMirror(validEnv)
			Eventually(session, 60 * time.Second).Should(gexec.Exit(0))

			By("capturing stories on the public backlog via http request")
			publicBacklogStories := fetchStoriesFromBacklog(PublicBacklogId)

			var storyNames []string
			for _, story := range publicBacklogStories {
				storyNames = append(storyNames, story.Name)
			}

			By("and verifying the (non-)existence of the (non-)public labeled stories on the public backlog")
			for _, publicStory := range newPublicStories {
				Expect(storyNames).To(ContainElement(publicStory))
			}

			for _, nonPublicStory := range nonPublicStories {
				Expect(storyNames).NotTo(ContainElement(nonPublicStory))
			}
		})

		It("Clears the public backlog prior to repopulating it", func() {
			By("running the backlog mirror")
			session := runBacklogMirror(validEnv)
			Eventually(session, 60 * time.Second).Should(gexec.Exit(0))

			By("capturing stories on the public backlog via http request")
			publicBacklogStories := fetchStoriesFromBacklog(PublicBacklogId)

			var storyIDsAfterFirstMirror []int
			for _, story := range publicBacklogStories {
				storyIDsAfterFirstMirror = append(storyIDsAfterFirstMirror, story.Id)
			}

			By("running the backlog mirror again")
			session = runBacklogMirror(validEnv)
			Eventually(session, 60 * time.Second).Should(gexec.Exit(0))

			By("capturing stories on the public backlog via http request again")
			publicBacklogStories = fetchStoriesFromBacklog(PublicBacklogId)

			By("and verifying no story ID from the first request is contained in the second")
			// this verification method can possibly flake as tracker can reuse story IDs once they're deleted
			var storyIDsAfterSecondMirror []int
			for _, story := range publicBacklogStories {
				storyIDsAfterSecondMirror = append(storyIDsAfterSecondMirror, story.Id)
			}

			for _, id := range storyIDsAfterFirstMirror {
				Expect(storyIDsAfterSecondMirror).ShouldNot(ContainElement(id))
			}

		})


		deleteStoriesFromBacklog := func(storyIDs []int, projectID int) {

			for _, id := range storyIDs {

				storyEndpoint := fmt.Sprintf(StoriesEndpoint, projectID) + fmt.Sprintf("/%d", id)

				req, err := http.NewRequest(http.MethodDelete, storyEndpoint, nil)
				Expect(err).ToNot(HaveOccurred())
				req.Header.Add("X-TrackerToken", APIToken)

				response, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(BeNumerically(">=", 200))
				Expect(response.StatusCode).To(BeNumerically("<", 300))
			}
		}

		flushBacklog := func(projectID int) {

			stories := fetchStoriesFromBacklog(projectID)
			var storyIDs []int
			for _, story := range stories {
				storyIDs = append(storyIDs, story.Id)
			}

			deleteStoriesFromBacklog(storyIDs, projectID)
		}

		AfterEach(func() {

			// our test private backlog has edge case stories, so we don't flush the entire private backlog
			//		we should make these cases explicit in our setup and flush the entire private backlog here
			deleteStoriesFromBacklog(createdStoryIDs, PrivateBacklogId)
			flushBacklog(PublicBacklogId)

		})
	})

	Context("with unset environment variables", func() {
		It("fails with a useful error message when token isn't set", func() {
			session := runBacklogMirror([]string{
				fmt.Sprintf("TRACKER_API_TOKEN="),
				fmt.Sprintf("TRACKER_ORIG_BACKLOG=%d", PrivateBacklogId),
				fmt.Sprintf("TRACKER_DEST_BACKLOG=%d", PublicBacklogId),
			})
			Eventually(session, 12 * time.Second).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("TRACKER_API_TOKEN must be set"))
		})

		It("fails with a useful error message when destination backlog isn't set", func() {
			session := runBacklogMirror([]string{
				fmt.Sprintf("TRACKER_API_TOKEN=%s", APIToken),
				fmt.Sprintf("TRACKER_ORIG_BACKLOG=%d", PrivateBacklogId),
				fmt.Sprintf("TRACKER_DEST_BACKLOG="),
			})
			Eventually(session, 12 * time.Second).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("TRACKER_ORIG_BACKLOG and TRACKER_DEST_BACKLOG must be set"))
		})

		It("fails with a useful error message original backlog isn't set", func() {
			session := runBacklogMirror([]string{
				fmt.Sprintf("TRACKER_API_TOKEN=%s", APIToken),
				fmt.Sprintf("TRACKER_ORIG_BACKLOG="),
				fmt.Sprintf("TRACKER_DEST_BACKLOG=%d", PublicBacklogId),
			})
			Eventually(session, 12 * time.Second).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("TRACKER_ORIG_BACKLOG and TRACKER_DEST_BACKLOG must be set"))
		})
	})
})
