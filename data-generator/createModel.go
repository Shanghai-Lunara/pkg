package data_generator

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	err1      error
	structDir string
	sourceDir string
	genGoFile string
)

// xlsx 类型转 go
var typeMap = map[string]string{
	"int":                          "int32",
	"string":                       "string",
	"float":                        "float32",
	"List<int>":                    "[]int32",
	"List<float>":                  "[]float32",
	"List<string>":                 "[]string",
	"Dictionary<int,List<int>>":    "map[int32][]int32",
	"Dictionary<int,List<string>>": "map[int32][]string",
	"Dictionary<int,int>":          "map[int32]int32",
	"Dictionary<int,string>":       "map[int32]string",
}

// go 转 protobuf tag type
var proMap = map[string]string{
	"int32":              "varint",
	"string":             "bytes",
	"float32":            "bytes",
	"[]int32":            "bytes",
	"[]float32":          "bytes",
	"map[int32][]int32":  "bytes",
	"map[int32][]string": "bytes",
	"map[int32]int32":    "bytes",
	"map[int32]string":   "bytes",
}

// 生成 structs 用的map
var structMap []string

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func ToStruct(fileName string, arr [][]string) {
	var f *os.File

	f = CreateFile(fileName, "go", structDir)
	Format(fileName, arr, f)
}

func CreateFile(fileName, fileType, path string) *os.File {
	filename := fmt.Sprintf("./%s/%s.%s", path, fileName, fileType)
	var f *os.File
	if checkFileIsExist(filename) { //如果文件存在
		_ = os.Remove(fileName)
	}
	f, err1 = os.Create(filename) //创建文件

	if err1 != nil {
		panic(err1)
	}
	return f
}

// 格式化struct,并写入文件
func Format(fileName string, arr [][]string, f *os.File) {
	writeString := fmt.Sprintf("package structure\n\ntype %s struct {\n", strings.ToUpper(fileName[0:1])+fileName[1:])
	_, err1 = io.WriteString(f, writeString) //写入文件(字符串)

	defer f.Close()

	for i := 0; i < len(arr[0]); i++ {
		if arr[0][i] == "" {
			break
		}

		if ok := IsIgnored(arr[0][i]); ok {
			continue
		}

		jsonName := strings.ToLower(arr[0][i])
		vname := returnName(jsonName)
		tname := typeMap[arr[1][i]]
		writeString = fmt.Sprintf("\t%-15s %-20s `json:\"%s\" protobuf:\"%s,%d,opt,name=%s\"`\n", vname, tname, jsonName, proMap[tname], i+1, jsonName)
		_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
	}

	writeString = "}"
	_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
}

// 检查是否为可忽略字段: remarks 或 _开头字段为可忽略字段
func IsIgnored(name string) bool {
	if ok := strings.HasPrefix(name, "_"); ok {
		return true
	}

	if name == "remarks" {
		return true
	}

	return false
}

// 返回驼峰格式命名
func returnName(str string) string {
	strByte := []byte(str)

	reg := regexp.MustCompile(`_`)
	temp := reg.FindAllIndex(strByte, -1)

	// 将下划线后的字母大写
	for _, value := range temp {
		str = string(strByte[0:value[1]]) + strings.ToUpper(string(strByte[value[1]:value[1]+1])) + string(strByte[value[1]+1:])
		strByte = []byte(str)
	}
	// 去掉下划线
	str = reg.ReplaceAllString(str, ``)
	// 首字母大写
	str = strings.ToUpper(str[0:1]) + str[1:]

	return str
}

func CreatModel(strDir, souDir, goFile, jsDir string) {

	structDir = strDir
	sourceDir = souDir
	genGoFile = goFile
	jsonDir = jsDir

	fileMap := GetFileMap(sourceDir)

	klog.Info("Create Model START")
	for sheetName, f := range fileMap {
		klog.Info(sheetName, "    start")
		rows, err := f.GetRows("Sheet1")
		klog.Info("rows: ", rows[0:2])
		// 生成 struct
		ToStruct(sheetName, rows[0:2])
		// 将filename 记录到structMap中
		structMap = append(structMap, sheetName)
		if err != nil {
			klog.Fatal(err)
		}
		klog.Info(sheetName, "    end")
	}
	klog.Info("Create Model END")
	klog.Info(structMap)

	createStructs()

	createData()

	createStructMap()
}

