# Run kubescore for konfig.
kube-score-konfig:
	$(MAKE) namespace=konfig kube-score-namespace

# Delete konfig.
delete-konfig:
	$(MAKE) namespace=konfig delete-namespace

# Setup konfig.
setup-konfig:
	$(MAKE) namespace=konfig setup-namespace
	kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace konfig

# Preview konfig.
preview-konfig:
	$(MAKE) namespace=konfig preview-namespace

# Update konfig.
update-konfig:
	$(MAKE) namespace=konfig update-namespace

# Rollout konfig.
rollout-konfig:
	$(MAKE) namespace=konfig rollout-namespace
