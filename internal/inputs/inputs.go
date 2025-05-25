package inputs

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

const (
	// Always for inputs.
	Always = pulumi.String("Always")

	// On for inputs.
	On = pulumi.String("on")

	// Off for inputs.
	Off = pulumi.String("off")

	// Yes for inputs.
	Yes = pulumi.Bool(true)

	// No for inputs.
	No = pulumi.Bool(false)

	// Enabled for inputs.
	Enabled = pulumi.String("enabled")

	// False for inputs.
	False = pulumi.String("false")

	// One for inputs.
	One = pulumi.Int(1)

	// Automatic for inputs.
	Automatic = pulumi.Float64(1)
)
