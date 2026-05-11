package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	debugPortName = "debug"
	debugPort     = 6060
	httpPortName  = "http"
	httpPort      = 8080
	grpcPortName  = "grpc"
	grpcPort      = 9090
)

// Port describes an application port exposed by Kubernetes resources.
type Port struct {
	// Name is the stable Kubernetes service port name.
	Name string
	// Number is the TCP port number used by the container and service.
	Number int
}

// Ports returns the ports exposed by app.
func Ports(app *App) []Port {
	if app.IsInternal() {
		return []Port{
			{Name: debugPortName, Number: debugPort},
			{Name: httpPortName, Number: httpPort},
			{Name: grpcPortName, Number: grpcPort},
		}
	}

	return []Port{{Name: httpPortName, Number: httpPort}}
}

func containerPorts(app *App) cv1.ContainerPortArray {
	ports := Ports(app)
	containerPorts := make(cv1.ContainerPortArray, 0, len(ports))
	for _, port := range ports {
		containerPorts = append(containerPorts, cv1.ContainerPortArgs{ContainerPort: pulumi.Int(port.Number)})
	}

	return containerPorts
}

func servicePorts(app *App) cv1.ServicePortArray {
	ports := Ports(app)
	servicePorts := make(cv1.ServicePortArray, 0, len(ports))
	for _, port := range ports {
		servicePorts = append(servicePorts, servicePort(port))
	}

	return servicePorts
}

func servicePort(port Port) cv1.ServicePortArgs {
	return cv1.ServicePortArgs{
		AppProtocol: pulumi.String("TCP"),
		Name:        pulumi.String(port.Name),
		Port:        pulumi.Int(port.Number),
		TargetPort:  pulumi.Int(port.Number),
	}
}
