package constant

const (
	AesKey    = "1~$c31kjtR^@@c2#"
	AesVector = "#$3456$890A54321"
)
const (
	DeployStateUninitialized = "uninitialized"
	DeployStateInitialized   = "initialized"
	DeployStateRunning       = "running"
	DeployStateStop          = "stop"
	DeployStateComplete      = "complete"
)

const (
	ClusterStateDeploying   = "deploying"
	ClusterStateDeployError = "deploy-error"
	ClusterStateRunning     = "running"
	ClusterStateUnhealthy   = "unhealthy"
)
