// build包提供插件构建和打包相关功能
package build

// 导入所需的包
import (
	// 导入zip文件处理包
	"archive/zip"
	"path/filepath"

	// 导入错误处理包
	"errors"
	// 导入格式化包
	"fmt"
	// 导入工具包
	"github.com/1340691923/eve-plugin-sdk-go/util"
	// 导入IO操作包
	"io"
	// 导入日志包
	"log"
	// 导入操作系统包
	"os"
	// 导入命令执行包
	"os/exec"
	// 导入文件路径处理包
	// 导入运行时包
	"runtime"
	// 导入字符串处理包
	"strings"
)

// BuildPluginSvr 构建插件服务器
func BuildPluginSvr(pluginJsonFile string,isUpx bool) (err error) {
	// 调用内部构建函数
	err = buildPluginSvr(pluginJsonFile,isUpx)
	// 重置环境变量
	resetEnv()
	// 如果构建失败，打印错误信息
	if err != nil {
		log.Println("编译失败", err)
		return
	}

	return
}

// Build 定义构建结构体
type Build struct {
	// EVE版本
	EvVersion string
	// 主Go文件
	MainGoFile string
	// 插件别名
	PluginAlias string
	IsUpx bool
}

// BuildConfig 定义构建配置结构体
type BuildConfig struct {
	IsUpx bool
	// 操作系统
	OS string // GOOS
	// 架构
	GOARCH string
	// 环境变量
	Env map[string]string
	// 输出路径
	OutputPath string
	// EVE版本
	EvVersion string
	// 主Go文件
	MainGoFile string
	// 插件别名
	PluginAlias string
}

