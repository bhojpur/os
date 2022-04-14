TARGETS := $(shell ls scripts)

.dapper:
	@echo Downloading Rancher Dapper
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-`uname -s`-`uname -m` > .dapper.tmp
	@@chmod +x .dapper.tmp
	@./.dapper.tmp -v
	@mv .dapper.tmp .dapper

$(TARGETS): .dapper
	@rm -rf ./dist ./build
	./.dapper $@

.DEFAULT_GOAL := default

.PHONY: $(TARGETS)