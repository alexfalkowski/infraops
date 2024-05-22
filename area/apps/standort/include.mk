# Run kubescore for standort.
kube-score-standort:
	$(MAKE) namespace=standort kube-score-namespace

# Delete standort.
delete-standort:
	$(MAKE) namespace=standort delete-namespace

# Setup standort.
setup-standort:
	$(MAKE) namespace=standort setup-namespace

# Preview standort.
preview-standort:
	$(MAKE) namespace=standort preview-namespace

# Update standort.
update-standort:
	$(MAKE) namespace=standort update-namespace

# Rollout standort.
rollout-standort:
	$(MAKE) namespace=standort rollout-namespace
