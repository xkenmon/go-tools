# Gnome壁纸配置生成器

## 简介

Ubuntu的Gnome桌面默认没有自动切换壁纸的功能，但可以通过添加主题配置文件来实现定时切换

文件太多写起来非常麻烦，使用这个生成器可以快速生成主题配置文件

理论上适用于所有的Gnome桌面


## 生成器用法

```
Usage of ./bg-generate:
  -f string
    	生成文件存放路径，为空输出到标准输出流
  -keep float
    	每张壁纸持续时间(S) (default 2000)
  -path string
    	壁纸文件所在文件夹 (default ".")
  -suffix string
    	壁纸文件后缀列表(以空格分割,如 ".jpg .png .jpeg") (default ".jpg .jpeg .png .gif .bmp")
  -transit float
    	切换壁纸时的动画时长(S) (default 5)
```

## 配置文件用法

编辑gnome桌面背景配置文件

```bash
vim /usr/share/gnome-background-properties/`lsb_release -sc`-wallpapers.xml
```

在该配置文件中将刚生成的配置文件的位置添加进去

```
  <wallpaper>
      <name>custome name</name>
      <filename>/your/config/file/path</filename>
      <options>zoom</options>
  </wallpaper>
```

然后换主题，就OK了.

