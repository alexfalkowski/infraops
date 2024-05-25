# Run kubescore for konfig.
kube-score-konfig:
	$(MAKE) namespace=konfig kube-score-namespace

# Delete konfig.
delete-konfig:
	$(MAKE) namespace=konfig delete-namespace

# Setup konfig.
setup-konfig:
	$(MAKE) namespace=konfig setup-namespace
	$(MAKE) namespace=konfig setup-otlp
	kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace konfig

# Rollout konfig.
rollout-konfig:
	$(MAKE) namespace=konfig rollout-namespace
