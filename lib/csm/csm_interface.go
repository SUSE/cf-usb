package csm

//CSM is the model to use for implementing a new CSM client
type CSM interface {
	Login(string, string) error
	CreateWorkspace(string) error
	WorkspaceExists(string) (bool, error)
	DeleteWorkspace(string) error
	CreateConnection(string, string) (interface{}, error)
	ConnectionExists(string, string) (bool, error)
	DeleteConnection(string, string) error
	GetStatus() error
}
