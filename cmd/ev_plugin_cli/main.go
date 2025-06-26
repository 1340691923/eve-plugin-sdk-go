package main

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ev_plugin_cli",
	Short: "ElasticView 插件开发工具",
	Long:  `ElasticView 插件开发工具，集成插件构建和打包功能`,
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "构建插件",
	Long:  `构建插件二进制文件，包括前端构建和后端编译`,
	RunE:  runBuildCmd,
}

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "打包源码",
	Long:  `将插件源码打包成zip文件`,
	RunE:  runZipCmd,
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "启动开发环境",
	Long:  `同时启动前端开发服务器和后端热编译服务，用于插件开发调试`,
	RunE:  runDevCmd,
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "检查开发环境",
	Long:  `检查插件开发所需的工具和环境是否正确安装和配置`,
	RunE:  runDoctorCmd,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化插件项目",
	Long:  `基于官方模板创建新的插件项目`,
	RunE:  runInitCmd,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装项目依赖",
	Long:  `安装插件项目的前后端依赖，包括 Go 模块和前端包`,
	RunE:  runInstallCmd,
}

// 全局参数变量
var (
	pluginJsonFile string
	execUpx        bool
	destZip        string
	excludeDir     string
	pluginName     string
	pluginAlias    string
)

func init() {
	// build 命令参数
	buildCmd.Flags().StringVarP(&pluginJsonFile, "plugin-json", "p", "plugin.json", "插件配置文件")
	buildCmd.Flags().BoolVarP(&execUpx, "upx", "u", false, "是否使用upx压缩")

	// zip 命令参数
	zipCmd.Flags().StringVarP(&destZip, "output", "o", "dist/source_code.zip", "输出的zip文件路径")
	zipCmd.Flags().StringVarP(&excludeDir, "exclude", "e", "node_modules,dist,.idea,.vscode,.git", "要排除的文件夹路径")

	// init 命令参数
	initCmd.Flags().StringVarP(&pluginName, "name", "n", "", "插件中文名称（必填）")
	initCmd.Flags().StringVarP(&pluginAlias, "alias", "a", "", "插件别名/项目名（必填）")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("alias")

	// 添加子命令
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(zipCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令时出错: %v\n", err)
		os.Exit(1)
	}
}

// runBuildCmd 执行构建命令
func runBuildCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 开始构建插件...")

	// 构建前端
	if err := buildVue(); err != nil {
		return fmt.Errorf("前端构建失败: %w", err)
	}

	// 构建后端
	if err := build.BuildPluginSvr(pluginJsonFile, execUpx); err != nil {
		return fmt.Errorf("后端构建失败: %w", err)
	}

	fmt.Println("✅ 插件构建成功!")
	return nil
}

// runZipCmd 执行打包命令
func runZipCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("📦 开始打包源码...")

	sourceDir, _ := os.Getwd()
	execFileName := "plugin"
	if runtime.GOOS == "windows" {
		execFileName = execFileName + ".exe"
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(destZip), os.ModePerm); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 执行打包
	if err := build.CompressPathToZip(sourceDir, excludeDir, execFileName, destZip); err != nil {
		return fmt.Errorf("打包失败: %w", err)
	}

	fmt.Printf("✅ 打包完成，输出文件: %s\n", destZip)
	return nil
}

