package cmd

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/dinhtp/lets-run-woocommerce-customer/customer"
	ppb "github.com/dinhtp/lets-run-pbtype/platform"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run woocommerce customer grpc command",
	Run:   runGrpcCommand,
}

func init() {
	grpcCmd.Flags().StringP("backend", "", "grpc-address", "gRPC address")

	_ = viper.BindPFlag("backend", grpcCmd.Flags().Lookup("backend"))

	serveCmd.AddCommand(grpcCmd)
}

func runGrpcCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// services
	grpcServer := initializeServices(grpc.NewServer())

	// init GRPC backend
	grpcAddr := viper.GetString("backend")
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	// Serve GRPC
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	logrus.WithFields(logrus.Fields{
		"service": "run-woocommerce-customer-service",
		"type":    "grpc",
		"address": grpcAddr,
	}).Info("run woocommerce customer service server started")

	<-c
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logrus.WithFields(logrus.Fields{
		"service": "run-woocommerce-customer-service",
		"type":    "grpc",
		"address": grpcAddr,
	}).Info("run woocommerce customer service gracefully shutdowns")
}


func initializeServices(grpcServer *grpc.Server) *grpc.Server {
	customerService := customer.NewService()
	ppb.RegisterCustomerServiceServer(grpcServer, customerService)

	return grpcServer
}