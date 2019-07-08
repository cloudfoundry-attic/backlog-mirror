package main_test

import (
	"encoding/json"
	"fmt"
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

const (
	privateBacklogId = 2345567
	publicBacklogId  = 2345570
)

type story struct {
	Id            int
	Name          string
	Current_state string
}

var _ = Describe("Backlog Mirror Application", func() {

	var apiToken string
	var testStory story

	runBacklogMirror := func() {
		backlogMirrorCmd := exec.Command("./backlog-mirror")

		backlogMirrorCmd.Env = []string{
			fmt.Sprintf("TRACKER_API_TOKEN=%s", apiToken),
			fmt.Sprintf("TRACKER_ORIG_BACKLOG=%d", privateBacklogId),
			fmt.Sprintf("TRACKER_DEST_BACKLOG=%d", publicBacklogId),
		}

		err := backlogMirrorCmd.Run()
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(1 * time.Second)
	}

	setupNewPublicLabelStory := func() {
		uuid, _ := uuid.NewV4()
		testStoryName := "story to be made public " + uuid.String()

		requestJson := fmt.Sprintf(
			`{
			"name": "%s",
			"current_state": "unstarted",
			"labels": ["public"]
		}`, testStoryName)

		privateBacklogStoriesEndpoint := fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories", privateBacklogId)
		req, _ := http.NewRequest("POST", privateBacklogStoriesEndpoint, strings.NewReader(requestJson))
		req.Header.Add("X-TrackerToken", apiToken)
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
		err := exec.Command("go", "build").Run()
		Expect(err).NotTo(HaveOccurred())

		apiToken = os.Getenv("TRACKER_API_TOKEN")
		Expect(apiToken).NotTo(BeEmpty())

		setupNewPublicLabelStory()

		runBacklogMirror()

		publicBacklogStoriesEndpoint := fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories", publicBacklogId)
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

	AfterSuite(func() {
		if testStory.Id != 0 {
			trackerStoriesEndpoint := "https://www.pivotaltracker.com/services/v5/projects/2345567/stories/" + strconv.Itoa(testStory.Id)

			req, _ := http.NewRequest(http.MethodDelete, trackerStoriesEndpoint, nil)
			req.Header.Add("X-TrackerToken", apiToken)

			response, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(BeNumerically(">=", 200))
			Expect(response.StatusCode).To(BeNumerically("<", 300))
		}
	})
})
