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

	// LockDuration optionally protects all objects from deletion and overwrites by age.
	// Configured values must be positive whole seconds in Go duration syntax.
	// Nil creates no lock resource. Existing remote rules persist when the field is removed.
	LockDuration *string

	// Name is the R2 bucket name.
	Name string

	// Region is the R2 data location hint.
	Region string
}

// GetExpiration parses the configured Go duration, returning zero when omitted.
// Invalid durations and values that are not positive whole seconds return an error.
func (b *Bucket) GetExpiration() (time.Duration, error) {
	return parseBucketDuration(b.Name, "expiration", b.Expiration)
}

// ConvertBucket converts a protobuf v2.Bucket into the internal Bucket model.
func ConvertBucket(bucket *v2.Bucket) *Bucket {
	b := &Bucket{
		Name:         bucket.GetName(),
		Region:       bucket.GetRegion(),
		Expiration:   bucket.Expiration,
		LockDuration: bucket.LockDuration,
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

// CreateBucket provisions an R2 bucket with optional object expiration, lock, and custom domain.
//
// Explicit expiration and lock durations use Go duration syntax and must be positive whole seconds.
// Invalid durations return an error before registering resources.
// The lifecycle owns the full lifecycle rule set and keeps seven-day incomplete multipart cleanup.
// The lock owns the full bucket lock rule set. The pinned provider leaves remote rules active
// when the corresponding lifecycle or lock setting is removed.
// A non-nil Zone also creates a custom domain using Zone.ID and Zone.Domain.
func CreateBucket(ctx *pulumi.Context, bucket *Bucket) error {
	expiration, err := bucket.GetExpiration()
	if err != nil {
		return err
	}
	lockDuration, err := parseBucketDuration(bucket.Name, "lock_duration", bucket.LockDuration)
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
		if err := newLifecycle(ctx, b, expiration, bucket); err != nil {
			return err
		}
	}
	if bucket.LockDuration != nil {
		if err := newLock(ctx, b, bucket.Name, lockDuration); err != nil {
			return err
		}
	}
	if bucket.Zone == nil {
		return nil
	}

	return newDomain(ctx, b, bucket)
}

func parseBucketDuration(bucket, field string, value *string) (time.Duration, error) {
	if value == nil {
		return 0, nil
	}

	duration, err := time.ParseDuration(*value)
	if err != nil {
		return 0, fmt.Errorf("bucket %s: %s: %w", bucket, field, err)
	}
	if duration <= 0 || duration%time.Second != 0 {
		return 0, fmt.Errorf("bucket %s: %s must be a positive whole number of seconds", bucket, field)
	}
	return duration, nil
}

// newLifecycle replaces the bucket's lifecycle rules after the bucket exists.
func newLifecycle(ctx *pulumi.Context, created *cloudflare.R2Bucket, expiration time.Duration, bucket *Bucket) error {
	args := &cloudflare.R2BucketLifecycleArgs{
		AccountId:  account,
		BucketName: created.Name,
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

func newLock(ctx *pulumi.Context, created *cloudflare.R2Bucket, name string, duration time.Duration) error {
	args := &cloudflare.R2BucketLockArgs{
		AccountId:  account,
		BucketName: created.Name,
		Rules: cloudflare.R2BucketLockRuleArray{
			&cloudflare.R2BucketLockRuleArgs{
				Id:      pulumi.String("protect-objects"),
				Enabled: inputs.Yes,
				Prefix:  pulumi.String(""),
				Condition: &cloudflare.R2BucketLockRuleConditionArgs{
					Type:          pulumi.String("Age"),
					MaxAgeSeconds: pulumi.Int(duration / time.Second),
				},
			},
		},
	}

	_, err := cloudflare.NewR2BucketLock(ctx, name, args)
	return err
}

// newDomain creates an R2 custom domain after the bucket exists.
//
// It enables the domain and sets a minimum TLS version of 1.2.
func newDomain(ctx *pulumi.Context, created *cloudflare.R2Bucket, bucket *Bucket) error {
	args := &cloudflare.R2CustomDomainArgs{
		AccountId:  account,
		BucketName: created.Name,
		Domain:     pulumi.String(bucket.Zone.Domain),
		Enabled:    inputs.Yes,
		ZoneId:     pulumi.String(bucket.Zone.ID),
		MinTls:     pulumi.String("1.2"),
	}

	_, err := cloudflare.NewR2CustomDomain(ctx, bucket.Name, args)
	return err
}
