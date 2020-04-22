package processor

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	mock_rabbitmq "github.com/canyanio/rating-agent-janus/client/rabbitmq/mock"
	dconfig "github.com/canyanio/rating-agent-janus/config"
	"github.com/canyanio/rating-agent-janus/model"
	"github.com/canyanio/rating-agent-janus/state"
)

func TestProcessSIPMessageStart(t *testing.T) {
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "invite.data")
	buffInvite, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "ack.data")
	buffAck, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "bye.data")
	buffBye, _ := ioutil.ReadFile(path)

	// mock rabbitmq client
	mockClient := &mock_rabbitmq.Client{}
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameBeginTransaction,
		mock.MatchedBy(func(req *model.BeginTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:08Z", req.Request.TimestampBegin)

			return true
		}),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.MatchedBy(func(req *model.EndTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:09Z", req.Request.TimestampEnd)

			return true
		}),
	).Return(nil)

	ctx := context.Background()

	stateManager := state.NewMemoryManager()
	processor := NewJanusProcessor(stateManager, mockClient)

	sipMessage := model.SIPMessageFromString(string(buffInvite))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:07Z")
	processor.processSIPMessage(ctx, sipMessage)

	sipMessage = model.SIPMessageFromString(string(buffAck))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:08Z")
	processor.processSIPMessage(ctx, sipMessage)

	sipMessage = model.SIPMessageFromString(string(buffBye))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:09Z")
	processor.processSIPMessage(ctx, sipMessage)

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestProcessSIPMessageStartWithRedis(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "invite.data")
	buffInvite, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "ack.data")
	buffAck, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "bye.data")
	buffBye, _ := ioutil.ReadFile(path)

	// mock rabbitmq client
	mockClient := &mock_rabbitmq.Client{}
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameBeginTransaction,
		mock.AnythingOfType("*model.BeginTransaction"),
	).Return(nil)
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameEndTransaction,
		mock.AnythingOfType("*model.EndTransaction"),
	).Return(nil)

	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	ctx := context.Background()

	stateManager := state.NewRedisManager(redisAddress, redisPassword, redisDb)
	stateManager.Connect(ctx)
	defer stateManager.Close(ctx)

	processor := NewJanusProcessor(stateManager, mockClient)

	sipMessage := model.SIPMessageFromString(string(buffInvite))
	processor.processSIPMessage(ctx, sipMessage)

	sipMessage = model.SIPMessageFromString(string(buffAck))
	processor.processSIPMessage(ctx, sipMessage)

	sipMessage = model.SIPMessageFromString(string(buffBye))
	processor.processSIPMessage(ctx, sipMessage)

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestProcessSIPMessageStartPublishFailure(t *testing.T) {
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "invite.data")
	buffInvite, _ := ioutil.ReadFile(path)

	path = filepath.Join(cwd, "..", "testdata", "ack.data")
	buffAck, _ := ioutil.ReadFile(path)

	// mock rabbitmq client
	mockClient := &mock_rabbitmq.Client{}
	mockClient.On("Publish",
		mock.MatchedBy(func(_ context.Context) bool {
			return true
		}),
		rabbitmq.QueueNameBeginTransaction,
		mock.MatchedBy(func(req *model.BeginTransaction) bool {
			assert.Equal(t, "2020-03-14T08:56:08Z", req.Request.TimestampBegin)

			return true
		}),
	).Return(errors.New("generic error"))

	ctx := context.Background()

	stateManager := state.NewMemoryManager()
	processor := NewJanusProcessor(stateManager, mockClient)

	sipMessage := model.SIPMessageFromString(string(buffInvite))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:07Z")
	processor.processSIPMessage(ctx, sipMessage)

	sipMessage = model.SIPMessageFromString(string(buffAck))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:08Z")
	processor.processSIPMessage(ctx, sipMessage)

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}

func TestProcessSIPMessageStartAckWithoutInvite(t *testing.T) {
	cwd, _ := os.Getwd()

	path := filepath.Join(cwd, "..", "testdata", "ack.data")
	buffAck, _ := ioutil.ReadFile(path)

	// mock rabbitmq client
	mockClient := &mock_rabbitmq.Client{}

	ctx := context.Background()

	stateManager := state.NewMemoryManager()
	processor := NewJanusProcessor(stateManager, mockClient)

	sipMessage := model.SIPMessageFromString(string(buffAck))
	sipMessage.Timestamp, _ = time.Parse("2006-01-02T15:04:05Z", "2020-03-14T08:56:07Z")
	processor.processSIPMessage(ctx, sipMessage)

	// assert expectations (processor)
	mockClient.AssertExpectations(t)
}
