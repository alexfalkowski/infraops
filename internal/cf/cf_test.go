package cf_test

import (
	"os"
	"path/filepath"
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/cf"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

const (
	r2BucketResourceType       = "cloudflare:index/r2Bucket:R2Bucket"
	r2LifecycleResourceType    = "cloudflare:index/r2BucketLifecycle:R2BucketLifecycle"
	r2CustomDomainResourceType = "cloudflare:index/r2CustomDomain:R2CustomDomain"
	recordResourceType         = "cloudflare:index/record:Record"
	zoneSettingResourceType    = "cloudflare:index/zoneSetting:ZoneSetting"
)

func TestCreateBalancerZone(t *testing.T) {
	stub := &test.ResourceStub{}
	require.NoError(t, createBalancerZone(t, stub))
	requireRecord(t, stub.Resources(recordResourceType), "A", "test.test.com", "127.0.0.1")
	requireRecord(t, stub.Resources(recordResourceType), "AAAA", "test.test.com", "::1")
	requireZoneSetting(t, stub.Resources(zoneSettingResourceType), "ssl", "full")

	stub = &test.ResourceStub{}
	stub.FailResource(zoneSettingResourceType)
	require.Error(t, createBalancerZone(t, stub))
}

func TestCreateBalancerZoneWithoutIPV6(t *testing.T) {
	stub := &test.ResourceStub{}
	err := test.RunWithMocks(func(ctx *pulumi.Context) error {
		zone := &cf.BalancerZone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			IPV4:        "127.0.0.1",
		}

		return cf.CreateBalancerZone(ctx, zone)
	}, stub)

	require.NoError(t, err)
	requireRecord(t, stub.Resources(recordResourceType), "A", "test.test.com", "127.0.0.1")
	require.Len(t, stub.Resources(recordResourceType), 1)
}

func TestCreateBalancerZoneReturnsResourceErrors(t *testing.T) {
	t.Run("setting", func(t *testing.T) {
		stub := &test.ResourceStub{}
		stub.FailResource(zoneSettingResourceType)

		require.Error(t, createBalancerZone(t, stub))
		require.NotEmpty(t, stub.Resources(zoneSettingResourceType))
	})

	t.Run("a record", func(t *testing.T) {
		stub := &test.ResourceStub{}
		stub.FailResourceAt(recordResourceType, 1)

		require.Error(t, createBalancerZone(t, stub))
		require.NotEmpty(t, stub.Resources(recordResourceType))
	})

	t.Run("aaaa record", func(t *testing.T) {
		stub := &test.ResourceStub{}
		stub.FailResourceAt(recordResourceType, 2)

		require.Error(t, createBalancerZone(t, stub))
		requireRecord(t, stub.Resources(recordResourceType), "A", "test.test.com", "127.0.0.1")
		require.Len(t, stub.Resources(recordResourceType), 2)
	})
}

func TestCreatePageZone(t *testing.T) {
	stub := &test.ResourceStub{}
	require.NoError(t, createPageZone(t, stub))
	requireRecord(t, stub.Resources(recordResourceType), "CNAME", "www.test.com", "test.github.io")
	requireZoneSetting(t, stub.Resources(zoneSettingResourceType), "ssl", "strict")

	stub = &test.ResourceStub{}
	stub.FailResource(zoneSettingResourceType)
	require.Error(t, createPageZone(t, stub))
}

func TestCreatePageZoneReturnsRecordError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(recordResourceType)

	require.Error(t, createPageZone(t, stub))
	requireZoneSetting(t, stub.Resources(zoneSettingResourceType), "ssl", "strict")
	require.Len(t, stub.Resources(recordResourceType), 1)
}

func TestCreateBucket(t *testing.T) {
	stub := &test.ResourceStub{}
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket := &cf.Bucket{
			Name:   "test",
			Region: "eeur",
			Zone: &cf.BucketZone{
				ID:     "test",
				Domain: "www.test.com",
			},
		}

		err := cf.CreateBucket(ctx, bucket)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", stub))
	require.NoError(t, err)
	requireBucket(t, stub.Resources(r2BucketResourceType), "test", "eeur")
	requireCustomDomain(t, stub.Resources(r2CustomDomainResourceType), "test", "www.test.com", "test")
	require.Empty(t, stub.Resources(r2LifecycleResourceType))

	stub = &test.ResourceStub{}
	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket := &cf.Bucket{
			Name:   "test",
			Region: "eeur",
		}

		err := cf.CreateBucket(ctx, bucket)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", stub))
	require.NoError(t, err)
	requireBucket(t, stub.Resources(r2BucketResourceType), "test", "eeur")
	require.Empty(t, stub.Resources(r2CustomDomainResourceType))
	require.Empty(t, stub.Resources(r2LifecycleResourceType))

	stub = &test.ResourceStub{}
	stub.FailResource(r2BucketResourceType)
	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket := &cf.Bucket{
			Name:       "test",
			Region:     "eeur",
			Expiration: new("24h"),
		}

		return cf.CreateBucket(ctx, bucket)
	}, pulumi.WithMocks("project", "stack", stub))
	require.Error(t, err)
	require.Len(t, stub.Resources(r2BucketResourceType), 1)
	require.Empty(t, stub.Resources(r2LifecycleResourceType), "lifecycle requires a successfully created bucket")
}

