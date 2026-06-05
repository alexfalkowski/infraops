# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Delete add-on namespaces.
delete:
	@kubectl delete namespaces nginx-ingress circleci metrics-server

# Set up Helm repos and install add-ons.
setup: setup-helm install-helm

# Set up Helm repos.
setup-helm:
	@helm repo add nginx https://kubernetes.github.io/ingress-nginx
	@helm repo add circleci https://circleci-public.github.io/cci-k8s-release-agent
	@helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
	@helm repo update

# Install CircleCI. This target is run manually from an operator workstation, not CI.
install-circleci:
	@helm upgrade --install circleci-release circleci/circleci-release-agent --set tokenSecret.token=$(CIRCLECI_K8S_TOKEN) --create-namespace --namespace circleci --set managedNamespaces="{lean}" --version 1.6.2

# Install metrics server.
install-metrics-server:
	@helm upgrade --install metrics-server-release metrics-server/metrics-server --create-namespace --namespace metrics-server --version 3.13.0

# Install nginx ingress controller.
install-nginx-ingress:
	@helm upgrade --install nginx-ingress-release nginx/ingress-nginx --namespace nginx-ingress --create-namespace --version 4.15.1 -f nginx/values.yaml

# Install Helm add-ons.
install-helm: install-circleci install-metrics-server install-nginx-ingress

# List all pods.
pods:
	@kubectl get pods --all-namespaces
