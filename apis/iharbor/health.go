package iharbor

type SystemHealth struct {
	Status     string            `json:"status"`
	Components []ComponentStatus `json:"components"`
}

type ComponentStatus struct {
	Name   string `json:"name"`   //组件名
	Status string `json:"status"` //healthy
}
