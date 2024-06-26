include secrets.mk
include konfig.mk
include standort.mk
include bezeichner.mk

kube-score: kube-score-konfig kube-score-standort kube-score-bezeichner

# Lint config.
lint: kube-linter kube-score

# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Delete namespace.
delete-namespace:
	kubectl delete namespaces $(namespace)

# Setup namespace.
setup-namespace:
	kubectl create namespace $(namespace)

# Rollout a namespace.
rollout-namespace:
	kubectl rollout restart deployment/$(namespace) -n $(namespace)

# Run kube-score for a namespace.
kube-score-namespace:
	kubectl api-resources --verbs=list --namespaced -o name \
		| xargs -I{} bash -c "kubectl get {} --namespace $(namespace) -oyaml && echo ---" \
		| kube-score score --ignore-test deployment-has-host-podantiaffinity  -

# Setup all apps.
setup: setup-konfig setup-standort setup-bezeichner

# Delete all apps.
delete: delete-konfig delete-standort delete-bezeichner

# Rollout all apps.
rollout: rollout-konfig rollout-standort rollout-bezeichner

# Verify all apps.
verify: verify-standort verify-bezeichner
