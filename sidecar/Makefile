NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

CLIENTS=mysql mssql
TARGETS=CreateConnection CreateWorkspace GetConnection GetWorkspace DeleteConnection DeleteWorkspace

ifdef CLIENT
	CLIENTS=$(CLIENT)
endif

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
