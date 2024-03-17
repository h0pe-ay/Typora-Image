# Typora-Image

## 前言

由于使用Typora插件配合Picgo上传图片时总是会卡死，因此利用go开发一款用于扫描Typroa的程序，将本地图片推送到图床上。

## 功能

- 扫描Typora文件
- 上传本地图片
- 下载图床图片
- ~~集成到Typora里（待开发）~~

# 使用方法

```shell
main.exe <image-path> <upload|download>

#上传图片
main.exe /path/test.md upload
#下载图片
main.exe /path/test.md download
```

