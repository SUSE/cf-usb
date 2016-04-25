NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

build: buildserver buildclients

buildserver:
ifdef CLIENT
	$(call buildme,./server/cmd/broker-server)
else
	echo "$(ERROR_COLOR)No target clients specified$(NO_COLOR),\r\nPlease use CLIENT=clientname as a parameter"
endif

buildclients:
ifdef  CLIENT
	$(call buildme,./clients/$(CLIENT)/standalonecommands/CreateConnection)
	$(call buildme,./clients/$(CLIENT)/standalonecommands/CreateWorkspace)
	$(call buildme,./clients/$(CLIENT)/standalonecommands/DeleteConnection)
	$(call buildme,./clients/$(CLIENT)/standalonecommands/DeleteWorkspace)
	$(call buildme,./clients/$(CLIENT)/standalonecommands/GetConnection)
	$(call buildme,./clients/$(CLIENT)/standalonecommands/GetWorkspace)
else
	echo "$(ERROR_COLOR)No target clients specified$(NO_COLOR),\r\nPlease use CLIENT=clientname as a parameter"
endif

buildme =  @(DIRNAME=$$(basename $(1));\
	echo "$(OK_COLOR)==> Building$(NO_COLOR) $$DIRNAME"; \
	DIRNAME=$$(basename $(1)); \
	env GOOS="linux" go build -o build/linux/$(CLIENT)/$$DIRNAME $(1); \
	env GOOS="darwin" go build -o build/darwin/$(CLIENT)/$$DIRNAME $(1)) 
