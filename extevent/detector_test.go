package extevent

import (
	"SyncthingHook/stclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

type MockCoreEventDetector struct {
	wg sync.WaitGroup
	mock.Mock
}

func (m *MockCoreEventDetector) subscribeUpstream(st *stclient.Syncthing) upstream {
	args := m.Called(st)
	return args.Get(0).(upstream)
}

func (m *MockCoreEventDetector) unsubscribeUpstream(st *stclient.Syncthing, up upstream) {
	m.Called(st, up)
}

func (m *MockCoreEventDetector) handleUpstream(u upstream, dispatcher downstreamDispatcher) {
	m.Called(u, dispatcher)
	m.wg.Done()
}

func TestSubscribe(t *testing.T) {
	setupMockDetector := func() (d *MockCoreEventDetector) {
		d = &MockCoreEventDetector{}
		stubUpstream := make(upstream)
		d.On("subscribeUpstream", mock.Anything).Return(stubUpstream)
		d.On("handleUpstream", stubUpstream, mock.Anything).Return()
		return
	}

	t.Run("firstSubscription", func(t *testing.T) {
		mockCore := setupMockDetector()
		mockCore.wg.Add(1)
		d := newDetector(nil, mockCore)
		assert.NotNil(t, d.Subscribe())

		mockCore.wg.Wait()
		mockCore.AssertExpectations(t)
		assert.Len(t, d.subscribers, 1)
	})

	t.Run("moreThanOneSubscription", func(t *testing.T) {
		mockCore := setupMockDetector()
		mockCore.wg.Add(1)
		d := newDetector(nil, mockCore)
		assert.NotNil(t, d.Subscribe())
		assert.NotNil(t, d.Subscribe())

		mockCore.wg.Wait()
		mockCore.AssertNumberOfCalls(t, "handleUpstream", 1)
		mockCore.AssertNumberOfCalls(t, "subscribeUpstream", 1)
		assert.Len(t, d.subscribers, 2)
	})

	t.Run("asyncSubscribe", func(t *testing.T) {
		const NUM = 500
		mockCore := setupMockDetector()
		mockCore.wg.Add(1)
		d := newDetector(nil, mockCore)
		var asyncWg sync.WaitGroup
		asyncWg.Add(NUM)
		for i := 0; i < NUM; i++ {
			go func() {
				assert.NotNil(t, d.Subscribe())
				asyncWg.Done()
			}()
		}

		asyncWg.Wait()
		mockCore.wg.Wait()
		mockCore.AssertNumberOfCalls(t, "handleUpstream", 1)
		mockCore.AssertNumberOfCalls(t, "subscribeUpstream", 1)
		assert.Len(t, d.subscribers, NUM)
	})
}

func TestUnsubscribe(t *testing.T) {
	setup := func(subscribeNum int) (mockCore *MockCoreEventDetector, d *detector, outs []<-chan Event) {
		mockCore = &MockCoreEventDetector{}
		outs = make([]<-chan Event, 0, subscribeNum)
		stubUpstream := make(upstream)
		mockCore.On("subscribeUpstream", mock.Anything).Return(stubUpstream)
		mockCore.On("handleUpstream", stubUpstream, mock.Anything).Return()
		mockCore.On("unsubscribeUpstream", mock.Anything, mock.Anything).Return()

		d = newDetector(nil, mockCore)
		mockCore.wg.Add(1)
		for i := 0; i < subscribeNum; i++ {
			outs = append(outs, d.Subscribe())
		}
		mockCore.wg.Wait()
		return
	}

	t.Run("removeNotExistListener", func(t *testing.T) {
		_, d, _ := setup(1)
		d.Unsubscribe(make(subscriber))

		assert.Len(t, d.subscribers, 1)
	})

	t.Run("removeTheOnlyListener", func(t *testing.T) {
		mockCore, d, outs := setup(1)
		d.Unsubscribe(outs[0])

		assert.Len(t, d.subscribers, 0)
		// should unsubscribe the upstream as well
		mockCore.AssertNumberOfCalls(t, "unsubscribeUpstream", 1)
		assert.Nil(t, d.upstream)
	})

	t.Run("removeOneOfTwoListeners", func(t *testing.T) {
		mockCore, d, outs := setup(2)
		d.Unsubscribe(outs[0])

		assert.Len(t, d.subscribers, 1)
		// should not unsubscribe the upstream
		mockCore.AssertNotCalled(t, "unsubscribeUpstream", mock.Anything, mock.Anything)
		assert.NotNil(t, d.upstream)
	})

	t.Run("asyncUnsubscribe", func(t *testing.T) {
		const NUM = 500
		mockCore, d, outs := setup(NUM)
		var asyncWg sync.WaitGroup
		asyncWg.Add(NUM)
		go func() {
			for i := 0; i < NUM; i++ {
				d.Unsubscribe(outs[i])
				asyncWg.Done()
			}
		}()
		asyncWg.Wait()

		assert.Len(t, d.subscribers, 0)
		mockCore.AssertNumberOfCalls(t, "unsubscribeUpstream", 1)
		assert.Nil(t, d.upstream)
	})
}
