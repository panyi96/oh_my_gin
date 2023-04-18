package nacosconfig

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

var Config nacosConfig

type nacosConfig struct {
	ProjectName string `yaml:"projectName"`
	Version     string `yaml:"version"`
	Redis       redis  `yaml:"redis"`
	Mysql       mysql  `yaml:"mysql"`
}
type redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}
type mysql struct {
	DSN string `yaml:"dsn"`
}

// Init TODO The nacos config ready get from command argument
func Init(ch chan int) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "ea50ebb8-bbe3-43e7-8d7d-7c23b0f87fa2", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//LogDir:              "./nacos/log",
		//CacheDir:            "./nacos/cache",
		LogLevel:       "debug",
		AppendToStdout: true, //打印到控制台
	}

	//// 创建clientConfig的另一种方式
	//clientConfig := *constant.NewClientConfig(
	//	constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"), //当namespace是public时，此处填空字符串。
	//	constant.WithTimeoutMs(5000),
	//	constant.WithNotLoadCacheAtStart(true),
	//	constant.WithLogDir("/tmp/nacos/log"),
	//	constant.WithCacheDir("/tmp/nacos/cache"),
	//	constant.WithLogLevel("debug"),
	//)

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:   "127.0.0.1",
			Port:     8848,
			GrpcPort: 9848,
		},
	}

	//// 创建serverConfig的另一种方式
	//serverConfigs := []constant.ServerConfig{
	//	*constant.NewServerConfig(
	//		"console1.nacos.io",
	//		80,
	//		constant.WithScheme("http"),
	//		constant.WithContextPath("/nacos"),
	//	),
	//	*constant.NewServerConfig(
	//		"console2.nacos.io",
	//		80,
	//		constant.WithScheme("http"),
	//		constant.WithContextPath("/nacos"),
	//	),
	//}

	// 创建服务发现客户端的另一种方式
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if nil != err {
		panic(err)
	}
	// 创建动态配置客户端的另一种方式
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if nil != err {
		panic(err)
	}

	//获取配置
	content, err := configClient.GetConfig(
		vo.ConfigParam{
			DataId: "oh-my-gin-dev.yaml",
			Group:  "study",
		},
	)

	if nil != err {
		panic(err)
	}

	logger.Infof("获取到的配置\n%s", content)
	err = yaml.Unmarshal([]byte(content), &Config)
	if nil != err {
		logger.Error("===yaml转换结构体错误%s\n===", err.Error())
	}
	ch <- 1
	//阻塞结束

	//监听配置变化
	listenErr := configClient.ListenConfig(vo.ConfigParam{
		DataId: "oh-my-gin-dev.yaml",
		Group:  "study",
		OnChange: func(namespace, group, dataId, data string) {
			logger.Debugf("配置发生变化ListenConfig：group:%s, dataId:%s, data:%s", group, dataId, data)
		},
	})
	if nil != listenErr {
		logger.Error("===监听配置发生错误%s\n===", err.Error())
	}

	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        8849,
		ServiceName: "oh_my_gin.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		//ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName: "study", // 默认值DEFAULT_GROUP
	})
	if !success {
		logger.Errorf("===注册实例发生错误%s\n===", err.Error())
	}
}
