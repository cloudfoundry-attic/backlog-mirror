package mirror_test

import (
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"github.com/cloudfoundry-incubator/backlog-mirror/mirror/mirrorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backlog Mirror", func() {
	It("Gets stories for a project", func() {
		fakeStoryApi := &mirrorfakes.FakeStoryApi{}
		m := mirror.NewMirror(fakeStoryApi)

		m.MirrorBacklog()

		Expect(fakeStoryApi.GetAllStoriesCallCount()).To(Equal(1))
	})
})