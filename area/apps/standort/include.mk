# Delete standort.
delete-standort:
	$(MAKE) namespace=standort delete-namespace

# Setup standort.
setup-standort:
	$(MAKE) namespace=standort setup-namespace

# Preview standort.
preview-standort:
	$(MAKE) namespace=standort preview-namespace

# Update standort.
update-standort:
	$(MAKE) namespace=standort update-namespace