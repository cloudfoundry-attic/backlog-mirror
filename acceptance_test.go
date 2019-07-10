package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/gexec"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type story struct {
	Id            int
	Name          string
	Current_state string
}

var _ = Describe("Backlog Mirror Application", func() {

	var testStory story

	runBacklogMirror := func() *gexec.Session {
		backlogMirrorCmd := exec.Command(BacklogMirrorExecutable)
		backlogMirrorCmd.Stderr = os.Stderr
		backlogMirrorCmd.Env = []string{
			fmt.Sprintf("TRACKER_API_TOKEN=%s", APIToken),
			fmt.Sprintf("TRACKER_ORIG_BACKLOG=%d", PrivateBacklogId),
			fmt.Sprintf("TRACKER_DEST_BACKLOG=%d", PublicBacklogId),
		}

		session, err := gexec.Start(backlogMirrorCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		return session
	}

	setupNewPublicLabelStory := func() {
		uuid, err := uuid.NewV4()
		Expect(err).ToNot(HaveOccurred())

		testStoryName := "story to be made public " + uuid.String()
		storyRequest := fmt.Sprintf(
			`{
			"name": "%s",
			"current_state": "unstarted",
			"labels": ["public"]
		}`, testStoryName)

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
		testStory = newStory
		time.Sleep(1 * time.Second)
	}

	It("Clones a story from private to public backlog", func() {

		setupNewPublicLabelStory()
		session := runBacklogMirror()
		Eventually(session, 12 * time.Second).Should(gexec.Exit(0))

		publicBacklogStoriesEndpoint := fmt.Sprintf(StoriesEndpoint, PublicBacklogId)
		storiesResponse, err := http.Get(publicBacklogStoriesEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(storiesResponse).NotTo(BeNil())
		storiesResponseBody, err := ioutil.ReadAll(storiesResponse.Body)

		var publicBacklogStories []story
		err = json.Unmarshal(storiesResponseBody, &publicBacklogStories)
		Expect(err).ToNot(HaveOccurred())

		var storyNames []string
		for _, story := range publicBacklogStories {
			storyNames = append(storyNames, story.Name)
		}
		Expect(storyNames).To(ContainElement(testStory.Name))
	})

	AfterEach(func() {
		if testStory.Id != 0 {
			trackerStoriesEndpoint := "https://www.pivotaltracker.com/services/v5/projects/2345567/stories/" + strconv.Itoa(testStory.Id)

			req, _ := http.NewRequest(http.MethodDelete, trackerStoriesEndpoint, nil)
			req.Header.Add("X-TrackerToken", APIToken)

			response, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(BeNumerically(">=", 200))
			Expect(response.StatusCode).To(BeNumerically("<", 300))
		}
	})
})
