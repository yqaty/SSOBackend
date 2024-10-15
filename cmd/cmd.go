package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/server"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/spf13/cobra"
	"github.com/xylonx/zapx"
	"github.com/xylonx/zapx/decoder"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "unique-sso",
	Short: "unique studio sso service",
	PreRunE: func(c *cobra.Command, args []string) (err error) {
		err = config.Setup(cfgFile)
		if err != nil {
			return err
		}

		err = core.Setup()
		if err != nil {
			return err
		}

		zapx.UseCtxDecoder(decoder.OpentelemetaryDecoder)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		shutdown, err := tracer.SetupTracing(
			config.Config.Apm.Name,
			config.Config.Application.Mode,
			config.Config.Apm.ReportBackend,
		)
		if err != nil {
			zapx.Warn("setup tracing report backend failed", zap.Error(err))
		}

		httpAddr := config.Config.Application.HttpHost + ":" + strconv.FormatInt(config.Config.Application.HttpPort, 10)
		httpServer := server.InitHttpServer(&server.HttpOption{
			Addr:         httpAddr,
			ReadTimeout:  time.Duration(config.Config.Application.HttpReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.Config.Application.HttpWriteTimeout) * time.Second,
			AllowOrigins: config.Config.Application.HttpAllowOrigin,
			Mode:         config.Config.Application.Mode,
		})

		grpcAddr := fmt.Sprintf("%s:%d", config.Config.Application.GrpcHost, config.Config.Application.GrpcPort)
		grpcLis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			zapx.Error("listen grpc addr failed", zap.String("addr", grpcAddr))
			return err
		}
		grpcServer := server.InitGRPCServer()

		httpErrCh := make(chan error)
		grpcErrCh := make(chan error)

		go func() {
			zapx.Info("start http server", zap.String("host", httpAddr))
			if err := httpServer.ListenAndServe(); err != nil {
				zapx.Error("http run error", zap.Error(err))
				httpErrCh <- err
			}
		}()

		go func() {
			zapx.Info("start grpc server", zap.String("addr", grpcAddr))
			if err := grpcServer.Serve(grpcLis); err != nil {
				zapx.Error("serve grpc listen failed", zap.Error(err))
				grpcErrCh <- err
			}
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		select {
		case <-sig:
		case <-httpErrCh:
		case <-grpcErrCh:
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(3)
		go func() {
			shutdown(ctx)
			wg.Done()
		}()
		go func() {
			if err := httpServer.Shutdown(ctx); err != nil {
				zapx.Error("shutdown http server failed", zap.Error(err))
			}
			wg.Done()
		}()
		go func() {
			grpcServer.GracefulStop()
			wg.Done()
		}()
		wg.Wait()

		return nil
	},
}

var cfgFile string

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "./config.local.yaml", "path to config file")
}

func Execute() error {
	return rootCmd.Execute()
}