func getSheet(fileName string) string {
	reg := regexp.MustCompile(`\.xlsx`)
	fileName = reg.ReplaceAllString(fileName, ``)

	fileName = strings.ToLower(fileName)

	return fileName
}

func GetFileMap(souDir string) map[string]*excelize.File {
	files, _ := ioutil.ReadDir(souDir)
	var fileMap = map[string]*excelize.File{}

	for _, file := range files {
		if path.Ext(file.Name()) != ".xlsx" || path.Base(file.Name()) == "Character.xlsx" || path.Base(file.Name()) == "Error.xlsx" || path.Base(file.Name()) == "LanguageCN.xlsx" || path.Base(file.Name()) == "Calendar.xlsx" || path.Base(file.Name()) == "Plot.xlsx" || path.Base(file.Name()) == "Dialog.xlsx" {
			continue
		} else {
			f, err := excelize.OpenFile(souDir + "/" + file.Name())
			if err != nil {
				klog.Fatal(err)
			}
			sheetName := getSheet(file.Name())
			fileMap[sheetName] = f
		}
	}
	return fileMap
}

// 生成 structs
func createStructs() {
	var f *os.File
	var writeString string

	f = CreateFile("structs", "go", structDir)

	_, err1 = io.WriteString(f, "package structure\n\ntype Data struct {\n") //写入文件(字符串)

	defer f.Close()

	for i, val := range structMap {
		name := strings.ToUpper(val[0:1]) + val[1:]
		writeString = fmt.Sprintf("\t%-25s %-35s `json:\"%s\" protobuf:\"bytes,%d,opt,name=%s\"`\n", name, "map[int32]*"+name, val, i+1, val)
		_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
	}
	writeString = "}"
	_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
}

// 生成 data/data.go
func createData() {

	var f *os.File
	var writeString string

	f = CreateFile("data", "go", genGoFile)
	_, err1 = io.WriteString(f, `package data

import (
	"encoding/json"
	"github.com/Shanghai-Lunara/hamster-server/data/structure"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

// auto-generated by generator

type Data struct {
	//Config []*structure.Config
	cache *structure.Data
}

var directory string
func NewData(dir string) *structure.Data {
	directory = dir
	d := &Data{
		cache: &structure.Data{},
	}
	d.runAll()
	return d.cache
}`) //写入文件头

	var loadTemplate = `

func (d *Data) load%s() {
	f, err := os.Open(directory +"/%s.json")

	if err != nil {
		klog.Fatal(err)
	}
	content, _ := ioutil.ReadAll(f)
	var employeeArr map[int32]*structure.%s
	json.Unmarshal(content, &employeeArr)

	d.cache.%s= employeeArr
}`

	var runAllTemplate = `

func (d *Data) runAll() {`

	defer f.Close()

	for _, val := range structMap {
		strByte := []byte(val)
		upperName := strings.ToUpper(string(strByte[0:1])) + string(strByte[1:])
		writeString = fmt.Sprintf(loadTemplate, upperName, val, upperName, upperName)
		_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
		runAllTemplate += fmt.Sprintf("\n\td.load%s()", upperName)
	}

	runAllTemplate += "\n}"
	_, err1 = io.WriteString(f, runAllTemplate)
	klog.Info(runAllTemplate)
}

// 生成structMap.go
func createStructMap() {
	var f *os.File
	var writeString string

	f = CreateFile("structMap", "go", structDir)

	_, err1 = io.WriteString(f, "package structure\n\nvar RegStruct = make(map[string]interface{})\n\nfunc InitStructMap() map[string]interface{} {\n") //写入文件(字符串)

	defer f.Close()

	klog.Info("structMap: ", structMap)

	for _, val := range structMap {
		name := strings.ToUpper(val[0:1]) + val[1:]
		klog.Info(name)
		writeString = fmt.Sprintf("\t RegStruct[\"%s\"] = %s\n", val, name+"{}")
		_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
	}
	writeString = "\n\t return RegStruct \n}"
	_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
}
