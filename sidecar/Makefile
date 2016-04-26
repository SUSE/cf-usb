NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

START_TARGET=darwin

CLIENTS=mysql mssql
TARGETS=CreateConnection CreateWorkspace GetConnection GetWorkspace DeleteConnection DeleteWorkspace

ifdef CLIENT
	CLIENTS=$(CLIENT)
endif

help:
	@echo "$(WARN_COLOR)##############################################################$(NO_COLOR)"
	@echo "$(WARN_COLOR)# $(OK_COLOR)build$(NO_COLOR) 		=> builds all clients and servers"
	@echo "$(WARN_COLOR)# $(OK_COLOR)buildserver$(NO_COLOR) 		=> builds all servers"
	@echo "$(WARN_COLOR)# $(OK_COLOR)buildclients$(NO_COLOR) 		=> builds all clients"
	@echo "$(WARN_COLOR)# $(OK_COLOR)clean$(NO_COLOR) 		=> deletes all the generated binaries"
	@echo "$(WARN_COLOR)# $(OK_COLOR)run$(NO_COLOR) CLIENT=client	=> starts the broker-server"
	@echo "$(WARN_COLOR)# available clients:$(NO_COLOR) $(CLIENTS)" 
	@echo "$(WARN_COLOR)# when any of the build commands have a parameter client"
	@echo "$(WARN_COLOR)# passed (i.e. CLIENT=mysql), only the binaries for the"
	@echo "$(WARN_COLOR)# respective client will be built$(NO_COLOR)"
	@echo "$(WARN_COLOR)##############################################################$(NO_COLOR)"

build: buildserver buildclients

clean: 
	@(if [ -d "./build" ]; then\
		rm -r ./build;\
	fi)
buildserver:
	@(for C in $(CLIENTS); do\
		$(call buildme,./server/cmd/broker-server,$$C);\
	done)

buildclients:
	@(for C in $(CLIENTS); do\
		for T in $(TARGETS); do\
			$(call buildme, ./clients/$$C/standalonecommands/$$T,$$C);\
		done;\
	done)

buildme =  (DIRNAME=$$(basename $(1));\
	echo "$(OK_COLOR)==> Building$(NO_COLOR) $$DIRNAME for $(2)";\
	DIRNAME=$$(basename $(1));\
	env GOOS="linux" go build -o build/linux/$(2)/$$DIRNAME $(1);\
	env GOOS="darwin" go build -o build/darwin/$(2)/$$DIRNAME $(1)) 

run: 
ifndef CLIENT
	@echo "$(WARN_COLOR)Please specify a CLIENT to run (i.e. make run CLIENT=mysql)$(NO_COLOR)"
else
ifndef START_TARGET
	@echo "$(WARN_COLOR)Please specify a target (i.e. make run CLIENT=client START_TARGET=lunux$(NO_COLOR)"
else
	@(cd ./build/$(START_TARGET)/$(CLIENT);\
	./broker-server --port 5333) 
endif
endif
