// Code generated by counterfeiter. DO NOT EDIT.
package mirrorfakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/backlog-mirror/mirror"
	"gopkg.in/salsita/go-pivotaltracker.v2/v5/pivotal"
)

type FakeTrackerClient struct {
	AddStoryToProjectStub        func(int, *pivotal.StoryRequest) error
	addStoryToProjectMutex       sync.RWMutex
	addStoryToProjectArgsForCall []struct {
		arg1 int
		arg2 *pivotal.StoryRequest
	}
	addStoryToProjectReturns struct {
		result1 error
	}
	addStoryToProjectReturnsOnCall map[int]struct {
		result1 error
	}
	GetFilteredStoriesStub        func(int, string) ([]*pivotal.Story, error)
	getFilteredStoriesMutex       sync.RWMutex
	getFilteredStoriesArgsForCall []struct {
		arg1 int
		arg2 string
	}
	getFilteredStoriesReturns struct {
		result1 []*pivotal.Story
		result2 error
	}
	getFilteredStoriesReturnsOnCall map[int]struct {
		result1 []*pivotal.Story
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTrackerClient) AddStoryToProject(arg1 int, arg2 *pivotal.StoryRequest) error {
	fake.addStoryToProjectMutex.Lock()
	ret, specificReturn := fake.addStoryToProjectReturnsOnCall[len(fake.addStoryToProjectArgsForCall)]
	fake.addStoryToProjectArgsForCall = append(fake.addStoryToProjectArgsForCall, struct {
		arg1 int
		arg2 *pivotal.StoryRequest
	}{arg1, arg2})
	fake.recordInvocation("AddStoryToProject", []interface{}{arg1, arg2})
	fake.addStoryToProjectMutex.Unlock()
	if fake.AddStoryToProjectStub != nil {
		return fake.AddStoryToProjectStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addStoryToProjectReturns
	return fakeReturns.result1
}

func (fake *FakeTrackerClient) AddStoryToProjectCallCount() int {
	fake.addStoryToProjectMutex.RLock()
	defer fake.addStoryToProjectMutex.RUnlock()
	return len(fake.addStoryToProjectArgsForCall)
}

func (fake *FakeTrackerClient) AddStoryToProjectCalls(stub func(int, *pivotal.StoryRequest) error) {
	fake.addStoryToProjectMutex.Lock()
	defer fake.addStoryToProjectMutex.Unlock()
	fake.AddStoryToProjectStub = stub
}

func (fake *FakeTrackerClient) AddStoryToProjectArgsForCall(i int) (int, *pivotal.StoryRequest) {
	fake.addStoryToProjectMutex.RLock()
	defer fake.addStoryToProjectMutex.RUnlock()
	argsForCall := fake.addStoryToProjectArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTrackerClient) AddStoryToProjectReturns(result1 error) {
	fake.addStoryToProjectMutex.Lock()
	defer fake.addStoryToProjectMutex.Unlock()
	fake.AddStoryToProjectStub = nil
	fake.addStoryToProjectReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTrackerClient) AddStoryToProjectReturnsOnCall(i int, result1 error) {
	fake.addStoryToProjectMutex.Lock()
	defer fake.addStoryToProjectMutex.Unlock()
	fake.AddStoryToProjectStub = nil
	if fake.addStoryToProjectReturnsOnCall == nil {
		fake.addStoryToProjectReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addStoryToProjectReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTrackerClient) GetFilteredStories(arg1 int, arg2 string) ([]*pivotal.Story, error) {
	fake.getFilteredStoriesMutex.Lock()
	ret, specificReturn := fake.getFilteredStoriesReturnsOnCall[len(fake.getFilteredStoriesArgsForCall)]
	fake.getFilteredStoriesArgsForCall = append(fake.getFilteredStoriesArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetFilteredStories", []interface{}{arg1, arg2})
	fake.getFilteredStoriesMutex.Unlock()
	if fake.GetFilteredStoriesStub != nil {
		return fake.GetFilteredStoriesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getFilteredStoriesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTrackerClient) GetFilteredStoriesCallCount() int {
	fake.getFilteredStoriesMutex.RLock()
	defer fake.getFilteredStoriesMutex.RUnlock()
	return len(fake.getFilteredStoriesArgsForCall)
}

func (fake *FakeTrackerClient) GetFilteredStoriesCalls(stub func(int, string) ([]*pivotal.Story, error)) {
	fake.getFilteredStoriesMutex.Lock()
	defer fake.getFilteredStoriesMutex.Unlock()
	fake.GetFilteredStoriesStub = stub
}

func (fake *FakeTrackerClient) GetFilteredStoriesArgsForCall(i int) (int, string) {
	fake.getFilteredStoriesMutex.RLock()
	defer fake.getFilteredStoriesMutex.RUnlock()
	argsForCall := fake.getFilteredStoriesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTrackerClient) GetFilteredStoriesReturns(result1 []*pivotal.Story, result2 error) {
	fake.getFilteredStoriesMutex.Lock()
	defer fake.getFilteredStoriesMutex.Unlock()
	fake.GetFilteredStoriesStub = nil
	fake.getFilteredStoriesReturns = struct {
		result1 []*pivotal.Story
		result2 error
	}{result1, result2}
}

func (fake *FakeTrackerClient) GetFilteredStoriesReturnsOnCall(i int, result1 []*pivotal.Story, result2 error) {
	fake.getFilteredStoriesMutex.Lock()
	defer fake.getFilteredStoriesMutex.Unlock()
	fake.GetFilteredStoriesStub = nil
	if fake.getFilteredStoriesReturnsOnCall == nil {
		fake.getFilteredStoriesReturnsOnCall = make(map[int]struct {
			result1 []*pivotal.Story
			result2 error
		})
	}
	fake.getFilteredStoriesReturnsOnCall[i] = struct {
		result1 []*pivotal.Story
		result2 error
	}{result1, result2}
}

func (fake *FakeTrackerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addStoryToProjectMutex.RLock()
	defer fake.addStoryToProjectMutex.RUnlock()
	fake.getFilteredStoriesMutex.RLock()
	defer fake.getFilteredStoriesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTrackerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ mirror.TrackerClient = new(FakeTrackerClient)
