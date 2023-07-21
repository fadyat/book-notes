help:
	@perl -e '\
        my %help; \
        while (<>) { \
            if (/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/) { \
                my ($$option_name, $$group, $$description) = ($$1, $$2 // "options", $$3); \
                push @{$$help{$$group}}, [$$option_name, $$description]; \
            } \
        } \
        foreach my $$group (keys %help) { \
            print "$$group:\n"; \
            foreach my $$option_desc (@{$$help{$$group}}) { \
                my ($$option, $$desc) = @{$$option_desc}; \
                printf "  %-20s %s\n", $$option, $$desc; \
            } \
            print "\n"; \
        }' $(MAKEFILE_LIST)

