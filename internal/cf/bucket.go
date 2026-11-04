package cf

import (
	"fmt"
	"time"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// BucketZone identifies the Cloudflare zone used to attach a public/custom domain to an R2 bucket.
//
// This is optional; when provided, CreateBucket will also provision an R2 custom domain.
// See: https://developers.cloudflare.com/r2/buckets/public-buckets/
type BucketZone struct {
	// ID is the Cloudflare zone identifier.
	ID string

	// Domain is the fully-qualified domain name to associate with the bucket.
	Domain string
}

// Bucket represents an R2 bucket configuration in Cloudflare.
//
// Regions are defined under: https://developers.cloudflare.com/r2/reference/data-location/#available-hints
type Bucket struct {
	// Zone is optional; when nil, no custom domain is created for the bucket.
	Zone *BucketZone

	// Expiration optionally expires all objects by age using Go duration syntax.
	// Configured values must be positive whole seconds within the time.Duration range.
	// Nil creates no lifecycle resource, but does not remove previously applied remote rules.
	Expiration *string

	// Name is the R2 bucket name.
	Name string

	// Region is the R2 data location hint.
	Region string
}

// GetExpiration parses the configured Go duration, returning zero when omitted.
// Invalid durations and values that are not positive whole seconds return an error.
func (b *Bucket) GetExpiration() (time.Duration, error) {
	if b.Expiration == nil {
		return 0, nil
	}

	expiration, err := time.ParseDuration(*b.Expiration)
	if err != nil {
		return 0, fmt.Errorf("bucket %s: expiration: %w", b.Name, err)
	}
	if expiration <= 0 || expiration%time.Second != 0 {
		return 0, fmt.Errorf("bucket %s: expiration must be a positive whole number of seconds", b.Name)
	}
	return expiration, nil
}

// ConvertBucket converts a protobuf v2.Bucket into the internal Bucket model.
func ConvertBucket(bucket *v2.Bucket) *Bucket {
	b := &Bucket{
		Name:       bucket.GetName(),
		Region:     bucket.GetRegion(),
		Expiration: bucket.Expiration,
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

// CreateBucket provisions an R2 bucket with optional object expiration and custom domain.
//
// Explicit expiration is parsed with time.ParseDuration and must be positive whole seconds.
// Invalid durations return an error before registering resources.
// The lifecycle owns all bucket rules and retains seven-day incomplete multipart cleanup.
// Removing its configuration leaves remote rules active with the pinned provider.
// A non-nil Zone also creates a custom domain using Zone.ID and Zone.Domain.
func CreateBucket(ctx *pulumi.Context, bucket *Bucket) error {
	expiration, err := bucket.GetExpiration()
	if err != nil {
		return err
	}

	args := &cloudflare.R2BucketArgs{
		AccountId:    account,
		Name:         pulumi.String(bucket.Name),
		Location:     pulumi.String(bucket.Region),
		StorageClass: pulumi.String("Standard"),
	}

	b, err := cloudflare.NewR2Bucket(ctx, bucket.Name, args)
	if err != nil {
		return err
	}

	if bucket.Expiration != nil {
		if err := newLifecycle(ctx, bucket, b, expiration); err != nil {
			return err
		}
	}
	if bucket.Zone == nil {
		return nil
	}

	return newDomain(ctx, bucket)
}

// newLifecycle replaces the bucket's lifecycle rules after the bucket exists.
func newLifecycle(ctx *pulumi.Context, bucket *Bucket, b *cloudflare.R2Bucket, expiration time.Duration) error {
	args := &cloudflare.R2BucketLifecycleArgs{
		AccountId:  account,
		BucketName: b.Name,
		Rules: cloudflare.R2BucketLifecycleRuleArray{
			&cloudflare.R2BucketLifecycleRuleArgs{
				Id:      pulumi.String("expire-objects"),
				Enabled: inputs.Yes,
				Conditions: &cloudflare.R2BucketLifecycleRuleConditionsArgs{
					Prefix: pulumi.String(""),
				},
				DeleteObjectsTransition: &cloudflare.R2BucketLifecycleRuleDeleteObjectsTransitionArgs{
					Condition: &cloudflare.R2BucketLifecycleRuleDeleteObjectsTransitionConditionArgs{
						Type:   pulumi.String("Age"),
						MaxAge: pulumi.Int(expiration / time.Second),
					},
				},
				// Replacing the rule set must retain R2's default incomplete multipart cleanup.
				AbortMultipartUploadsTransition: &cloudflare.R2BucketLifecycleRuleAbortMultipartUploadsTransitionArgs{
					Condition: &cloudflare.R2BucketLifecycleRuleAbortMultipartUploadsTransitionConditionArgs{
						Type:   pulumi.String("Age"),
						MaxAge: pulumi.Int(604800),
					},
				},
			},
		},
	}

	_, err := cloudflare.NewR2BucketLifecycle(ctx, bucket.Name, args)
	return err
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
