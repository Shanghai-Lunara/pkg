package main

import (
	"context"
	"flag"
	"github.com/Shanghai-Lunara/pkg/websocket/echoserver"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	s := echoserver.Init(context.Background())
	<-stopCh
	s.Shutdown()
}
