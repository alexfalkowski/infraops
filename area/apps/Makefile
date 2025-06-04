include lean.mk

# Run kube score.
kube-score: kube-score-lean

# Run kubescape.
kubescape: kubescape-lean

# Lint config.
lint: kube-score kubescape

# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save 291cdfec-94f7-40d5-9552-2e3a3e3ded17

# Setup all apps.
setup: setup-lean

# Delete all apps.
delete: delete-lean

# Rollout all apps.
rollout: rollout-lean

# Verify all apps.
verify: verify-lean

# Load all apps.
load: load-lean
