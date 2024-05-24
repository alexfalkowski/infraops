# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Delete all.
delete:
	kubectl delete namespaces nginx-ingress grafana

# Setup everything.
setup: setup-helm install-helm setup-grafana

# Setup helm repos.
setup-helm:
	helm repo add nginx-stable https://helm.nginx.com/stable
	helm repo add grafana https://grafana.github.io/helm-charts
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
	helm install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --version 1.2.1
	helm install grafana-agent-release --create-namespace grafana/grafana-agent-operator -n grafana --version v0.3.22

# Update deps.
update: update-grafana

# All pods.
pods:
	kubectl get pods --all-namespaces
