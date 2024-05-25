include secrets.mk
include konfig.mk
include standort/include.mk
include bezeichner/include.mk

kube-linter:
	kube-linter lint .

kube-score: kube-score-konfig kube-score-standort kube-score-bezeichner

# Lint config.
lint: kube-linter kube-score

# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Delete namespace.
delete-namespace:
	kubectl delete namespaces $(namespace)

# Setup namespace.
setup-namespace:
	kubectl create namespace $(namespace)

# Preview namespace.
preview-namespace:
	kubectl apply -f $(namespace) -n $(namespace) --dry-run=server

# Update namespace.
update-namespace:
	kubectl apply -f $(namespace) -n $(namespace)

# Rollout a namespace.
rollout-namespace:
	kubectl rollout restart deployment/$(namespace) -n $(namespace)

# Run kube-score for a namespace.
kube-score-namespace:
	kube-score score --ignore-test deployment-has-host-podantiaffinity $(namespace)/*.yaml

# Delete all apps.
delete: delete-konfig delete-standort delete-bezeichner

# Preview all apps.
preview: preview-konfig preview-standort preview-bezeichner

# Update all apps.
update: update-konfig update-standort update-bezeichner

# Rollout all apps.
rollout: rollout-konfig rollout-standort rollout-bezeichner

# Verify all apps.
verify: verify-standort verify-bezeichner
