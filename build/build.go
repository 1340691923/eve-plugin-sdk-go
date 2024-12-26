package build

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func BuildPluginSvr(pluginJsonFile string) (err error) {

	err = buildPluginSvr(pluginJsonFile)
	resetEnv()
	if err != nil {
		log.Println("编译失败", err)
		return
	}

	return
}

type Build struct {
	EvVersion   string
	MainGoFile  string
	PluginAlias string
}

type BuildConfig struct {
	OS          string // GOOS
	GOARCH      string
	Env         map[string]string
	OutputPath  string
	EvVersion   string
	MainGoFile  string
	PluginAlias string
}

func getExecutableName(exname, os, arch string) (string, error) {

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

func buildBackend(cfg BuildConfig) error {

	exeName, err := getExecutableName(cfg.PluginAlias, cfg.OS, cfg.GOARCH)
	if err != nil {
		return err
	}

	ldFlags := fmt.Sprintf("-w -s%s%s ", " ", `-extldflags "-static"`)

	outputPath := cfg.OutputPath

	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	args = append(args, "-ldflags", ldFlags)

	rootPackage := cfg.MainGoFile

	args = append(args, rootPackage)

	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["CGO_ENABLED"] = "0"
	cfg.Env["GOARCH"] = cfg.GOARCH
	return RunGoBuild(cfg.Env, args...)
}

func newBuildConfig(os, arch, evVersion, mainGoFile, pluginAlias string) BuildConfig {
	return BuildConfig{
		OS:          os,
		GOARCH:      arch,
		EvVersion:   evVersion,
		OutputPath:  fmt.Sprintf("dist/%s", evVersion),
		Env:         map[string]string{},
		MainGoFile:  mainGoFile,
		PluginAlias: pluginAlias,
	}
}

func (this *Build) LinuxArm64() error {
	return buildBackend(newBuildConfig("linux", "arm64",
		this.EvVersion, this.MainGoFile, this.PluginAlias))
}

func (this *Build) LinuxAmd64() error {
	return buildBackend(newBuildConfig("linux", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias))
}

func (this *Build) WindowsAmd64() error {
	return buildBackend(newBuildConfig("windows", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias))
}

func (this *Build) DarwinAmd64() error {
	return buildBackend(newBuildConfig("darwin", "amd64",
		this.EvVersion, this.MainGoFile, this.PluginAlias))
}

func (this *Build) DarwinArm64() error {
	return buildBackend(newBuildConfig("darwin", "arm64",
		this.EvVersion, this.MainGoFile, this.PluginAlias))
}


func resetEnv() {
	env := map[string]string{}
	env["GOOS"] = runtime.GOOS
	env["CGO_ENABLED"] = "0"
	env["GOARCH"] = runtime.GOARCH
	setEnv(env)
}

func buildPluginSvr(pluginJsonFile string) (err error) {
	var pluginData PluginJsonData
	err = util.LoadJsonAndParse(pluginJsonFile, &pluginData)
	if err != nil {
		return
	}
	if pluginData.Version == "" {
		err = errors.New("版本号[version]不可为空")
		return
	}
	if pluginData.PluginName == "" {
		err = errors.New("插件名[plugin_name]不能为空")
		return
	}

	if pluginData.PluginAlias == "" {
		err = errors.New("插件别名[plugin_alias]不能为空")
	}

	if pluginData.MainGoFile == "" {
		err = errors.New("main文件位置[main_go_file]不能为空")
		return
	}

	version := pluginData.Version
	fmt.Println("开始编译插件二进制文件,请不要停止该操作，", pluginData.String())

	b := Build{
		EvVersion:   version,
		MainGoFile:  pluginData.MainGoFile,
		PluginAlias: pluginData.PluginAlias,
	}

	runErrFn := func(fnArr ...func() error) (err error) {
		for _, fn := range fnArr {
			if err = fn(); err != nil {
				fmt.Println("err", err)
				return
			}
		}
		return
	}

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

	outputPath := fmt.Sprintf("dist/%s", version)

	outputZipPath := fmt.Sprintf("dist/%s", strings.ReplaceAll(version, ".", "_"))
	fmt.Println("开始检测是否已有该版本打包")
	if CheckFileIsExist(outputZipPath) {
		fmt.Println("检测到已经该版本打包，正在删除老包")
		os.RemoveAll(outputZipPath)
	} else {
		fmt.Println("暂无该版本打包")
	}
	fmt.Println("开始打包")

	os.Mkdir(outputZipPath, 0755)

	outputZip := filepath.Join(outputZipPath, "plugin.zip")

	err = CompressPathToZip(outputPath, "", "", outputZip)
	if err != nil {
		return
	}
	fmt.Println("打包成功，开始清理临时文件")
	os.RemoveAll(outputPath)
	fmt.Println("清理临时文件完毕")

	return
}

func setEnv(env map[string]string) (err error) {
	if len(env) > 0 {
		envArr := []string{"env", "-w"}
		for k, v := range env {
			envArr = append(envArr, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Println(fmt.Sprintf("start build cmd: go %v", envArr))

		cmd := exec.Command("go", envArr...)

		err = cmd.Run()
		if err != nil {
			return
		}
	}
	return
}

func RunGoBuild(env map[string]string, args ...string) (err error) {
	if err = setEnv(env); err != nil {
		return
	}

	fmt.Println(fmt.Sprintf("start build cmd: go %v", args))
	cmd := exec.Command("go", args...)
	err = cmd.Run()
	return
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// CompressPathToZip 压缩文件夹
func CompressPathToZip(path, excludeDir string, excludeFile string, targetFile string) error {
	d, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	err = compress(f, "", "", w, strings.Split(excludeDir, ","), excludeFile)

	return err
}

func compress(file *os.File, prefix string, fileName string, zw *zip.Writer, excludeDirs []string, excludeFile string) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if fileName != "" {
			for _, excludeDir := range excludeDirs {
				if fileName == excludeDir {
					log.Println("excludeDir", excludeDir)
					return nil
				}
			}
		}

		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}

			err = compress(f, prefix, fi.Name(), zw, excludeDirs, excludeFile)
			if err != nil {
				return err
			}
		}
	} else {

		if excludeFile != "" && info.Name() == excludeFile {
			log.Println("excludeFile", info.Name())
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
