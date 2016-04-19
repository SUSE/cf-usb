package main

import (
	"fmt"

	"github.com/hpcloud/cf-usb/sidecar/server/executablecaller"
)

func main() {

	var caller executablecaller.IExecutableCaller
	caller = executablecaller.DefaultCaller{}
	//response, _ := caller.CreateConnectionCaller("testdb2", "mvcv")
	//response, _ := caller.CreateWorkspaceCaller("testdb2")
	//response, _ := caller.DeleteWorkspaceCaller("testdb2")
	response, _ := caller.DeleteConnectionCaller("mvcv")
	//response, _ := caller.GetConnectionCaller("testdb2", "mvcv")
	//response, _ := caller.GetWorkspaceCaller("testdb2")
	fmt.Println(fmt.Sprintf("%s", response))

}
