otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

include standort/include.mk
include konfig/include.mk
include bezeichner/include.mk

kube-linter:
	kube-linter lint .

kube-score: kube-score-konfig kube-score-standort kube-score-bezeichner

# Lint config.
lint: kube-linter kube-score

# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Setup otlp.
setup-otlp:
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace $(namespace)

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

# Preview all apps.
preview: preview-konfig preview-standort preview-bezeichner

# Update all apps.
update: update-konfig update-standort update-bezeichner

# Rollout all apps.
rollout: rollout-konfig rollout-standort rollout-bezeichner
