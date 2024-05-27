otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

# Setup otlp secret.
setup-otlp-secret:
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace $(namespace)

# Setup konfig secret.
setup-konfig-secret:
	kubectl create secret generic konfig-secret --from-literal=token=$(KONFIG_TOKEN) --namespace $(namespace)
