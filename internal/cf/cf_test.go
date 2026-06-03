package cf_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/cf"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

const (
	r2BucketResourceType       = "cloudflare:index/r2Bucket:R2Bucket"
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

func TestCreatePagerZone(t *testing.T) {
	stub := &test.ResourceStub{}
	require.NoError(t, createPageZone(t, stub))
	requireRecord(t, stub.Resources(recordResourceType), "CNAME", "www.test.com", "test.github.io")
	requireZoneSetting(t, stub.Resources(zoneSettingResourceType), "ssl", "strict")

	stub = &test.ResourceStub{}
	stub.FailResource(zoneSettingResourceType)
	require.Error(t, createPageZone(t, stub))
}

func TestCreatePagerZoneReturnsRecordError(t *testing.T) {
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

	stub = &test.ResourceStub{}
	stub.FailResource(r2BucketResourceType)
	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket := &cf.Bucket{
			Name:   "test",
			Region: "eeur",
		}

		return cf.CreateBucket(ctx, bucket)
	}, pulumi.WithMocks("project", "stack", stub))
	require.Error(t, err)
	require.Len(t, stub.Resources(r2BucketResourceType), 1)
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
		z := &cf.BalancerZone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			IPV4:        "127.0.0.1",
			IPV6:        "::1",
		}

		return cf.CreateBalancerZone(ctx, z)
	}, mocks)
}

func createPageZone(t *testing.T, mocks pulumi.MockResourceMonitor) error {
	t.Helper()

	return test.RunWithMocks(func(ctx *pulumi.Context) error {
		z := &cf.PageZone{
			Name:   "test",
			Domain: "test.com",
			Host:   "test.github.io",
		}

		return cf.CreatePageZone(ctx, z)
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
