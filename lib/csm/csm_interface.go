package csm

//CSM is the model to use for implementing a new CSM client
type CSM interface {
	Login(string, string) error
	CreateWorkspace(string) (bool, error)
	WorkspaceExists(string) (bool, bool, error)
	DeleteWorkspace(string) (bool, error)
	CreateConnection(string, string) (interface{}, bool, error)
	ConnectionExists(string, string) (bool, bool, error)
	DeleteConnection(string, string) (bool, error)
	GetStatus() error
}
