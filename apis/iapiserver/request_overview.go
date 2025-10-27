package iapiserver

type NodeOverview struct {
	// 节点名
	NodeName string `json:"node_name"`
	// 节点地址
	NodeAddr string `json:"node_addr"`
	// 总CPU
	CpuCapacity float64 `json:"cpu_capacity"`
	// 已使用CPU
	CpuUsed float64 `json:"cpu_used"`
	// CPU使用率
	CpuUsedRatio float64 `json:"cpu_used_ratio"`
	// CPU资源请求
	CpuResourceRequest int64 `json:"cpu_resource_request"`
	// CPU资源请求率
	CpuResourceRequestRatio float64 `json:"cpu_resource_request_ratio"`
	// CPU资源限制
	CpuResourceLimit int64 `json:"cpu_resource_limit"`
	// CPU资源限制率ListRequestG
	CpuResourceLimitRatio float64 `json:"cpu_resource_limit_ratio"`
	// 总内存
	MemoryCapacity float64 `json:"memory_capacity"`
	// 已使用内存
	MemoryUsed float64 `json:"memory_used"`
	// 内存使用率
	MemoryUsedRatio float64 `json:"memory_used_ratio"`
	// 内存资源请求
	MemoryResourceRequest int64 `json:"memory_resource_request"`
	// 内存资源请求率
	MemoryResourceRequestRatio float64 `json:"memory_resource_request_ratio"`
	// 内存资源限制
	MemoryResourceLimit int64 `json:"memory_resource_limit"`
	// 内存资源限制
	MemoryResourceLimitRatio float64 `json:"memory_resource_limit_ratio"`
	// 总存储容量
	DiskCapacity float64 `json:"disk_capacity"`
	// 已使用存储容量
	DiskUsed float64 `json:"disk_used"`
	// 存储容量使用率
	DiskUsedRatio float64 `json:"disk_used_ratio"`
	// Pod总数
	PodCapacity int `json:"pod_capacity"`
	// 已使用Pod总数
	PodUsed int `json:"pod_used"`
	// Pod使用率
	PodUsedRatio float64 `json:"pod_used_ratio"`
	// 是否网关节点
	IsGateway bool `json:"is_gateway"`
	// 运行状态
	Status string `json:"status"`
	// 节点角色
	Role string `json:"role"`
	// 系统镜像
	OsImage string `json:"os_image"`
	// 容器运行时
	ContainerRuntimeVersion string `json:"container_runtime_version"`
	// 错误信息
	ErrorMsg string `json:"error_msg"`
	// 创建时间
	CreateTime int64             `json:"create_time"`
	Extra      map[string]string `json:"extra"`

	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
}

type ComponentInfo struct {
	Namespace  string  `json:"namespace" description:"命名空间"`
	Name       string  `json:"name" description:"名称"`
	Addr       string  `json:"addr" description:"地址"`
	Status     string  `json:"status" description:"运行状态"` // 健康  不健康
	Image      string  `json:"image" description:"容器镜像"`
	CpuUsed    float64 `json:"cpu_used" description:"cpu使用率"`
	MemoryUsed float64 `json:"memory_used" description:"内存使用率"`
	DiskUsed   float64 `json:"disk_used" description:"磁盘使用率"`
}
