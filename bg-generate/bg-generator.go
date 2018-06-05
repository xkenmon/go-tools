package main

import "bytes"
import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "path/filepath"
import "strings"

const (
	bgconfig_path  = "~/.bg-generator"
	sysconfig_path = "/usr/share/gnome-background-properties"
)

func main() {

	path := flag.String("path", "", "壁纸文件所在文件夹")
	keep := flag.Float64("keep", 2000.0, "每张壁纸持续时间(S)")
	transit := flag.Float64("transit", 5, "切换壁纸时的动画时长(S)")
	suffix := flag.String("suffix", ".jpg .jpeg .png .gif .bmp", "壁纸文件后缀列表(以空格分割,如 \".jpg .png .jpeg\")")
	f := flag.String("f", "", "生成文件存放路径，为空输出到标准输出流")
	w := flag.Bool("w", false, "是否写入到系统配置中（只适用于Gnome）,若指定该选项，-f 参数失效")

	flag.Parse()

	if *path == "" {
		fmt.Println("请输入图片路径， -h查看帮助")
		return
	}

	// suffix array
	sfx := strings.Split(*suffix, " ")

	content, err := generate(*path, *keep, *transit, sfx)
	if err != nil {
		fmt.Errorf("生成配置文件出错:%v\n", err)
		return
	}

	// need write to sys config
	if *w {
		abs_config_path, err := filepath.Abs(bgconfig_path)

		CheckAndMakeDir(abs_config_path)

		if err := ioutil.WriteFile(bgconfig_path+"/bg-generator.xml", []byte(content), 0644); err != nil {
			fmt.Errorf("写入到路径%s出错!", bgconfig_path)
			return
		}
		syscontent := `<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE wallpapers SYSTEM "gnome-wp-list.dtd">
		<wallpapers>
		<wallpaper>
		<name>bg-generator</name>
		<filename>%s</filename>
		<options>zoom</options>
		</wallpaper>
		</wallpapers>
		`

		CheckAndMakeDir(sysconfig_path)

		err = ioutil.WriteFile(sysconfig_path+"/bg-generator.xml", []byte(fmt.Sprintf(syscontent, abs_config_path+"/bg-generator.xml")), 0644)
		if err != nil {
			fmt.Errorf("写入到系统配置中出错:%v，写入路径%s", err, sysconfig_path)
		}
		fmt.Printf("success!, sys config path:%s\n", sysconfig_path+"/bg-generator.xml")
		return
	}

	if err != nil {
		fmt.Errorf("生成文件出错:%s", err)
		return
	}

	// store to file
	if *f == "" {
		fmt.Println(content)
	} else {
		fp, err := filepath.Abs(*f)
		if err != nil {
			fmt.Errorf("获取文件路径出错：%s", err)
			return
		}

		CheckAndMakeDir(filepath.Dir(fp))
		fmt.Printf("-f 绝对路径:%s\n", fp)
		if err = ioutil.WriteFile(fp, []byte(content), 0644); err != nil {
			fmt.Errorf("写入文件出错:%s", err)
			return
		}
		fmt.Printf("写入成功，-> %s\n", fp)
	}
}
func generate(path string, keep, transit float64, suffix []string) (string, error) {
	start :=
		`<background>
	<starttime>
	<year>2009</year>
	<month>08</month>
	<day>04</day>
	<hour>00</hour>
	<minute>00</minute>
	<second>00</second>
	</starttime>
	<!-- This animation will start at midnight. -->
	`
	end := "</background>"

	static := `  
	<static>
	<duration>%.1f</duration>
	<file>%s</file>
	</static>
	`
	transition := `
	<transition>
	<duration>%.1f</duration>
	<from>%s</from>
	<to>%s</to>
	</transition>
	`

	last := ""
	buf := bytes.Buffer{}

	buf.WriteString(start)

	filenames, err := ListFile(path, suffix)
	if err != nil {
		return "", err
	}

	for _, fn := range filenames {
		if last == "" {
			buf.WriteString(fmt.Sprintf(static, keep, fn))
			last = fn
		} else {
			buf.WriteString(fmt.Sprintf(transition, transit, last, fn))
			buf.WriteString(fmt.Sprintf(static, keep, fn))
			last = fn
		}
	}

	buf.WriteString(end)

	return buf.String(), nil
}

func ListFile(dir string, suffix []string) (filenames []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fi := range files {
		fn, _ := filepath.Abs(dir + string(os.PathSeparator) + fi.Name())
		for _, s := range suffix {
			if !fi.IsDir() && filepath.Ext(fn) == s {
				filenames = append(filenames, fn)
			}
		}
	}
	return
}

// 判断文件夹是否存在
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

func CheckAndMakeDir(dir string) {
	if b, _ := PathExists(dir); !b {
		os.MkdirAll(dir, 0744)
	}
}
