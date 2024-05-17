otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Delete all.
delete: delete-standort

# Setup all apps.
setup: setup-standort

# Setup otlp.
setup-otlp:
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace $(namespace)

# Delete standort.
delete-standort:
	kubectl delete namespaces standort

# Setup standort.
setup-standort:
	kubectl create namespace standort

# Preview standort.
preview-standort:
	kubectl apply -f standort -n standort --dry-run=server

# Update standort.
update-standort:
	kubectl apply -f standort -n standort