func TestBucketExpiration(t *testing.T) {
	for _, tt := range []struct {
		expiration *string
		name       string
		seconds    int64
		public     bool
	}{
		{name: "omitted"},
		{name: "seven days", expiration: new("168h"), seconds: 604800},
		{name: "one day with custom domain", expiration: new("24h"), seconds: 86400, public: true},
		{name: "minimum", expiration: new("1s"), seconds: 1},
		{name: "milliseconds forming a whole second", expiration: new("1000ms"), seconds: 1},
		{name: "fractional hours", expiration: new("1.5h"), seconds: 5400},
		{name: "maximum whole seconds", expiration: new("2562047h47m16s"), seconds: 9223372036},
	} {
		t.Run(tt.name, func(t *testing.T) {
			configuration, err := cf.ReadConfiguration("../test/cf.hjson")
			require.NoError(t, err)
			configured := configuration.GetBuckets()[0]
			configured.Expiration = tt.expiration
			if !tt.public {
				configured.Zone = nil
			}

			file := filepath.Join(t.TempDir(), "cf.hjson")
			require.NoError(t, os.WriteFile(file, nil, 0o600))
			require.NoError(t, config.Write(file, configuration))
			rewritten, err := cf.ReadConfiguration(file)
			require.NoError(t, err)
			require.True(t, proto.Equal(configuration, rewritten), "HJSON rewrites must preserve the full configuration")
			require.Equal(t, tt.expiration, rewritten.GetBuckets()[0].Expiration)
			bucket := cf.ConvertBucket(rewritten.GetBuckets()[0])
			require.Equal(t, tt.expiration, bucket.Expiration)
			stub := &test.ResourceStub{}
			err = test.RunWithMocks(func(ctx *pulumi.Context) error {
				return cf.CreateBucket(ctx, bucket)
			}, stub)
			require.NoError(t, err)
			requireBucket(t, stub.Resources(r2BucketResourceType), "atlas-assets", "EEUR")
			if tt.public {
				requireCustomDomain(t, stub.Resources(r2CustomDomainResourceType), "atlas-assets", "assets.example.net", "test-zone-id")
			} else {
				require.Empty(t, stub.Resources(r2CustomDomainResourceType))
			}
			if tt.expiration == nil {
				require.Empty(t, stub.Resources(r2LifecycleResourceType))
				return
			}
			requireLifecycle(t, stub.Resources(r2LifecycleResourceType), "atlas-assets", tt.seconds)
		})
	}
}

func TestBucketRejectsInvalidExpiration(t *testing.T) {
	for _, tt := range []struct {
		name       string
		expiration string
	}{
		{name: "empty", expiration: ""},
		{name: "zero", expiration: "0s"},
		{name: "negative", expiration: "-1s"},
		{name: "fractional seconds", expiration: "1.5s"},
		{name: "less than a second", expiration: "999ms"},
		{name: "unsupported day unit", expiration: "7d"},
		{name: "malformed", expiration: "invalid"},
		{name: "overflow", expiration: "2562048h"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			err := test.RunWithMocks(func(ctx *pulumi.Context) error {
				return cf.CreateBucket(ctx, cf.ConvertBucket(&v2.Bucket{
					Name:       "test",
					Region:     "EEUR",
					Expiration: new(tt.expiration),
				}))
			}, stub)
			require.ErrorContains(t, err, "bucket test: expiration")
			require.Empty(t, stub.Resources(r2BucketResourceType))
			require.Empty(t, stub.Resources(r2LifecycleResourceType))
			require.Empty(t, stub.Resources(r2CustomDomainResourceType))
		})
	}
}

func TestBucketRejectsNonStringExpiration(t *testing.T) {
	for _, value := range []string{"86400", "1.5", "true"} {
		t.Run(value, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "cf.hjson")
			contents := `{"buckets":[{"name":"test","expiration":` + value + `}]}`
			require.NoError(t, os.WriteFile(file, []byte(contents), 0o600))
			_, err := cf.ReadConfiguration(file)
			require.Error(t, err)
		})
	}
}

func TestCreateBucketReturnsLifecycleError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(r2LifecycleResourceType)
	configuration, err := cf.ReadConfiguration("../test/cf.hjson")
	require.NoError(t, err)
	err = test.RunWithMocks(func(ctx *pulumi.Context) error {
		return cf.CreateBucket(ctx, cf.ConvertBucket(configuration.GetBuckets()[2]))
	}, stub)
	require.Error(t, err)
	requireBucket(t, stub.Resources(r2BucketResourceType), "atlas-archive", "EEUR")
	require.Len(t, stub.Resources(r2LifecycleResourceType), 1)
}

