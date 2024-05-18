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
