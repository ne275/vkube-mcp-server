package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const vkubeFileInteractiveGuide = `# VKubeFile Interactive Generation Guide

## 🎯 Your Task
You need to assist the user in generating a complete VKubeFile configuration. Follow these steps:

## 📋 Step 1: Understand the Schema
The JSON Schema below defines the complete structure of the VKubeFile. Carefully review each field's:
- **required**: Whether the field is mandatory
- **userInput**: Whether user input is needed
- **description**: Explanation of the field
- **example**: Example value
- **allowedRange** / **allowed**: Allowed value range

## 💬 Step 2: Collect Information Interactively
For all fields marked as **"userInput": true** and **"required": true**, you need to:

1. **Ask the user questions one by one**, using friendly natural language.
2. **Provide examples** to help the user understand.
3. **Validate user input** to ensure it meets the requirements (e.g., resourceUnit can only be 1/2/4/8/16).
4. **Infer intelligently**: If the user provides project code or description, try to extract information from it.

### Core Information to Collect:
- ✅ **vkubeToken**: VKube platform service token
- ✅ **containers[].imageName** or **registryImagePath**: Image information
  - **Ask the user whether the image is from a local build or a public registry.**
  - **If the image is local**, use "imageName" and "tag" fields to specify the image.
  - **If the image is from a public registry**, use "registryImagePath" to specify the full path (e.g., "nginx:latest").
- ✅ **containers[].deploy.containerName**: Container name
- ✅ **containers[].deploy.resourceUnit**: Resource unit (1/2/4/8/16)

### Optional but Recommended Information:
- 🔹 Port mappings (ports)
- 🔹 Environment variables (env)
- 🔹 Persistent storage paths (persistStorage)

## 🔧 Step 3: Generate the VKubeFile
Once all necessary information is collected, the AI should directly generate the complete VKubeFile content in YAML format and return it (do NOT call the GenerateVKubeFile API or any other endpoint). The generated output must be:
- Pure YAML text (no surrounding prose).
- Suitable to save with a ".yaml" extension (suggested filename: "vkubefile.yaml").
- Strictly follow the schema and examples collected from the user.

For example:
Kind: vkube
vkubeToken: "your-vkube-token"
containers:
	- imageName: "nginx:latest"
		deploy:
			containerName: "myapp"
			resourceUnit: 2
			ports:
				- containerPort: 80
					hostPort: 8080

## 💡 Tips
- If the user is unsure about a value, suggest a default.
- If the user mentions a Dockerfile or project repository, analyze this information to infer configurations.
- For complex configurations (e.g., multiple containers), configure one first, then ask if more need to be added.

---
Now, start interacting with the user! Begin by asking for the "vkubeToken" and, when all required fields are collected, directly output the final VKubeFile YAML.`

func registerTools(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_vkube_file_schema",
		Description: "Returns the VKubeFile JSON schema and an interactive guide for collecting user input to generate vkubefile.yaml",
	}, toolGetVKubeFileSchema)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "deploy_vkube_file",
		Description: "Builds a suggested vkube deploy -f <path> command for a local VKubeFile",
	}, toolDeployVKubeFile)
}

func toolGetVKubeFileSchema(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	schema := defaultVKubeFileSchema()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: vkubeFileInteractiveGuide},
			&mcp.TextContent{Text: toJSON(schema)},
		},
	}, nil, nil
}

type deployToolInput struct {
	VkubeFilePath string `json:"vkubeFilePath" jsonschema:"Local path to vkubefile.yaml (optional; default placeholder if empty)"`
}

func toolDeployVKubeFile(_ context.Context, _ *mcp.CallToolRequest, in deployToolInput) (*mcp.CallToolResult, any, error) {
	const defaultPath = "/path/to/vkubefile.yaml"
	path := in.VkubeFilePath
	if path == "" {
		path = defaultPath
	}
	cmd := fmt.Sprintf("vkube deploy -f %s", path)
	msg := "以下是生成的部署命令，请确认路径是否正确并决定是否执行。如果您有本地的 vkubefile.yaml 文件，请提供其完整路径以替换默认路径。"
	inst := "请在本地查找 vkubefile.yaml 文件（或您自定义的文件名），并将其路径提供给程序。例如：/home/user/vkubefile.yaml"
	out := fmt.Sprintf("%s\n\ncommand: %s\ndefaultPath: %s\n\n%s", msg, cmd, path, inst)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: out}},
	}, nil, nil
}
