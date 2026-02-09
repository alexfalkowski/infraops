# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Delete all.
delete:
	@kubectl delete namespaces nginx-ingress circleci metrics-server

# Setup everything.
setup: setup-helm install-helm

# Setup helm repos.
setup-helm:
	@helm repo add nginx https://kubernetes.github.io/ingress-nginx
	@helm repo add circleci https://circleci-public.github.io/cci-k8s-release-agent
	@helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
	@helm repo update

# Install circleci.
install-circleci:
	@helm upgrade --install circleci-release circleci/circleci-release-agent --set tokenSecret.token=$(CIRCLECI_K8S_TOKEN) --create-namespace --namespace circleci --set managedNamespaces="{lean}" --version v1.3.7

# Install metrics server.
install-metrics-server:
	@helm upgrade --install metrics-server-release metrics-server/metrics-server --create-namespace --namespace metrics-server --version 3.13.0

# Install nginx ingress controller.
install-nginx-ingress:
	@helm upgrade --install nginx-ingress-release nginx/ingress-nginx --namespace nginx-ingress --create-namespace --version 4.14.1 -f nginx/values.yaml

# Install helm.
install-helm: install-circleci install-metrics-server install-nginx-ingress

# All pods.
pods:
	@kubectl get pods --all-namespaces
