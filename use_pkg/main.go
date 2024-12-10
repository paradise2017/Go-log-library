package main

/*
1:如果结构体字段的首字母大写（如 Log_file_path），该字段是 exported（可导出的），也就是说，其他包可以访问和修改该字段。
2:如果结构体字段的首字母小写（如 log_file_path），该字段是 unexported（不可导出的），只能在定义该结构体的包内访问。
*/
import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	mylogger "github.com/paradise2017/hello/package/mylogger"
)

// 定义抽象接口 执行多态
var logger mylogger.Logger

// conf 是 structtag
type LogConfig struct {
	Log_file_path string `conf:"log_file_path"`
	Log_file_name string `conf:"log_file_name"`
	Max_size      int64  `conf:"max_size"`
}

// conf读取内容赋值给结构体指针
func ParseConf(file_name string, result interface{}) (err error) {
	// 文件反射

	t := reflect.TypeOf(result)  //result类型
	v := reflect.ValueOf(result) //result值
	if t.Kind() != reflect.Ptr {
		err = errors.New("need to convey painter")
	}
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("Must be a struct pointer")
	}

	// 1: 打开文件
	data, err := ioutil.ReadFile(file_name) // data:[]byte
	if err != nil {
		err = fmt.Errorf("打开配置文件%s失败", file_name)
		return err
	}
	// 2：配置文件解析 ，读取的文件数据按照行分割，
	file_line_data := strings.Split(string(data), "\n")

	// 3：分行解析  range(index,value)
	for index, line := range file_line_data {
		line = strings.TrimSpace(line) // 去除字符串首尾的空白
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// 忽略空白和注释
			continue
		}
		// 解析配置文件
		equal_index := strings.Index(line, "=")
		if equal_index == -1 {
			fmt.Printf("line%d error\n", index+1)
		}
		// =分割  [）
		key := line[:equal_index]
		value := line[equal_index+1:]
		if len(key) == 0 {
			err = fmt.Errorf("line:%d failed", index+1)
			return
		}
		// 利用反射赋值,
		//拿到结构体的字段给 conf赋值
		// NumField()结构体字段总数
		// t.Elem().Field(i) 拿到第i个字段
		for i := 0; i < t.Elem().NumField(); i++ {
			field := t.Elem().Field(i)
			tag := field.Tag.Get("conf") // 拿到具体的key
			if key == tag {
				//匹配正确 value赋值
				// 拿到每个字段的类型
				field_type := field.Type
				switch field_type.Kind() {
				case reflect.String:
					// field_value := v.Elem().FieldByName(field.Name) // 根据字段名找到对应的值
					// field_value.SetString(value)                    // 配置文件数据数据设置给value
					v.Elem().Field(i).SetString(value)
				case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int32:
					// strconv字符串转换为任意字段的整形
					// 将value按照10进制转化为int64 类型的整数。
					value64, _ := strconv.ParseInt(value, 10, 64)
					v.Elem().Field(i).SetInt(value64)
				}
			}
		}
	}
	return err
}
func main() {
	// "./" 编译成二进制的当前路径
	// 日志级别越高，写的越少
	// logger := mylogger.NewFileLogger(mylogger.DEBUGLevel, "./", "test.log")

	// logger.Debug("测试demo")

	// 设置的error级别，小于的都不用写
	// logger := mylogger.NewConsoleLogger(mylogger.ERRORLevel)
	// logger.ConsoleCritical("test")

	// 配置文件解析

	var c = &LogConfig{}
	err := ParseConf("log.conf", c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", c)
	fmt.Println(err)
}

// 日志切割
