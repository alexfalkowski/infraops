otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

# Setup otlp.
setup-otlp:
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace $(namespace)
