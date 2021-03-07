package hooks

type DefaultRemoteHooks struct{}

func (drh *DefaultRemoteHooks) BeforeRemote(info RemoteHooksInfo) error {
	return nil
}

func (drh *DefaultRemoteHooks) AfterRemote(info RemoteHooksInfo) error {
	return nil
}

type DefaultOperationHooks struct{}

func (drh *DefaultOperationHooks) BeforeOperation(info OperationHooksInfo) error {
	return nil
}

func (drh *DefaultOperationHooks) AfterOperation(info OperationHooksInfo) error {
	return nil
}
