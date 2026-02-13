package inputs

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

const (
	// Always is the string value "Always" as a Pulumi input.
	Always = pulumi.String("Always")

	// On is the string value "on" as a Pulumi input.
	On = pulumi.String("on")

	// Off is the string value "off" as a Pulumi input.
	Off = pulumi.String("off")

	// Yes is the boolean value true as a Pulumi input.
	Yes = pulumi.Bool(true)

	// No is the boolean value false as a Pulumi input.
	No = pulumi.Bool(false)

	// Enabled is the string value "enabled" as a Pulumi input.
	Enabled = pulumi.String("enabled")

	// False is the string value "false" as a Pulumi input.
	False = pulumi.String("false")

	// One is the integer value 1 as a Pulumi input.
	One = pulumi.Int(1)

	// Automatic is the float64 value 1 as a Pulumi input.
	Automatic = pulumi.Float64(1)
)