// runDevCmd 执行开发环境启动命令
func runDevCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 启动插件开发环境...")
	fmt.Println("📋 这将同时启动前端开发服务器和后端热编译服务")
	fmt.Println("---")

	// 检查前端目录是否存在
	frontendDir := filepath.Join("frontend")
	hasFrontend := true
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		fmt.Println("⚠️  未找到 frontend 目录，将只启动后端服务")
		hasFrontend = false
	}

	// 检查 package.json 是否存在
	if hasFrontend {
		packageJson := filepath.Join(frontendDir, "package.json")
		if _, err := os.Stat(packageJson); os.IsNotExist(err) {
			fmt.Println("⚠️  未找到 package.json，将只启动后端服务")
			hasFrontend = false
		}
	}

	// 检查是否安装了 gowatch
	if _, err := exec.LookPath("gowatch"); err != nil {
		return fmt.Errorf("❌ 未找到 gowatch 工具，请先安装: go install github.com/silenceper/gowatch@latest")
	}

	// 创建上下文和信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	var frontendCmd, backendCmd *exec.Cmd

	// 启动前端开发服务器（如果存在）
	if hasFrontend {
		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Println("🎨 [FRONTEND] 启动前端开发服务器...")

			frontendCmd = exec.CommandContext(ctx, "npm", "run", "dev")
			frontendCmd.Dir = frontendDir

			// 创建带前缀的输出
			stdout, err := frontendCmd.StdoutPipe()
			if err != nil {
				fmt.Printf("❌ [FRONTEND] 无法创建输出管道: %v\n", err)
				return
			}

			stderr, err := frontendCmd.StderrPipe()
			if err != nil {
				fmt.Printf("❌ [FRONTEND] 无法创建错误输出管道: %v\n", err)
				return
			}

			if err := frontendCmd.Start(); err != nil {
				fmt.Printf("❌ [FRONTEND] 启动失败: %v\n", err)
				return
			}

			fmt.Println("✅ [FRONTEND] 前端开发服务器已启动")

			// 处理输出
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					fmt.Printf("📦 [FRONTEND] %s\n", scanner.Text())
				}
			}()

			go func() {
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					fmt.Printf("🔥 [FRONTEND] %s\n", scanner.Text())
				}
			}()

			if err := frontendCmd.Wait(); err != nil {
				if ctx.Err() == nil { // 只有在非主动取消时才报错
					fmt.Printf("❌ [FRONTEND] 进程退出: %v\n", err)
				}
			}
		}()

		// 等待前端服务器启动
		time.Sleep(2 * time.Second)
	}

	// 启动后端热编译服务
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("⚙️  [BACKEND] 启动后端热编译服务...")

		backendCmd = exec.CommandContext(ctx, "gowatch")

		// 创建带前缀的输出
		stdout, err := backendCmd.StdoutPipe()
		if err != nil {
			fmt.Printf("❌ [BACKEND] 无法创建输出管道: %v\n", err)
			return
		}

		stderr, err := backendCmd.StderrPipe()
		if err != nil {
			fmt.Printf("❌ [BACKEND] 无法创建错误输出管道: %v\n", err)
			return
		}

		if err := backendCmd.Start(); err != nil {
			fmt.Printf("❌ [BACKEND] 启动失败: %v\n", err)
			return
		}

		fmt.Println("✅ [BACKEND] 后端热编译服务已启动")

		// 处理输出
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Printf("🔧 [BACKEND] %s\n", scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Printf("🚨 [BACKEND] %s\n", scanner.Text())
			}
		}()

		if err := backendCmd.Wait(); err != nil {
			if ctx.Err() == nil { // 只有在非主动取消时才报错
				fmt.Printf("❌ [BACKEND] 进程退出: %v\n", err)
			}
		}
	}()

	fmt.Println("---")
	fmt.Println("🎯 开发环境已启动！")
	if hasFrontend {
		fmt.Println("   前端服务: http://localhost:3000 (或查看前端输出)")
	}
	fmt.Println("   后端服务: 热编译模式")
	fmt.Println("📝 使用 Ctrl+C 停止所有服务")
	fmt.Println("---")

	// 等待信号或所有服务结束
	select {
	case <-sigChan:
		fmt.Println("\n🛑 收到中断信号，正在关闭服务...")
		cancel()

		// 给进程一些时间优雅关闭
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			fmt.Println("✅ 所有服务已优雅关闭")
		case <-time.After(5 * time.Second):
			fmt.Println("⚠️  强制关闭服务...")
			if frontendCmd != nil && frontendCmd.Process != nil {
				frontendCmd.Process.Kill()
			}
			if backendCmd != nil && backendCmd.Process != nil {
				backendCmd.Process.Kill()
			}
		}
	case <-ctx.Done():
		wg.Wait()
	}

	return nil
}

