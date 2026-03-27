package main

// VKubeFileSchema 定义 VKubeFile 的结构化 schema
type VKubeFileSchema struct {
	Type       string                    `json:"type"`
	Required   []string                  `json:"required"`
	Properties map[string]PropertySchema `json:"properties"`
}

type PropertySchema struct {
	Type         string                    `json:"type"`
	Description  string                    `json:"description"`
	Required     bool                      `json:"required,omitempty"`
	UserInput    bool                      `json:"userInput,omitempty"`
	Default      interface{}               `json:"default,omitempty"`
	Fixed        string                    `json:"fixed,omitempty"`
	Allowed      []string                  `json:"allowed,omitempty"`
	AllowedRange []int                     `json:"allowedRange,omitempty"`
	Example      interface{}               `json:"example,omitempty"`
	Items        *PropertySchema           `json:"items,omitempty"`
	Properties   map[string]PropertySchema `json:"properties,omitempty"`
	Rules        string                    `json:"rules,omitempty"`
}

// GenerateVKubeFileRequest 定义生成 VKubeFile 的请求结构（供文档/扩展使用）
type GenerateVKubeFileRequest struct {
	Kind          string      `json:"Kind" binding:"required"`
	VkubeToken    string      `json:"vkubeToken" binding:"required"`
	ImageRegistry string      `json:"imageRegistry"`
	Containers    []Container `json:"containers" binding:"required"`
}

type Container struct {
	RegistryImagePath string `json:"registryImagePath,omitempty"`
	ImageName         string `json:"imageName,omitempty"`
	Tag               string `json:"tag,omitempty"`
	Deploy            Deploy `json:"deploy" binding:"required"`
	Build             *Build `json:"build,omitempty"`
}

type Deploy struct {
	ContainerName  string            `json:"containerName" binding:"required"`
	ResourceUnit   int               `json:"resourceUnit" binding:"required"`
	Ports          []Port            `json:"ports,omitempty"`
	Env            []EnvVar          `json:"env,omitempty"`
	Command        []string          `json:"command,omitempty"`
	Args           []string          `json:"args,omitempty"`
	Configurations map[string]string `json:"configurations,omitempty"`
	PersistStorage string            `json:"persistStorage,omitempty"`
}

type Port struct {
	ContainerPort int `json:"containerPort" binding:"required"`
	HostPort      int `json:"hostPort" binding:"required"`
}

type EnvVar struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Build struct {
	DockerfilePath string     `json:"dockerfilePath" binding:"required"`
	ContextPath    string     `json:"contextPath,omitempty"`
	BuildArgs      []BuildArg `json:"buildArgs,omitempty"`
}

type BuildArg struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}
