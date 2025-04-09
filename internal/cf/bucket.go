package cf

import (
	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Bucket represents an R2 bucket in Cloudflare.
// Regions are defined under https://developers.cloudflare.com/r2/reference/data-location/#available-hints.
type Bucket struct {
	Name   string
	Region string
}

// ConvertBucket converts a v2.Bucket to a Bucket.
func ConvertBucket(bucket *v2.Bucket) *Bucket {
	return &Bucket{
		Name:   bucket.GetName(),
		Region: bucket.GetRegion(),
	}
}

// CreateBucket for cf.
func CreateBucket(ctx *pulumi.Context, bucket *Bucket) error {
	args := &cloudflare.R2BucketArgs{
		AccountId: account,
		Name:      pulumi.String(bucket.Name),
		Location:  pulumi.String(bucket.Region),
	}

	_, err := cloudflare.NewR2Bucket(ctx, bucket.Name, args)

	return err
}