func TestCreateBucketReturnsCustomDomainError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(r2CustomDomainResourceType)

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket := &cf.Bucket{
			Name:   "test",
			Region: "eeur",
			Zone: &cf.BucketZone{
				ID:     "test",
				Domain: "www.test.com",
			},
		}

		return cf.CreateBucket(ctx, bucket)
	}, pulumi.WithMocks("project", "stack", stub))
	require.Error(t, err)
	requireBucket(t, stub.Resources(r2BucketResourceType), "test", "eeur")
	require.Len(t, stub.Resources(r2CustomDomainResourceType), 1)
}

func createBalancerZone(t *testing.T, mocks pulumi.MockResourceMonitor) error {
	t.Helper()

	return test.RunWithMocks(func(ctx *pulumi.Context) error {
		zone := &cf.BalancerZone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			IPV4:        "127.0.0.1",
			IPV6:        "::1",
		}

		return cf.CreateBalancerZone(ctx, zone)
	}, mocks)
}

func createPageZone(t *testing.T, mocks pulumi.MockResourceMonitor) error {
	t.Helper()

	return test.RunWithMocks(func(ctx *pulumi.Context) error {
		zone := &cf.PageZone{
			Name:   "test",
			Domain: "test.com",
			Host:   "test.github.io",
		}

		return cf.CreatePageZone(ctx, zone)
	}, mocks)
}

func requireRecord(t *testing.T, records []resource.PropertyMap, kind, name, content string) {
	t.Helper()

	for _, record := range records {
		if test.Property(t, record, "type").StringValue() != kind {
			continue
		}
		if test.Property(t, record, "name").StringValue() != name {
			continue
		}

		require.Equal(t, content, test.Property(t, record, "content").StringValue())
		require.True(t, test.Property(t, record, "proxied").BoolValue())
		require.Equal(t, 1, int(test.Property(t, record, "ttl").NumberValue()))
		return
	}

	require.Failf(t, "missing DNS record", "%s %s", kind, name)
}

func requireBucket(t *testing.T, buckets []resource.PropertyMap, name, region string) {
	t.Helper()

	require.Len(t, buckets, 1)

	bucket := buckets[0]
	require.Equal(t, name, test.Property(t, bucket, "name").StringValue())
	require.Equal(t, region, test.Property(t, bucket, "location").StringValue())
	require.Equal(t, "Standard", test.Property(t, bucket, "storageClass").StringValue())
}

func requireLifecycle(t *testing.T, lifecycles []resource.PropertyMap, bucket string, seconds int64) {
	t.Helper()

	require.Len(t, lifecycles, 1)
	lifecycle := lifecycles[0]
	require.Equal(t, bucket, test.Property(t, lifecycle, "bucketName").StringValue())
	rules := test.Property(t, lifecycle, "rules").ArrayValue()
	require.Len(t, rules, 1)
	rule := rules[0].ObjectValue()
	require.Equal(t, "expire-objects", test.Property(t, rule, "id").StringValue())
	require.True(t, test.Property(t, rule, "enabled").BoolValue(), "object expiration must be enabled")
	conditions := test.Property(t, rule, "conditions").ObjectValue()
	require.Empty(t, test.Property(t, conditions, "prefix").StringValue(), "expiration must cover all objects")
	deletion := test.Property(t, rule, "deleteObjectsTransition").ObjectValue()
	condition := test.Property(t, deletion, "condition").ObjectValue()
	require.Equal(t, "Age", test.Property(t, condition, "type").StringValue())
	require.Equal(t, seconds, int64(test.Property(t, condition, "maxAge").NumberValue()))
	multipart := test.Property(t, rule, "abortMultipartUploadsTransition").ObjectValue()
	condition = test.Property(t, multipart, "condition").ObjectValue()
	require.Equal(t, "Age", test.Property(t, condition, "type").StringValue())
	require.Equal(t, int64(604800), int64(test.Property(t, condition, "maxAge").NumberValue()), "retain default multipart cleanup")
}

func requireCustomDomain(t *testing.T, domains []resource.PropertyMap, bucket, domain, zone string) {
	t.Helper()

	require.Len(t, domains, 1)

	customDomain := domains[0]
	require.Equal(t, bucket, test.Property(t, customDomain, "bucketName").StringValue())
	require.Equal(t, domain, test.Property(t, customDomain, "domain").StringValue())
	require.True(t, test.Property(t, customDomain, "enabled").BoolValue())
	require.Equal(t, zone, test.Property(t, customDomain, "zoneId").StringValue())
	require.Equal(t, "1.2", test.Property(t, customDomain, "minTls").StringValue())
}

func requireZoneSetting(t *testing.T, settings []resource.PropertyMap, name, value string) {
	t.Helper()

	for _, setting := range settings {
		if test.Property(t, setting, "settingId").StringValue() != name {
			continue
		}

		require.Equal(t, value, test.Property(t, setting, "value").StringValue())
		return
	}

	require.Failf(t, "missing zone setting", "%s", name)
}
