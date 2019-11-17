GITHUB := "github.com"
PROJECTNAME := $(shell basename "$(PWD)")
PACKAGENAME := $(GITHUB)/$(shell basename "$(shell dirname "$(PWD)")")/$(PROJECTNAME)
GOBASE := $(shell pwd)

.PHONY .SILENT: test
## test: Run coverage tests
test: 
	# Example: make test
	@echo " *** Running Coverage Tests ***"
	$(GOBASE)/test.sh
	@echo " *** Completed *** "

.PHONY .SILENT: relocate
## relocate: Relocate packages
relocate:
	# Example: make relocate TARGET=github.com/wingkwong/myproject
	@test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
	@echo " *** Relocating packages to $(TARGET) *** "
	$(eval ESCAPED_PACKAGENAME := $(shell echo "${PACKAGENAME}" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_PROJECTNAME := $(shell echo "${PROJECTNAME}" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_TARGET_PACKAGENAME := $(shell echo "${TARGET}" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_TARGET_PROJECTNAME := $(shell basename "$(shell dirname ${TARGET})" | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_PARENT_DIRECTORY:= $(shell cd ../ && pwd | sed -e 's/[\/&]/\\&/g'))
	$(eval ESCAPED_PARENT_DIRECTORYNAME:= $(shell basename $(ESCAPED_PARENT_DIRECTORY)))

	# Replacing ${ESCAPED_PACKAGENAME} to ${ESCAPED_TARGET_PACKAGENAME}
	@echo " *** Replacing ${ESCAPED_PACKAGENAME} to ${ESCAPED_TARGET_PACKAGENAME} *** "
	@grep -rlI '${PACKAGENAME}' --include=*.go ./ | xargs -I@ sed -i '' 's/${ESCAPED_PACKAGENAME}/${ESCAPED_TARGET_PACKAGENAME}/g' @
	# Replacing ${PROJECTNAME} to ${ESCAPED_TARGET_PROJECTNAME}
	@echo " *** Replacing ${PROJECTNAME} to ${ESCAPED_TARGET_PROJECTNAME} *** "
	@grep -rlI '${PROJECTNAME}' --include=*.go ./ | xargs -I@ sed -i '' 's/${ESCAPED_PROJECTNAME}/${ESCAPED_TARGET_PROJECTNAME}/g' @
	@echo " *** Completed *** "

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo