package main

import (
	"flag"
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	nacosHost      = flag.String("nacos.host", "nacos.test.infra.ww5sawfyut0k.bitsvc.io", "nacos host")
	nacosPort      = flag.Uint64("nacos.port", 8848, "nacos port")
	nacosNamespace = flag.String("ns", "none", "nacos namespace")
	nacosUser      = flag.String("user", "bybit-nacos", "nacos user")
	nacosPassword  = flag.String("password", "bybit-nacos", "nacos password")
	nacosGroup     = flag.String("group", "DEFAULT_GROUP", "nacos group")
	serviceName    = flag.String("service.name", "none", " service name")
	serviceIP      = flag.String("service.ip", "none", "service IP")
	servicePort    = flag.Uint64("service.port", 9090, "service port")
)

func init() {
	flag.Parse()
	if *nacosNamespace == "none" {
		log.Fatalln("ns not set")
	}
	if *serviceName == "none" {
		log.Fatalln("service.name not set")
	}
	if *serviceIP == "none" {
		log.Fatalln("service.ip not set")
	}
}

func main() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			*nacosHost,
			*nacosPort,
			constant.WithContextPath("/nacos"),
		),
	}

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(*nacosNamespace),
		constant.WithUsername(*nacosUser),
		constant.WithPassword(*nacosPassword),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUpdateCacheWhenEmpty(true),
		constant.WithLogDir("tmp/log"),
		constant.WithCacheDir("tmp/cache"),
		constant.WithLogLevel("info"),
	)

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	dreparam := vo.DeregisterInstanceParam{
		Ip:          *serviceIP,
		Port:        *servicePort,
		ServiceName: *serviceName,
		GroupName:   *nacosGroup,
		Ephemeral:   true,
	}

	ok, err := namingClient.DeregisterInstance(dreparam)
	if err != nil {
		panic(err)
	}
	log.Println("DeregisterInstance ok", ok, dreparam)
}
