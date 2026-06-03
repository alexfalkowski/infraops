package test

import (
	"errors"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RunWithMocks runs a Pulumi program with the standard project/stack names and mocks.
func RunWithMocks(run pulumi.RunFunc, mocks pulumi.MockResourceMonitor) error {
	return pulumi.RunErr(run, pulumi.WithMocks("project", "stack", mocks))
}

// Stub is a pulumi.Mocks implementation that returns inputs unchanged and never errors.
type Stub struct{}

// Call returns the call arguments as the result without modification.
func (*Stub) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

// NewResource returns the provided inputs as the resource state without modification.
func (*Stub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", args.Inputs, nil
}

// ResourceStub is a pulumi.Mocks implementation that records created resource inputs by type token.
//
// Resource attempts are recorded before any configured failure is returned, so Resources includes
// failed attempts as well as successful ones.
type ResourceStub struct {
	Stub

	failures  map[string][]resourceFailure
	failAll   error
	resources map[string][]resource.PropertyMap
	mu        sync.Mutex
}

type resourceFailure struct {
	err        error
	occurrence int
}

// ErrResource is the default error returned by ResourceStub for targeted resource failures.
var ErrResource = errors.New("bad resource")

// FailResource makes every resource creation for token fail with ErrResource.
func (s *ResourceStub) FailResource(token string) {
	s.FailResourceWith(token, ErrResource)
}

// FailAllResources makes every resource creation fail with ErrResource.
func (s *ResourceStub) FailAllResources() {
	s.FailAllResourcesWith(ErrResource)
}

// FailAllResourcesWith makes every resource creation fail with err.
func (s *ResourceStub) FailAllResourcesWith(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.failAll = err
}

// FailResourceWith makes every resource creation for token fail with err.
func (s *ResourceStub) FailResourceWith(token string, err error) {
	s.failResource(token, 0, err)
}

// FailResourceAt makes the specified occurrence of resource creation for token fail with ErrResource.
//
// Occurrence is one-based.
func (s *ResourceStub) FailResourceAt(token string, occurrence int) {
	s.FailResourceAtWith(token, occurrence, ErrResource)
}

// FailResourceAtWith makes the specified occurrence of resource creation for token fail with err.
//
// Occurrence is one-based.
func (s *ResourceStub) FailResourceAtWith(token string, occurrence int, err error) {
	s.failResource(token, occurrence, err)
}

// NewResource records the resource inputs by Pulumi type token and returns the inputs unchanged.
//
// If a failure is configured for the resource type or occurrence, the attempt is still recorded
// before the error is returned.
func (s *ResourceStub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	s.mu.Lock()

	if s.resources == nil {
		s.resources = make(map[string][]resource.PropertyMap)
	}
	s.resources[args.TypeToken] = append(s.resources[args.TypeToken], args.Inputs)
	occurrence := len(s.resources[args.TypeToken])
	err := s.resourceError(args.TypeToken, occurrence)
	s.mu.Unlock()

	if err != nil {
		return "", args.Inputs, err
	}

	return s.Stub.NewResource(args)
}

// Resources returns the resource inputs recorded for token.
func (s *ResourceStub) Resources(token string) []resource.PropertyMap {
	s.mu.Lock()
	defer s.mu.Unlock()

	resources := s.resources[token]
	return append([]resource.PropertyMap(nil), resources...)
}

func (s *ResourceStub) failResource(token string, occurrence int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.failures == nil {
		s.failures = make(map[string][]resourceFailure)
	}
	s.failures[token] = append(s.failures[token], resourceFailure{err: err, occurrence: occurrence})
}

func (s *ResourceStub) resourceError(token string, occurrence int) error {
	if s.failAll != nil {
		return s.failAll
	}
	for _, failure := range s.failures[token] {
		if failure.occurrence == 0 || failure.occurrence == occurrence {
			return failure.err
		}
	}
	return nil
}

// Property returns the property value for key or fails the test.
func Property(t *testing.T, properties resource.PropertyMap, key string) resource.PropertyValue {
	t.Helper()

	value, ok := properties[resource.PropertyKey(key)]
	if !ok {
		t.Fatalf("missing %s", key)
	}

	return value
}

// StringValues converts Pulumi property values to strings.
func StringValues(values []resource.PropertyValue) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.StringValue())
	}
	return result
}
