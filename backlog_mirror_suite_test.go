package main_test

import (
	"os"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestBacklogMirror(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BacklogMirror Suite")
}

var (
	APIToken string
	BacklogMirrorExecutable string
	PrivateBacklogId int
	PublicBacklogId int
	StoriesEndpoint string
	err error
)

var _ = BeforeSuite(func() {
	APIToken = os.Getenv("TRACKER_API_TOKEN")
	Expect(APIToken).To(Not(BeEmpty()))
	PrivateBacklogId, err = strconv.Atoi(os.Getenv("TEST_TRACKER_ORIG_BACKLOG"))
	Expect(err).ToNot(HaveOccurred())
	PublicBacklogId, err = strconv.Atoi(os.Getenv("TEST_TRACKER_DEST_BACKLOG"))
	Expect(err).ToNot(HaveOccurred())
	StoriesEndpoint = "https://www.pivotaltracker.com/services/v5/projects/%d/stories"


	BacklogMirrorExecutable, err = gexec.Build("github.com/cloudfoundry-incubator/backlog-mirror")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	gexec.KillAndWait()
})