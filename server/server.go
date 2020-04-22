package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/unix"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"

	"github.com/canyanio/rating-agent-janus/api"
	"github.com/canyanio/rating-agent-janus/client/rabbitmq"
	dconfig "github.com/canyanio/rating-agent-janus/config"
	"github.com/canyanio/rating-agent-janus/state"
)

// InitAndRun initializes the server and runs it
func InitAndRun(conf config.Reader) error {
	ctx := context.Background()
	l := log.FromContext(ctx)

	messagebusURI := config.Config.GetString(dconfig.SettingMessageBusURI)
	stateManagerType := config.Config.GetString(dconfig.SettingStateManager)
	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	var stateManager state.ManagerInterface
	if stateManagerType == dconfig.StateManagerRedis {
		stateManager = state.NewRedisManager(redisAddress, redisPassword, redisDb)
	} else {
		stateManager = state.NewMemoryManager()
	}
	if err := stateManager.Connect(ctx); err != nil {
		l.Error(err)
		return err
	}
	defer stateManager.Close(ctx)

	client := rabbitmq.NewClient(messagebusURI)
	if err := client.Connect(ctx); err != nil {
		l.Error(err)
		return err
	}
	defer client.Close(ctx)

	var router = api.NewRouter(stateManager, client)

	var listen = conf.GetString(dconfig.SettingListen)
	l.Infof("listening: %s", listen)
	srv := &http.Server{
		Addr:    listen,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)
	<-quit

	l.Info("Shutdown Server ...")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxWithTimeout); err != nil {
		l.Fatal("Server Shutdown: ", err)
	}

	return nil
}
