# Run kubescore for bezeichner.
kube-score-bezeichner:
	$(MAKE) namespace=bezeichner kube-score-namespace

# Delete bezeichner.
delete-bezeichner:
	$(MAKE) namespace=bezeichner delete-namespace

# Setup bezeichner.
setup-bezeichner:
	$(MAKE) namespace=bezeichner setup-namespace
	$(MAKE) namespace=bezeichner setup-otlp-secret
	$(MAKE) namespace=bezeichner setup-konfig-secret

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
	curl -svf --header "Content-Type: application/json" --request POST --data '{ "application": "uuid", "count": 10 }'  https://bezeichner.lean-thoughts.com/v1/generate
