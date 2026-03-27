package main

func defaultVKubeFileSchema() VKubeFileSchema {
	return VKubeFileSchema{
		Type:     "object",
		Required: []string{"Kind", "vkubeToken", "containers"},
		Properties: map[string]PropertySchema{
			"Kind": {
				Type:        "string",
				Description: "固定为 vkube，不能修改",
				Required:    true,
				UserInput:   false,
				Fixed:       "vkube",
				Example:     "vkube",
			},
			"vkubeToken": {
				Type:        "string",
				Description: "用户 VKube 平台该服务的 token，用于指定特定的服务（service）",
				Required:    true,
				UserInput:   true,
				Example:     "your-vkube-token-here",
			},
			"imageRegistry": {
				Type:        "string",
				Description: "使用的镜像仓库前缀",
				Required:    false,
				UserInput:   false,
				Default:     "docker",
				Allowed:     []string{"docker", "ghcr"},
				Example:     "docker",
			},
			"containers": {
				Type:        "array",
				Description: "要部署的容器列表",
				Required:    true,
				UserInput:   true,
				Items: &PropertySchema{
					Type: "object",
					Properties: map[string]PropertySchema{
						"registryImagePath": {
							Type:        "string",
							Description: "从镜像仓库拉取的完整镜像路径(如果是 ghcr，则需要包含 ghcr.io 前缀)",
							Required:    false,
							UserInput:   true,
							Rules:       "与 imageName + tag 互斥",
							Example:     "ghcr.io/user/repo:tag",
						},
						"imageName": {
							Type:        "string",
							Description: "本地镜像名称",
							Required:    false,
							UserInput:   true,
							Rules:       "与 registryImagePath 互斥",
							Example:     "nginx",
						},
						"tag": {
							Type:        "string",
							Description: "镜像标签",
							Required:    false,
							UserInput:   true,
							Default:     "latest",
							Example:     "1.23",
						},
						"deploy": {
							Type:        "object",
							Description: "容器部署配置",
							Required:    true,
							UserInput:   true,
							Properties: map[string]PropertySchema{
								"containerName": {
									Type:        "string",
									Description: "容器运行时名称",
									Required:    true,
									UserInput:   true,
									Example:     "myapp",
								},
								"resourceUnit": {
									Type:         "integer",
									Description:  "资源单位，只允许填 1、2、4、8、16",
									Required:     true,
									UserInput:    true,
									AllowedRange: []int{1, 2, 4, 8, 16},
									Example:      2,
								},
								"ports": {
									Type:        "array",
									Description: "端口映射配置",
									Required:    false,
									UserInput:   true,
									Items: &PropertySchema{
										Type: "object",
										Properties: map[string]PropertySchema{
											"containerPort": {
												Type:        "integer",
												Description: "容器端口",
												Required:    true,
												UserInput:   true,
												Example:     80,
											},
											"hostPort": {
												Type:        "integer",
												Description: "宿主机端口",
												Required:    true,
												UserInput:   true,
												Example:     8080,
											},
										},
									},
								},
								"env": {
									Type:        "array",
									Description: "环境变量配置",
									Required:    false,
									UserInput:   true,
									Items: &PropertySchema{
										Type: "object",
										Properties: map[string]PropertySchema{
											"name": {
												Type:        "string",
												Description: "环境变量名",
												Required:    true,
												UserInput:   true,
												Example:     "MODE",
											},
											"value": {
												Type:        "string",
												Description: "环境变量值",
												Required:    true,
												UserInput:   true,
												Example:     "PROD",
											},
										},
									},
								},
								"command": {
									Type:        "array",
									Description: "覆盖容器的 ENTRYPOINT",
									Required:    false,
									UserInput:   true,
									Items: &PropertySchema{
										Type: "string",
									},
									Example: []string{"./start.sh"},
								},
								"args": {
									Type:        "array",
									Description: "容器启动参数",
									Required:    false,
									UserInput:   true,
									Items: &PropertySchema{
										Type: "string",
									},
								},
								"configurations": {
									Type:        "object",
									Description: "本地配置文件路径到容器内路径的映射",
									Required:    false,
									UserInput:   true,
								},
								"persistStorage": {
									Type:        "string",
									Description: "容器内需要持久化存储的路径",
									Required:    false,
									UserInput:   true,
									Example:     "/data",
								},
							},
						},
						"build": {
							Type:        "object",
							Description: "配置后可自动构建镜像（可选）",
							Required:    false,
							UserInput:   false,
							Properties: map[string]PropertySchema{
								"dockerfilePath": {
									Type:        "string",
									Description: "Dockerfile 路径",
									Required:    true,
									UserInput:   true,
									Example:     "./build/docker/Dockerfile",
								},
								"contextPath": {
									Type:        "string",
									Description: "构建上下文路径",
									Required:    false,
									UserInput:   true,
								},
								"buildArgs": {
									Type:        "array",
									Description: "构建参数",
									Required:    false,
									UserInput:   true,
									Items: &PropertySchema{
										Type: "object",
										Properties: map[string]PropertySchema{
											"name": {
												Type:      "string",
												Required:  true,
												UserInput: true,
											},
											"value": {
												Type:      "string",
												Required:  true,
												UserInput: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