// runDoctorCmd 执行环境检查命令
func runDoctorCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("🏥 ElasticView 插件开发环境检查")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	allPassed := true

	// 定义检查项
	checks := []struct {
		name        string
		command     string
		args        []string
		required    bool
		description string
		installCmd  string
	}{
		{
			name:        "Go",
			command:     "go",
			args:        []string{"version"},
			required:    true,
			description: "Go 编程语言运行时",
			installCmd:  "https://golang.org/dl/",
		},
		{
			name:        "Node.js",
			command:     "node",
			args:        []string{"--version"},
			required:    true,
			description: "Node.js JavaScript 运行时",
			installCmd:  "https://nodejs.org/",
		},
		{
			name:        "npm",
			command:     "npm",
			args:        []string{"--version"},
			required:    true,
			description: "Node.js 包管理器",
			installCmd:  "随 Node.js 一起安装",
		},
		{
			name:        "pnpm",
			command:     "pnpm",
			args:        []string{"--version"},
			required:    false,
			description: "高性能的 Node.js 包管理器",
			installCmd:  "npm install -g pnpm",
		},
		{
			name:        "gowatch",
			command:     "gowatch",
			args:        []string{"-v"},
			required:    true,
			description: "Go 代码热重载工具",
			installCmd:  "go install github.com/silenceper/gowatch@latest",
		},
		{
			name:        "git",
			command:     "git",
			args:        []string{"--version"},
			required:    true,
			description: "版本控制系统",
			installCmd:  "https://git-scm.com/downloads",
		},
	}

	// 执行检查
	for _, check := range checks {
		passed := checkTool(check.name, check.command, check.args, check.required, check.description, check.installCmd)
		if !passed && check.required {
			allPassed = false
		}
	}

	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 50))

	if allPassed {
		fmt.Println("✅ 恭喜！您的开发环境配置完整，可以开始插件开发了！")
		fmt.Println()
		fmt.Println("🚀 下一步:")
		fmt.Println(" ev_plugin_cli dev")
	} else {
		fmt.Println("❌ 发现一些问题，请根据上述提示安装缺失的工具")
		fmt.Println()
		fmt.Println("📚 更多帮助:")
		fmt.Println("   - 开发环境搭建: http://www.elastic-view.cn/plugin-dev/setup.html")
		fmt.Println("   - 开发流程: http://www.elastic-view.cn/plugin-dev/workflow.html")
		return fmt.Errorf("环境检查未通过")
	}

	return nil
}

// runInitCmd 执行初始化插件项目命令
func runInitCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 初始化 ElasticView 插件项目...")
	fmt.Printf("📋 插件名称: %s\n", pluginName)
	fmt.Printf("📋 插件别名: %s\n", pluginAlias)
	fmt.Println("---")

	// 检查目标目录是否已存在
	if _, err := os.Stat(pluginAlias); !os.IsNotExist(err) {
		return fmt.Errorf("❌ 目录 '%s' 已存在，请选择其他名称或删除现有目录", pluginAlias)
	}

	// 下载并解压模板
	templateURL := "https://github.com/1340691923/eve-plugin-vue3-template/archive/refs/tags/latest.zip"
	fmt.Printf("📥 下载模板压缩包: %s\n", templateURL)

	if err := downloadAndExtractTemplate(templateURL, pluginAlias); err != nil {
		return fmt.Errorf("❌ 下载模板失败: %w", err)
	}

	// 项目目录
	projectDir := filepath.Join(".", pluginAlias)

	// 修改 plugin.json 文件
	pluginJsonPath := filepath.Join(projectDir, "plugin.json")
	if err := updatePluginJson(pluginJsonPath, pluginName, pluginAlias); err != nil {
		return fmt.Errorf("❌ 更新 plugin.json 失败: %w", err)
	}

	fmt.Println("---")
	fmt.Println("✅ 插件项目初始化完成!")
	fmt.Printf("📁 项目目录: %s\n", pluginAlias)
	fmt.Println()
	fmt.Println("🚀 下一步:")
	fmt.Printf("   1. 进入项目目录: cd %s\n", pluginAlias)
	fmt.Println("   2. 安装项目依赖: ev_plugin_cli install")
	fmt.Println("   3. 启动开发环境: ev_plugin_cli dev")
	fmt.Println()
	fmt.Println("📚 更多帮助:")
	fmt.Println("   - 开发流程: http://www.elastic-view.cn/plugin-dev/workflow.html")
	fmt.Println("   - API 文档: http://www.elastic-view.cn/plugin-dev/api.html")

	return nil
}

