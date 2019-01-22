# markdownServer

## feature

* /update接口锁，保证同时只有一个groutine在更新fileList。
* 缓存，除非更新，已经解析过的文件会放在缓存中，无需重复解析

## 工程结构

* ./markdowns/* 存放需要解析的markdown文件
* ./templates/* 存放工程需要的html模板

## 接口

* / index显示扫描到的文件列表
* /file/${fileName}/ 显示markdown转换后的HTML页面
* /update 调用后，重新扫描文件目录

## 使用

>由于使用packr将html模板编译入静态文件内，编译后得到的可执行包可以单独使用。

## 运行参数

* -f 配置需要扫描的markdown文件目录，默认为./markdowns;
* -t 配置html模板(目前只有index.html)目录，设置时使用包内编译的默认资源;
* -p 端口

## 编译

```
go-bindata -o assets/asset.go -pkg=assets assets/...
go build
```

## 锁说明

文件缓存所时通过tryLock实现，读取缓存失败后会尝试获取锁。获取成功，解析文件；获取失败，等待解析groutine解析完成。这个等待是通过调用tryLock对象的Lock方法，尝试去获取锁，获取到时马上释放并从缓存中获取内容。

_ps:能够获取到锁，说明另外一个groutine已经解析完成_
