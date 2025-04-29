package cf

import (
	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type (
	// BucketZone defines a zone to be used by buckets.
	// This is optional and defined under https://developers.cloudflare.com/r2/buckets/public-buckets/
	BucketZone struct {
		ID     string
		Domain string
	}

	// Bucket represents an R2 bucket in Cloudflare.
	// Regions are defined under https://developers.cloudflare.com/r2/reference/data-location/#available-hints.
	Bucket struct {
		Zone   *BucketZone
		Name   string
		Region string
	}
)

// ConvertBucket converts a v2.Bucket to a Bucket.
func ConvertBucket(bucket *v2.Bucket) *Bucket {
	b := &Bucket{
		Name:   bucket.GetName(),
		Region: bucket.GetRegion(),
	}

	zone := bucket.GetZone()
	if zone != nil {
		b.Zone = &BucketZone{
			ID:     zone.GetId(),
			Domain: zone.GetDomain(),
		}
	}

	return b
}

// CreateBucket for cf.
func CreateBucket(ctx *pulumi.Context, bucket *Bucket) error {
	args := &cloudflare.R2BucketArgs{
		AccountId:    account,
		Name:         pulumi.String(bucket.Name),
		Location:     pulumi.String(bucket.Region),
		StorageClass: pulumi.String("Standard"),
	}

	if _, err := cloudflare.NewR2Bucket(ctx, bucket.Name, args); err != nil {
		return err
	}

	if bucket.Zone == nil {
		return nil
	}

	return newDomain(ctx, bucket)
}

func newDomain(ctx *pulumi.Context, bucket *Bucket) error {
	args := &cloudflare.R2CustomDomainArgs{
		AccountId:  account,
		BucketName: pulumi.String(bucket.Name),
		Domain:     pulumi.String(bucket.Zone.Domain),
		Enabled:    inputs.Yes,
		ZoneId:     pulumi.String(bucket.Zone.ID),
		MinTls:     pulumi.String("1.2"),
	}
	_, err := cloudflare.NewR2CustomDomain(ctx, bucket.Name, args)

	return err
}
