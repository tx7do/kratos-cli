package render

import (
	"fmt"
	"strings"
)

type InitWireTemplateData struct {
	Package      string
	Postfix      string
	ServiceNames []string
}

type WireTemplateData struct {
	Project string
	Service string
}

type MainTemplateData struct {
	Project string
	Service string
	Servers []string
}

func (d MainTemplateData) ServerImportPath() string {
	var str string

	for _, server := range d.Servers {
		switch strings.TrimSpace(strings.ToLower(server)) {
		case "grpc":
			str += fmt.Sprintf("\t\"github.com/go-kratos/kratos/v2/transport/grpc\"\n")
		case "rest":
			str += fmt.Sprintf("\t\"github.com/go-kratos/kratos/v2/transport/http\"\n")

		case "activemq":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/activemq\"\n")
		case "asynq":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/asynq\"\n")
		case "fasthttp":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/fasthttp\"\n")
		case "gin":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/gin\"\n")
		case "gozero":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/gozero\"\n")
		case "graphql":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/graphql\"\n")
		case "hertz":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/hertz\"\n")
		case "iris":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/iris\"\n")
		case "kafka":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/kafka\"\n")
		case "machinery":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/machinery\"\n")
		case "mqtt":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/mqtt\"\n")
		case "nats":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/nats\"\n")
		case "nsq":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/nsq\"\n")
		case "pulsar":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/pulsar\"\n")
		case "rabbitmq":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/rabbitmq\"\n")
		case "redis":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/redis\"\n")
		case "rocketmq":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/rocketmq\"\n")
		case "signalr":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/signalr\"\n")
		case "socketio":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/socketio\"\n")
		case "sse":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/sse\"\n")
		case "tcp":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/tcp\"\n")
		case "thrift":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/thrift\"\n")
		case "trpc":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/trpc\"\n")
		case "websocket":
			str += fmt.Sprintf("\t\"github.com/tx7do/kratos-transport/transport/websocket\"\n")
		}
	}

	return str
}

// ServerFormalParameters 形参
func (d MainTemplateData) ServerFormalParameters() string {
	var str string

	for _, server := range d.Servers {
		switch strings.TrimSpace(strings.ToLower(server)) {
		case "grpc":
			str += fmt.Sprintf("\tgs *grpc.Server,\n")
		case "rest":
			str += fmt.Sprintf("\ths *http.Server,\n")

		case "activemq":
			str += fmt.Sprintf("\tts *activemq.Server,\n")
		case "asynq":
			str += fmt.Sprintf("\tas *asynq.Server,\n")
		case "fasthttp":
			str += fmt.Sprintf("\tfs *fasthttp.Server,\n")
		case "gin":
			str += fmt.Sprintf("\tis *gin.Server,\n")
		case "gozero":
			str += fmt.Sprintf("\tos *gozero.Server,\n")
		case "graphql":
			str += fmt.Sprintf("\tqs *graphql.Server,\n")
		case "hertz":
			str += fmt.Sprintf("\tzs *hertz.Server,\n")
		case "iris":
			str += fmt.Sprintf("\trs *iris.Server,\n")
		case "kafka":
			str += fmt.Sprintf("\tks *kafka.Server,\n")
		case "machinery":
			str += fmt.Sprintf("\tys *machinery.Server,\n")
		case "mqtt":
			str += fmt.Sprintf("\tms *mqtt.Server,\n")
		case "nats":
			str += fmt.Sprintf("\tnats *nats.Server,\n")
		case "nsq":
			str += fmt.Sprintf("\tnsqs *nsq.Server,\n")
		case "pulsar":
			str += fmt.Sprintf("\tpulsars *pulsar.Server,\n")
		case "rabbitmq":
			str += fmt.Sprintf("\trabbitmqs *rabbitmq.Server,\n")
		case "redis":
			str += fmt.Sprintf("\trediss *redis.Server,\n")
		case "rocketmq":
			str += fmt.Sprintf("\trocketmqs *rocketmq.Server,\n")
		case "signalr":
			str += fmt.Sprintf("\tsignalrs *signalr.Server,\n")
		case "socketio":
			str += fmt.Sprintf("\tsocketios *socketio.Server,\n")
		case "sse":
			str += fmt.Sprintf("\tss *sse.Server,\n")
		case "tcp":
			str += fmt.Sprintf("\ttcps *tcp.Server,\n")
		case "thrift":
			str += fmt.Sprintf("\tthrifts *thrift.Server,\n")
		case "trpc":
			str += fmt.Sprintf("\ttrpcs *trpc.Server,\n")
		case "websocket":
			str += fmt.Sprintf("\tws *websocket.Server,\n")
		}
	}

	return str
}

// ServerTransferParameters 传参
func (d MainTemplateData) ServerTransferParameters() string {
	var str string

	for _, server := range d.Servers {
		switch strings.TrimSpace(strings.ToLower(server)) {
		case "grpc":
			str += fmt.Sprintf("\tgs,\n")
		case "rest":
			str += fmt.Sprintf("\ths,\n")

		case "activemq":
			str += fmt.Sprintf("\tts,\n")
		case "asynq":
			str += fmt.Sprintf("\tas,\n")
		case "fasthttp":
			str += fmt.Sprintf("\tfs,\n")
		case "gin":
			str += fmt.Sprintf("\tis,\n")
		case "gozero":
			str += fmt.Sprintf("\tos,\n")
		case "graphql":
			str += fmt.Sprintf("\tqs,\n")
		case "hertz":
			str += fmt.Sprintf("\tzs,\n")
		case "iris":
			str += fmt.Sprintf("\trs,\n")
		case "kafka":
			str += fmt.Sprintf("\tks,\n")
		case "machinery":
			str += fmt.Sprintf("\tys,\n")
		case "mqtt":
			str += fmt.Sprintf("\tms,\n")
		case "nats":
			str += fmt.Sprintf("\tnats,\n")
		case "nsq":
			str += fmt.Sprintf("\tnsqs,\n")
		case "pulsar":
			str += fmt.Sprintf("\tpulsars,\n")
		case "rabbitmq":
			str += fmt.Sprintf("\trabbitmqs,\n")
		case "redis":
			str += fmt.Sprintf("\trediss,\n")
		case "rocketmq":
			str += fmt.Sprintf("\trocketmqs,\n")
		case "signalr":
			str += fmt.Sprintf("\tsignalrs,\n")
		case "socketio":
			str += fmt.Sprintf("\tsocketios,\n")
		case "sse":
			str += fmt.Sprintf("\tss,\n")
		case "tcp":
			str += fmt.Sprintf("\ttcps,\n")
		case "thrift":
			str += fmt.Sprintf("\tthrifts,\n")
		case "trpc":
			str += fmt.Sprintf("\ttrpcs,\n")
		case "websocket":
			str += fmt.Sprintf("\tws,\n")
		}
	}

	return str
}

type ProtoItem struct {
	Module string
	Name   string
}

type ServerTemplateData struct {
	Module   string
	Project  string
	Service  string
	Type     string
	Services map[string]string
}

func (d ServerTemplateData) Modules() []string {
	modulesMap := make(map[string]struct{})
	for _, module := range d.Services {
		modulesMap[module] = struct{}{}
	}

	var modules []string
	for module := range modulesMap {
		modules = append(modules, module)
	}
	return modules
}
