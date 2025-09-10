package utils

import (
	"fmt"
	"os/exec"
)

// RunSwagInit 运行swag init命令生成swagger文档
func RunSwagInit() error {
	// 检查swag命令是否存在
	if _, err := exec.LookPath("swag"); err != nil {
		return fmt.Errorf("swag命令未找到，请先安装: go install github.com/swaggo/swag/cmd/swag@latest")
	}

	// 运行swag init命令
	cmd := exec.Command("swag", "init", "-g", "main.go", "-o", "docs")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("swag init执行失败: %v, 输出: %s", err, string(output))
	}

	fmt.Println("Swagger文档生成成功")
	return nil
}