// getExecutableName 获取可执行文件名称
func getExecutableName(exname, os, arch string) (string, error) {
	// 根据操作系统和架构格式化可执行文件名
	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	// 如果是Windows系统，添加.exe后缀
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

// buildBackend 构建后端
func buildBackend(cfg BuildConfig) error {
	// 获取可执行文件名
	exeName, err := getExecutableName(cfg.PluginAlias, cfg.OS, cfg.GOARCH)
	if err != nil {
		return err
	}

	// 设置链接标志
	ldFlags := fmt.Sprintf("-w -s%s%s ", " ", `-extldflags "-static"`)

	// 如果是Windows系统，使用特定的链接标志
	if cfg.OS == "windows" {
		ldFlags = "-H windowsgui -w -s"
	}

	// 获取输出路径
	outputPath := cfg.OutputPath

	// 准备构建参数
	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	// 添加链接标志参数
	args = append(args, "-ldflags", ldFlags)

	// 添加主包路径
	rootPackage := cfg.MainGoFile

	args = append(args, rootPackage)

	// 设置环境变量
	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["GOARCH"] = cfg.GOARCH
	// 执行Go构建
	return RunGoBuild(cfg.IsUpx,cfg.Env, args...)
}

// newBuildConfig 创建新的构建配置
func newBuildConfig(os, arch, evVersion, mainGoFile, pluginAlias string,isUpx bool ) BuildConfig {
	// 返回初始化的构建配置
	return BuildConfig{
		OS:          os,
		GOARCH:      arch,
		EvVersion:   evVersion,
		OutputPath:  fmt.Sprintf("dist/%s", evVersion),
		Env:         map[string]string{},
		MainGoFile:  mainGoFile,
		PluginAlias: pluginAlias,
		IsUpx: isUpx,
	}
}

// LinuxArm64 构建Linux ARM64架构的插件
func (this *Build) LinuxArm64() error {
	// 使用newBuildConfig创建构建配置并执行buildBackend
	return buildBackend(newBuildConfig("linux", "arm64",
		this.EvVersion, this.MainGoFile, this.PluginAlias,this.IsUpx))
}

// LinuxAmd64 构建Linux AMD64架构的插件
func (this *Build) LinuxAmd64() error {
	// 使用newBuildConfig创建构建配置并执行buildBackend
	return buildBackend(newBuildConfig("linux", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias,this.IsUpx))
}

// WindowsAmd64 构建Windows AMD64架构的插件
func (this *Build) WindowsAmd64() error {
	// 使用newBuildConfig创建构建配置并执行buildBackend
	return buildBackend(newBuildConfig("windows", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias,this.IsUpx))
}

// DarwinAmd64 构建macOS AMD64架构的插件
func (this *Build) DarwinAmd64() error {
	// 使用newBuildConfig创建构建配置并执行buildBackend
	return buildBackend(newBuildConfig("darwin", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias,this.IsUpx))
}

// DarwinArm64 构建macOS ARM64架构的插件
func (this *Build) DarwinArm64() error {
	// 使用newBuildConfig创建构建配置并执行buildBackend
	return buildBackend(newBuildConfig("darwin", "arm64",
		this.EvVersion, this.MainGoFile, this.PluginAlias,this.IsUpx))
}

// resetEnv 重置环境变量
func resetEnv() {
	// 创建环境变量映射
	env := map[string]string{}
	// 设置GOOS为当前系统
	env["GOOS"] = runtime.GOOS
	// 禁用CGO
	env["CGO_ENABLED"] = "0"
	// 设置GOARCH为当前架构
	env["GOARCH"] = runtime.GOARCH
	// 应用环境变量设置
	setEnv(env)
}

// buildPluginSvr 构建插件服务器的内部实现
func buildPluginSvr(pluginJsonFile string,isUpx bool) (err error) {
	// 初始化插件数据结构
	var pluginData PluginJsonData
	// 从JSON文件加载插件配置
	err = util.LoadJsonAndParse(pluginJsonFile, &pluginData)
	if err != nil {
		return
	}
	// 验证版本号是否为空
	if pluginData.Version == "" {
		err = errors.New("版本号[version]不可为空")
		return
	}
	// 验证插件名称是否为空
	if pluginData.PluginName == "" {
		err = errors.New("插件名[plugin_name]不能为空")
		return
	}

	// 验证插件别名是否为空
	if pluginData.PluginAlias == "" {
		err = errors.New("插件别名[plugin_alias]不能为空")
	}

	// 验证主Go文件路径是否为空
	if pluginData.MainGoFile == "" {
		err = errors.New("main文件位置[main_go_file]不能为空")
		return
	}

	// 获取版本号
	version := pluginData.Version
	// 输出开始编译提示
	fmt.Println("开始编译插件二进制文件,请不要停止该操作，", pluginData.String())

	// 创建构建实例
	b := Build{
		EvVersion:   version,
		MainGoFile:  pluginData.MainGoFile,
		PluginAlias: pluginData.PluginAlias,
		IsUpx: isUpx,
	}

	// 创建错误函数运行器
	runErrFn := func(fnArr ...func() error) (err error) {
		// 依次执行函数数组中的每个函数
		for _, fn := range fnArr {
			// 如果执行出错，输出错误并返回
			if err = fn(); err != nil {
				fmt.Println("err", err)
				return
			}
		}
		return
	}

	// 执行各平台的构建函数
	err = runErrFn(
		b.LinuxArm64,
		b.LinuxAmd64,
		b.DarwinAmd64,
		b.DarwinArm64,
		b.WindowsAmd64,
	)
	if err != nil {
		return
	}

	// 设置输出路径
	outputPath := fmt.Sprintf("dist/%s", version)

	// 设置ZIP输出路径
	outputZipPath := fmt.Sprintf("dist/%s", strings.ReplaceAll(version, ".", "_"))
	// 开始检测是否已有打包
	fmt.Println("开始检测是否已有该版本打包")
	// 如果已存在打包，则删除旧包
	if CheckFileIsExist(outputZipPath) {
		fmt.Println("检测到已经该版本打包，正在删除老包")
		os.RemoveAll(outputZipPath)
	} else {
		fmt.Println("暂无该版本打包")
	}
	// 开始打包
	fmt.Println("开始打包")

	// 创建输出目录
	os.Mkdir(outputZipPath, 0755)

	// 设置ZIP文件路径
	outputZip := filepath.Join(outputZipPath, "plugin.zip")

	// 压缩文件夹
	err = CompressPathToZip(outputPath, "", "", outputZip)
	if err != nil {
		return
	}
	// 清理临时文件
	fmt.Println("打包成功，开始清理临时文件")
	os.RemoveAll(outputPath)
	fmt.Println("清理临时文件完毕")

	return
}

// setEnv 设置环境变量
func setEnv(env map[string]string) (err error) {
	// 如果环境变量不为空
	if len(env) > 0 {
		// 创建命令参数数组
		envArr := []string{"env", "-w"}
		// 格式化环境变量参数
		for k, v := range env {
			envArr = append(envArr, fmt.Sprintf("%s=%s", k, v))
		}

		// 输出构建命令
		fmt.Println(fmt.Sprintf("start build cmd: go %v", envArr))

		// 创建命令对象
		cmd := exec.Command("go", envArr...)

		// 执行命令
		err = cmd.Run()
		if err != nil {
			return
		}
	}
	return
}

// RunGoBuild 运行Go构建命令
func RunGoBuild(isUpx bool,env map[string]string, args ...string) (err error) {
	// 设置环境变量
	if err = setEnv(env); err != nil {
		return
	}

	// 输出构建命令
	fmt.Println(fmt.Sprintf("start build cmd: go %v", args))
	// 创建命令对象
	cmd := exec.Command("go", args...)
	// 执行命令
	err = cmd.Run()

	if err!=nil{
		return
	}

	if isUpx{
		fmt.Println(fmt.Sprintf("start cmd %s %s %s", "upx", "--best", args[2]))
		cmd = exec.Command("upx", "--best", args[2])
		err = cmd.Run()
	}

	return
}

// CheckFileIsExist 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	// 默认文件存在
	var exist = true
	// 检查文件状态，如果不存在则设置exist为false
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// CompressPathToZip 压缩文件夹到ZIP文件
func CompressPathToZip(path, excludeDir string, excludeFile string, targetFile string) error {
	// 创建目标文件
	d, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	// 关闭文件
	defer d.Close()
	// 创建ZIP写入器
	w := zip.NewWriter(d)
	// 关闭ZIP写入器
	defer w.Close()

	// 打开源路径
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	// 压缩文件夹
	err = compress(f, "", "", w, strings.Split(excludeDir, ","), excludeFile)

	return err
}

// compress 递归压缩文件/文件夹
func compress(file *os.File, prefix string, fileName string, zw *zip.Writer, excludeDirs []string, excludeFile string) error {
	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// 如果是目录
	if info.IsDir() {
		// 检查目录是否需要排除
		if fileName != "" {
			for _, excludeDir := range excludeDirs {
				if fileName == excludeDir {
					log.Println("excludeDir", excludeDir)
					return nil
				}
			}
		}

		// 更新路径前缀
		prefix = prefix + "/" + info.Name()
		// 读取目录内容
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		// 递归处理目录内的每个文件
		for _, fi := range fileInfos {
			// 打开文件
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}

			// 递归调用compress处理子文件
			err = compress(f, prefix, fi.Name(), zw, excludeDirs, excludeFile)
			if err != nil {
				return err
			}
		}
	} else {
		// 如果是文件，检查是否需要排除
		if excludeFile != "" && info.Name() == excludeFile {
			log.Println("excludeFile", info.Name())
			return nil
		}
		// 创建ZIP文件头
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		// 创建ZIP条目
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		// 复制文件内容到ZIP
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
