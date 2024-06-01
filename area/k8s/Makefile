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
	helm repo update

# Setup grafana.
setup-grafana:
	kubectl create secret generic metrics-secret --from-literal=username=1289421 --from-literal=password=$(GRAFANA_METRICS_TOKEN) --namespace grafana
	kubectl create secret generic logs-secret --from-literal=username=743756 --from-literal=password=$(GRAFANA_LOGS_TOKEN) --namespace grafana

# Update grafana.
update-grafana:
	kubectl apply -f grafana --namespace grafana

# Install helm.
install-helm:
	helm upgrade --install nginx-ingress-release nginx/ingress-nginx --namespace nginx-ingress --create-namespace --version 4.10.1 -f nginx/values.yaml
	helm upgrade --install grafana-agent-release --create-namespace grafana/grafana-agent-operator -n grafana --version v0.3.23
	helm upgrade --install circleci-release circleci/circleci-release-agent --set tokenSecret.token=$(CIRCLECI_K8S_TOKEN) --create-namespace --namespace circleci --set managedNamespaces="{konfig,standort,bezeichner}" --version v1.2.1

# Update deps.
update: update-grafana

# All pods.
pods:
	kubectl get pods --all-namespaces