// runInstallCmd 执行安装依赖命令
func runInstallCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("📦 开始安装插件项目依赖...")
	fmt.Println("---")

	// 检查当前目录是否是插件项目
	if !isPluginProject() {
		return fmt.Errorf("❌ 当前目录不是有效的插件项目，请确保存在 plugin.json 和 go.mod 文件")
	}

	// 检查必需工具
	if err := checkRequiredTools(); err != nil {
		return err
	}

	// 安装 Go 依赖
	fmt.Println("🔧 安装 Go 依赖...")
	if err := installGoDependencies(); err != nil {
		return fmt.Errorf("❌ 安装 Go 依赖失败: %w", err)
	}

	// 检查是否有前端项目
	frontendDir := "frontend"
	if _, err := os.Stat(frontendDir); err == nil {
		// 检查 package.json 是否存在
		packageJsonPath := filepath.Join(frontendDir, "package.json")
		if _, err := os.Stat(packageJsonPath); err == nil {
			fmt.Println("🎨 安装前端依赖...")
			if err := installFrontendDependencies(frontendDir); err != nil {
				return fmt.Errorf("❌ 安装前端依赖失败: %w", err)
			}
		} else {
			fmt.Println("⚠️  未找到 frontend/package.json，跳过前端依赖安装")
		}
	} else {
		fmt.Println("⚠️  未找到 frontend 目录，跳过前端依赖安装")
	}

	fmt.Println("---")
	fmt.Println("✅ 依赖安装完成!")
	fmt.Println()
	fmt.Println("🚀 下一步:")
	fmt.Println("   启动开发环境: ev_plugin_cli dev")
	fmt.Println()
	fmt.Println("📚 更多帮助:")
	fmt.Println("   - 开发流程: http://www.elastic-view.cn/plugin-dev/workflow.html")

	return nil
}

// isPluginProject 检查当前目录是否是插件项目
func isPluginProject() bool {
	// 检查 plugin.json 和 go.mod 是否存在
	if _, err := os.Stat("plugin.json"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return false
	}
	return true
}

// checkRequiredTools 检查必需工具
func checkRequiredTools() error {
	// 检查 Go
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("❌ 未找到 go 命令，请先安装 Go: https://golang.org/dl/")
	}

	// 检查 Node.js
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("❌ 未找到 node 命令，请先安装 Node.js: https://nodejs.org/")
	}

	// 检查 npm
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("❌ 未找到 npm 命令，请先安装 npm（通常随 Node.js 一起安装）")
	}

	fmt.Println("✅ 必需工具检查通过")
	return nil
}

// installGoDependencies 安装 Go 依赖
func installGoDependencies() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("✅ Go 依赖安装完成")
	return nil
}

// installFrontendDependencies 安装前端依赖
func installFrontendDependencies(frontendDir string) error {
	// 优先使用 pnpm，如果不存在则使用 npm
	var cmd *exec.Cmd
	if _, err := exec.LookPath("pnpm"); err == nil {
		fmt.Println("   使用 pnpm 安装前端依赖...")
		cmd = exec.Command("pnpm", "install")
	} else {
		fmt.Println("   使用 npm 安装前端依赖...")
		cmd = exec.Command("npm", "install")
	}

	cmd.Dir = frontendDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("✅ 前端依赖安装完成")
	return nil
}

