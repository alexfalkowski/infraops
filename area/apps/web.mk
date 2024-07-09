# Run kubescore for web.
kube-score-web:
	$(MAKE) namespace=web kube-score-namespace

# Delete web.
delete-web:
	$(MAKE) namespace=web delete-namespace

# Setup web.
setup-web:
	$(MAKE) namespace=web setup-namespace
	$(MAKE) namespace=web setup-otlp-secret
	$(MAKE) namespace=web setup-konfig-secret

# Preview web.
preview-web:
	$(MAKE) namespace=web preview-namespace

# Update web.
update-web:
	$(MAKE) namespace=web update-namespace

# Rollout web.
rollout-web:
	$(MAKE) namespace=web rollout-namespace

# Verify web.
verify-web:
	curl -svf https://web.lean-thoughts.com
