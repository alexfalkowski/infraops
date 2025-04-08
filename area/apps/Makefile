include lean.mk

# Run kube score.
kube-score: kube-score-lean

# Run kubescape.
kubescape: kubescape-lean

# Lint config.
lint: kube-score kubescape

# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save 9bccba18-cc46-4b4b-9692-5cbb74d75b04

# Setup all apps.
setup: setup-lean

# Delete all apps.
delete: delete-lean

# Rollout all apps.
rollout: rollout-lean

# Verify all apps.
verify: verify-lean
