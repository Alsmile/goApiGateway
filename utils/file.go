package utils

import (
  "strings"
  "os"
  "path/filepath"
  "io/ioutil"
)

// 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
  files = make([]string, 0, 30)
  suffix = strings.ToUpper(suffix)
  err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
    // 忽略目录
    if fi.IsDir() {
      return nil
    }

    if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
      files = append(files, filename)
    }
    return nil
  })
  return files, err
}

func ReadFile(path string) (error, string) {
  fileHandle, err := os.Open(path)
  if err != nil {
    return err, ""
  }
  defer fileHandle.Close()
  fileBytes, err := ioutil.ReadAll(fileHandle)
  return err, string(fileBytes)
}
