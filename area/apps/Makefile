otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

include standort/include.mk
include konfig/include.mk
include bezeichner/include.mk

# Lint config.
lint:
	kube-linter lint .

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

# Preview all apps.
preview: preview-standort preview-konfig preview-bezeichner

# Update all apps.
update: update-standort update-konfig update-bezeichner
