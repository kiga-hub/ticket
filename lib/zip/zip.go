package zip

import (
    "archive/zip"
    "path/filepath"
    "bytes"
    "io"
    "os"
)

func IsZip(zipPath string) bool {
    f, err := os.Open(zipPath)
    if err != nil {
        return false
    }
    defer f.Close()

    buf := make([]byte, 4)
    if n, err := f.Read(buf); err != nil || n < 4 {
        return false
    }

    return bytes.Equal(buf, []byte("PK\x03\x04"))
}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func DeCompressByPath(tarFile, dest string) ([]string, error) {
    srcFile, err := os.Open(tarFile)
    if err != nil {
        return nil, err
    }
    defer srcFile.Close()
    return DeCompress(srcFile, dest)
}

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/
func DeCompress(srcFile *os.File, target string) ([]string, error) {
    var result []string
    zipFile, err := zip.OpenReader(srcFile.Name())
    if err != nil {
        return nil, err
    }
    if err := os.MkdirAll(target, 0775); err != nil {
        return nil, err
    }
    defer zipFile.Close()
    for _, innerFile := range zipFile.File {
        result = append(result, innerFile.Name)
        path := filepath.Join(target, innerFile.Name)
        if innerFile.FileInfo().IsDir() {
            err = os.MkdirAll(path, os.ModePerm)
            if err != nil {
                return nil, err
            }
            continue
        }
        srcFile, err := innerFile.Open()
        if err != nil {
            continue
        }
        defer srcFile.Close()

        newFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, innerFile.Mode())
        if err != nil {
            continue
        }
        if _, err := io.Copy(newFile, srcFile); err != nil {
            return nil, err
        }
        newFile.Close()
    }
    return result, nil
}

/**
@files：需要压缩的文件
@compreFile：压缩之后的文件
*/
func Compress_zip(files []*os.File, compreFile *os.File) (err error) {
    zw := zip.NewWriter(compreFile)
    defer zw.Close()
    for _, file := range files {
        err := compress_zip(file, zw)
        if err != nil {
            return err
        }
        file.Close()
    }
    return nil
}

/**
功能：压缩文件
@file:压缩文件
@prefix：压缩文件内部的路径
@tw：写入压缩文件的流
*/
func compress_zip(file *os.File, zw *zip.Writer) error {
    info, err := file.Stat()
    if err != nil {
        return err
    }
    // 获取压缩头信息
    head, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }
    // 指定文件压缩方式 默认为 Store 方式 该方式不压缩文件 只是转换为zip保存
    head.Method = zip.Deflate
    fw, err := zw.CreateHeader(head)
    if err != nil {
        return err
    }
    // 写入文件到压缩包中
    _, err = io.Copy(fw, file)
    file.Close()
    if err != nil {
        return err
    }
    return nil
}