// updatePluginJson 更新 plugin.json 文件
func updatePluginJson(filePath, pluginName, pluginAlias string) error {
	fmt.Printf("📝 更新 plugin.json 文件...\n")

	// 读取原文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析 JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 更新配置
	config["plugin_name"] = pluginName
	config["plugin_alias"] = pluginAlias
	config["developer"] = "开发者名称" // 可以后续手动修改

	// 更新前端路由中的路径
	if routes, exists := config["frontend_routes"].([]interface{}); exists {
		for _, route := range routes {
			if routeMap, ok := route.(map[string]interface{}); ok {
				// 更新路径前缀
				if path, exists := routeMap["path"].(string); exists {
					// 保持原有路径，只是确保配置正确
					routeMap["path"] = path
				}
			}
		}
	}

	// 重新编码 JSON
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("编码 JSON 失败: %w", err)
	}

	// 写回文件
	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("   ✅ plugin_name: %s\n", pluginName)
	fmt.Printf("   ✅ plugin_alias: %s\n", pluginAlias)
	return nil
}

// downloadAndExtractTemplate 下载并解压模板
func downloadAndExtractTemplate(templateURL, projectName string) error {
	// 创建临时文件
	tempFile := "template.zip"
	defer os.Remove(tempFile)

	// 下载文件
	fmt.Println("📥 正在下载模板文件...")
	resp, err := http.Get(templateURL)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 创建临时文件
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Println("📦 正在解压模板文件...")

	// 解压文件
	if err := extractZip(tempFile, projectName); err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}

	fmt.Println("✅ 模板下载并解压完成")
	return nil
}

// extractZip 解压zip文件
func extractZip(src, dest string) error {
	// 打开zip文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// 提取文件
	for _, f := range r.File {
		// 跳过根目录（通常是仓库名）
		parts := strings.Split(f.Name, "/")
		if len(parts) <= 1 {
			continue
		}

		// 构建目标路径（去掉第一级目录）
		relativePath := strings.Join(parts[1:], "/")
		if relativePath == "" {
			continue
		}

		destPath := filepath.Join(dest, relativePath)

		// 确保目标目录存在
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		// 提取文件
		if err := extractFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractFile 提取单个文件
func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// checkTool 检查单个工具是否安装
func checkTool(name, command string, args []string, required bool, description, installCmd string) bool {
	fmt.Printf("🔍 检查 %s...", name)

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if required {
			fmt.Printf(" ❌ 未安装\n")
			fmt.Printf("   描述: %s\n", description)
			fmt.Printf("   安装: %s\n", installCmd)
		} else {
			fmt.Printf(" ⚠️  未安装 (可选)\n")
			fmt.Printf("   描述: %s\n", description)
			if strings.Contains(installCmd, "跳过安装") {
				fmt.Printf("   安装: %s\n", installCmd)
			} else {
				fmt.Printf("   安装: %s (可跳过)\n", installCmd)
			}
		}
		return false
	}

	// 提取版本信息
	version := strings.TrimSpace(string(output))
	if len(version) > 100 {
		version = version[:100] + "..."
	}

	fmt.Printf(" ✅ 已安装\n")
	fmt.Printf("   版本: %s\n", version)
	return true
}

// checkPluginJson 检查 plugin.json 配置
func checkPluginJson() bool {
	data, err := os.ReadFile("plugin.json")
	if err != nil {
		return false
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return false
	}

	// 检查必需字段
	requiredFields := []string{"plugin_alias", "plugin_name", "version", "main_go_file"}
	for _, field := range requiredFields {
		if _, exists := config[field]; !exists {
			fmt.Printf("\n   ⚠️  缺少必需字段: %s", field)
			return false
		}
	}

	// 检查 frontend_debug 配置
	if frontendDebug, exists := config["frontend_debug"]; exists {
		if frontendDebug == true {
			fmt.Printf("\n   💡 提示: frontend_debug 为 true，适合开发调试")
		}
	}

	return true
}

// buildVue 构建前端项目
func buildVue() error {
	frontendDir := filepath.Join("frontend")

	// 检查前端目录是否存在
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		fmt.Println("⚠️  未找到 frontend 目录，跳过前端构建")
		return nil
	}

	// 检查 package.json 是否存在
	packageJson := filepath.Join(frontendDir, "package.json")
	if _, err := os.Stat(packageJson); os.IsNotExist(err) {
		fmt.Println("⚠️  未找到 package.json，跳过前端构建")
		return nil
	}

	fmt.Println("=================构建前端================")

	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = frontendDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithStack(err)
	}

	if err = cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err = cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
