# Run kubescore for bezeichner.
kube-score-bezeichner:
	$(MAKE) namespace=bezeichner kube-score-namespace

# Delete bezeichner.
delete-bezeichner:
	$(MAKE) namespace=bezeichner delete-namespace

# Setup bezeichner.
setup-bezeichner:
	$(MAKE) namespace=bezeichner setup-namespace

# Preview bezeichner.
preview-bezeichner:
	$(MAKE) namespace=bezeichner preview-namespace

# Update bezeichner.
update-bezeichner:
	$(MAKE) namespace=bezeichner update-namespace

# Rollout bezeichner.
rollout-bezeichner:
	$(MAKE) namespace=bezeichner rollout-namespace

# Verify bezeichner.
verify-bezeichner:
	curl -sf https://bezeichner.lean-thoughts.com/v1/generate/uuid/1
