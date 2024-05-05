package controller

import (
	"backend/model"
	"backend/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileListResponse struct {
	Files []string `json:"files"`
}

func GetTestCaseList(context *gin.Context) {
	path := GetPathByPid(context)
	if strings.Contains(path, "Error") {
		response.Response(context, http.StatusBadRequest, 400, nil, path)
		return
	}
	// 构建文件列表
	var fileList []model.TestCaseFile
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		response.Response(context, http.StatusBadRequest, 400, nil, "发生错误")
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue // 跳过子目录
		}

		fileExt := filepath.Ext(entry.Name())
		fileType := fileExt[1:] // 去除前导的点号

		filePath := filepath.Join(path, entry.Name())

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("无法获取文件信息:", err)
			continue
		}
		modTime := fileInfo.ModTime()
		file := model.TestCaseFile{
			Name:     entry.Name(),
			Bytes:    uint64(fileInfo.Size()),
			ModTime:  modTime.Format("2006-01-02 15:04:05"),
			FileType: fileType,
		}

		fileList = append(fileList, file)
	}

	// 发送JSON响应
	context.JSON(http.StatusOK, gin.H{
		"data": fileList,
	})
}

func UploadTestCase(context *gin.Context) {
	path := GetPathByPid(context)
	if strings.Contains(path, "Error") {
		response.Response(context, http.StatusBadRequest, 400, nil, path)
		return
	}
	// Get the uploaded file
	err := context.Request.ParseMultipartForm(10 << 20) // 解析表单数据, 限制数据大小为 10mb
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "解析表单数据错误")
		return
	}
	files := context.Request.MultipartForm.File["files[]"]
	if len(files) == 0 {
		response.Response(context, http.StatusBadRequest, 400, nil, "未找到上传文件")
		return
	}
	// 处理每个上传的文件
	for _, file := range files {
		// 打开文件
		fileContent, err := file.Open()
		if err != nil {
			response.Response(context, http.StatusInternalServerError, 500, nil, "打开文件错误")
			return
		}
		defer fileContent.Close()

		// 生成目标文件路径
		filePath := filepath.Join(path+"/", file.Filename)

		// 创建目标文件并保存文件内容
		outFile, err := os.Create(filePath)
		if err != nil {
			response.Response(context, http.StatusInternalServerError, 500, nil, "保存文件错误")
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, fileContent)
		if err != nil {
			response.Response(context, http.StatusInternalServerError, 500, nil, "保存文件错误")
			return
		}
	}

	// 文件保存成功
	response.Response(context, http.StatusOK, 200, nil, "文件上传成功")
}

func DeleteTestCase(context *gin.Context) {
	path := GetPathByPid(context)
	if strings.Contains(path, "Error") {
		response.Response(context, http.StatusBadRequest, 400, nil, path)
		return
	}
	err := os.Remove(path)
	if err != nil {
		// 处理删除文件失败的错误
		response.Response(context, http.StatusBadRequest, 400, nil, "删除文件错误")
		return
	}
	// 文件删除成功
	response.Response(context, http.StatusOK, 200, nil, "文件上传成功")
}

func GetTestCaseDetail(context *gin.Context) {
	path := GetPathByPid(context)
	if strings.Contains(path, "Error") {
		response.Response(context, http.StatusBadRequest, 400, nil, path)
		return
	}
	log.Println(path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件错误:", err)
		response.Response(context, http.StatusInternalServerError, 500, nil, "打开文件错误")
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("读取文件内容错误:", err)
		response.Response(context, http.StatusInternalServerError, 500, nil, "读取文件内容错误")
		return
	}
	if err != nil {
		response.Response(context, http.StatusInternalServerError, 500, nil, "解码数据失败")
		return
	}
	response.Response(context, http.StatusOK, 200, gin.H{
		"data": string((fileContent)),
	}, "读取文件内容成功")
}

func UpdateCase(context *gin.Context) {
	path := GetPathByPid(context)
	if strings.Contains(path, "Error") {
		response.Response(context, http.StatusBadRequest, 400, nil, path)
		return
	}
	// 获取 content
	content := context.PostForm("content")

	// 将 content 写入文件
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		// 处理写入文件错误
		response.Response(context, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	// 返回响应
	response.Response(context, http.StatusOK, 200, nil, "文件更新成功")
}

func GetPathByPid(context *gin.Context) string {
	pid := context.Query("pid")
	FileName := context.DefaultQuery("fname", "0")
	log.Println(FileName)
	problemID, _ := strconv.ParseUint(pid, 10, 64)
	path, err := GetPathByProblemId(problemID)
	if err != nil {
		return "Error, 获取测试点路径错误"
	}
	path = strings.Replace(path, "backend/", "./", 1)

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("文件夹不存在")
			return "Error, 文件夹不存在"
		} else {
			log.Println("发生错误:", err)
			return "Error, 发生错误"
		}
	}
	if len(FileName) > 1 {
		return path + "/" + FileName
	}
	return path
}
