# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Delete all.
delete:
	kubectl delete namespaces nginx-ingress grafana circleci

# Setup everything.
setup: setup-helm install-helm setup-grafana

# Setup helm repos.
setup-helm:
	helm repo add nginx https://kubernetes.github.io/ingress-nginx
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add circleci https://circleci-public.github.io/cci-k8s-release-agent
	helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
	helm repo update

# Setup grafana.
setup-grafana:
	kubectl create secret generic metrics-secret --from-literal=username=1289421 --from-literal=password=$(GRAFANA_METRICS_TOKEN) --namespace grafana
	kubectl create secret generic logs-secret --from-literal=username=743756 --from-literal=password=$(GRAFANA_LOGS_TOKEN) --namespace grafana

# Install helm.
install-helm:
	helm upgrade --install nginx-ingress-release nginx/ingress-nginx --namespace nginx-ingress --create-namespace --version 4.11.4 -f nginx/values.yaml
	helm upgrade --install alloy-metrics-release grafana/alloy --create-namespace -n grafana -f grafana/metrics.yaml --set-file alloy.configMap.content=grafana/metrics.alloy --version 0.11.0
	helm upgrade --install alloy-logs-release grafana/alloy --create-namespace -n grafana -f grafana/logs.yaml --set-file alloy.configMap.content=grafana/logs.alloy --version 0.11.0
	helm upgrade --install circleci-release circleci/circleci-release-agent --set tokenSecret.token=$(CIRCLECI_K8S_TOKEN) --create-namespace --namespace circleci --set managedNamespaces="{lean}" --version v1.3.3
	helm upgrade --install metrics-server-release metrics-server/metrics-server --create-namespace --namespace metrics-server --version 3.12.2

# All pods.
pods:
	kubectl get pods --all-namespaces
