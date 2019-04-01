package rest

import (
	"fmt"
	"go/build"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

//判断文件或文件夹是否已经存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取GOPATH路径
func GetGoPath() string {

	return build.Default.GOPATH

}

//获取该应用可执行文件的位置。
//例如：C:\Users\lishuang\AppData\Local\Temp
func GetHomePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	//如果exPath中包含了 /private/var/folders 我们认为是在Mac的开发环境中
	macDev := strings.HasPrefix(exPath, "/private/var/folders")
	if macDev {
		exPath = GetGoPath() + "/src/tank/tmp"
	}

	//如果exPath中包含了 \\AppData\\Local\\Temp 我们认为是在Win的开发环境中
	systemUser, err := user.Current()
	winDev := strings.HasPrefix(exPath, systemUser.HomeDir+"\\AppData\\Local\\Temp")
	if winDev {
		exPath = GetGoPath() + "/src/tank/tmp"
	}

	return exPath
}

//获取前端静态资源的位置。如果你在开发模式下，可以将这里直接返回tank/build下面的html路径。
//例如：C:/Users/lishuang/AppData/Local/Temp/html
func GetHtmlPath() string {

	homePath := GetHomePath()
	filePath := homePath + "/html"
	exists, err := PathExists(filePath)
	if err != nil {
		panic("判断上传文件是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			panic("创建上传文件夹时出错！")
		}
	}

	return filePath
}

//如果文件夹存在就不管，不存在就创建。 例如：/var/www/matter
func MakeDirAll(dirPath string) string {

	exists, err := PathExists(dirPath)
	if err != nil {
		panic("判断文件是否存在时出错！")
	}
	if !exists {
		//TODO:文件权限需要进一步考虑
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			panic("创建文件夹时出错！")
		}
	}

	return dirPath
}

//移除某个文件夹。 例如：/var/www/matter => /var/www
func RemoveDirectory(dirPath string) string {

	exists, err := PathExists(dirPath)
	if err != nil {
		panic("判断文件是否存在时出错！")
	}
	if exists {

		err = os.Remove(dirPath)
		if err != nil {
			panic("删除文件夹时出错！")
		}
	}

	return dirPath
}

//获取配置文件存放的位置
//例如：C:\Users\lishuang\AppData\Local\Temp/conf
func GetConfPath() string {

	homePath := GetHomePath()
	filePath := homePath + "/conf"
	exists, err := PathExists(filePath)
	if err != nil {
		panic("判断日志文件夹是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			panic("创建日志文件夹时出错！")
		}
	}

	return filePath
}

//获取日志的路径
//例如：默认存放于 home/log
func GetLogPath() string {

	homePath := GetHomePath()
	filePath := homePath + "/log"
	exists, err := PathExists(filePath)
	if err != nil {
		panic("判断日志文件夹是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			panic("创建日志文件夹时出错！")
		}
	}

	return filePath
}

//获取某个用户文件应该存放的位置。这个是相对GetFilePath的路径
//例如：/zicla/root
func GetUserFileDir(username string, cache bool) (absolutePath string, relativePath string) {

	filePath := CONFIG.MatterPath

	if cache {
		//如果是缓存文件夹，那么统一放在cache这个文件夹下面
		firstDir := MATTER_CACHE

		now := time.Now()
		datePath := now.Format("2006-01-02")
		//毫秒时间戳
		timestamp := now.UnixNano() / 1e6

		absolutePath = fmt.Sprintf("%s/%s/%s/%s/%d", filePath, username, firstDir, datePath, timestamp)
		relativePath = fmt.Sprintf("/%s/%s/%s/%d", username, firstDir, datePath, timestamp)

	} else {
		firstDir := MATTER_ROOT
		absolutePath = fmt.Sprintf("%s/%s/%s", filePath, username, firstDir)
		relativePath = fmt.Sprintf("/%s/%s", username, firstDir)
	}

	exists, err := PathExists(absolutePath)
	if err != nil {
		panic("判断上传文件是否存在时出错！请检查文件夹 " + filePath + " 的访问权限。")
	}
	if !exists {
		err = os.MkdirAll(absolutePath, 0777)
		if err != nil {
			panic("创建上传文件夹时出错！")
		}
	}

	return absolutePath, relativePath
}
