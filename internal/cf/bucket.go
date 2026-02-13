package cf

import (
	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type (
	// BucketZone identifies the Cloudflare zone used to attach a public/custom domain to an R2 bucket.
	//
	// This is optional; when provided, CreateBucket will also provision an R2 custom domain.
	// See: https://developers.cloudflare.com/r2/buckets/public-buckets/
	BucketZone struct {
		// ID is the Cloudflare zone identifier.
		ID string
		// Domain is the fully-qualified domain name to associate with the bucket.
		Domain string
	}

	// Bucket represents an R2 bucket configuration in Cloudflare.
	//
	// Regions are defined under: https://developers.cloudflare.com/r2/reference/data-location/#available-hints
	Bucket struct {
		// Zone is optional; when nil, no custom domain is created for the bucket.
		Zone *BucketZone
		// Name is the R2 bucket name.
		Name string
		// Region is the R2 data location hint.
		Region string
	}
)

// ConvertBucket converts a protobuf v2.Bucket into the internal Bucket model.
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

// CreateBucket provisions an R2 bucket and, optionally, an R2 custom domain.
//
// If bucket.Zone is nil, CreateBucket only creates the bucket.
// If bucket.Zone is non-nil, CreateBucket also creates an R2 custom domain using Zone.ID and Zone.Domain.
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

// newDomain creates an R2 custom domain for bucket.
//
// It enables the domain and sets a minimum TLS version of 1.2.
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
