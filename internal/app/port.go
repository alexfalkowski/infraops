package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Port describes an application port exposed by Kubernetes resources.
type Port struct {
	// Name is the stable Kubernetes service port name.
	Name string
	// Protocol is the application protocol exposed on the service port.
	Protocol string
	// Number is the TCP port number used by the container and service.
	Number int
}

// Ports returns the ports exposed by app.
func Ports(app *App) []Port {
	if app.IsInternal() {
		return []Port{
			{Name: "debug", Number: 6060, Protocol: "http"},
			{Name: "http", Number: 8080, Protocol: "http"},
			{Name: "grpc", Number: 9090, Protocol: "grpc"},
		}
	}

	return []Port{{Name: "http", Number: 8080, Protocol: "http"}}
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
		AppProtocol: pulumi.String(port.Protocol),
		Name:        pulumi.String(port.Name),
		Port:        pulumi.Int(port.Number),
		Protocol:    pulumi.String("TCP"),
		TargetPort:  pulumi.Int(port.Number),
	}
}
