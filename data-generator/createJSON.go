package data_generator

import (
	"encoding/json"
	"io"
	"k8s.io/klog/v2"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	jsonDir string
	structure map[string]interface{}
)

func newStruct(mod string) interface{} {
	if structure[mod] != nil {
		return structure[mod]
	} else {
		return nil
	}
}

func ToArray(fileName string, arr [][]string) map[int32]interface{} {
	var Employees = make(map[int32]interface{})

Loop:
	for _, row := range arr[4:] {
		klog.Info(row, len(row))
		if len(row) == 0 {
			break Loop
		}
		if row[0] == "" {
			break Loop
		}

		employee := newStruct(fileName)
		//klog.Info("employee: ", employee)
		//klog.Info("type employee: ", reflect.TypeOf(employee))
		getType := reflect.TypeOf(employee)
		getValue := reflect.ValueOf(&employee).Elem()
		//klog.Info("type getValue: ", reflect.TypeOf(getValue))
		newValue := reflect.New(getValue.Elem().Type()).Elem()
		newValue.Set(getValue.Elem())
		//klog.Info("type newValue: ", reflect.TypeOf(newValue))

		var id int32

		j := 0
		for i := 0; i < len(row); i++ {

			if j == getType.NumField() {
				break
			}

			if ok := IsIgnored(arr[0][i]); ok {
				continue
			}

			field := getType.Field(j)
			switch field.Type.String() {
			case "int32":
				parseInt, _ := strconv.ParseInt(row[i], 10, 64)
				if i == 0 {
					id = int32(parseInt)
				}
				newValue.FieldByName(field.Name).SetInt(parseInt)
			case "string":
				newValue.FieldByName(field.Name).SetString(row[i])
			case "float32":
				parseFloat, _ := strconv.ParseFloat(row[i], 32)
				newValue.FieldByName(field.Name).SetFloat(parseFloat)
			case "[]int32":
				tmp := toStringSlice(row[i])
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toIntSlice(tmp)))
			case "[]float32":
				tmp := toStringSlice(row[i])
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toFloatSlice(tmp)))
			case "map[int32][]int32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toListIntMap(row[i])))
			case "map[int32][]string":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toListStringMap(row[i])))
			case "map[int32]int32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toIntMap(row[i])))
			case "map[int32]string":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toStringMap(row[i])))
			case "map[int32][]float32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(toListFloatMap(row[i])))
			}

			j++
		}
		getValue.Set(newValue)

		Employees[id] = employee

	}

	return Employees
}

func CreateJSON(jsDir, souDir string, s map[string]interface{}) {
	structure = s
	jsonDir = jsDir

	klog.Info("Create JSON START")
	fileMap := GetFileMap(souDir)

	for sheetName, f := range fileMap {
		klog.Info(sheetName, "    start")
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			klog.Fatal(err)
		}
		//structure.InitStructMap()
		// TODO 生成数据[]
		temp := ToArray(sheetName, rows)

		// TODO 转换成JSON文件并存储
		jsons, err := json.MarshalIndent(temp, "", "  ")

		f := CreateFile(sheetName, "json", jsonDir)
		_, err = io.WriteString(f, string(jsons))
		if err != nil {
			klog.Fatal(err)
		}
		f.Close()
		klog.Info(sheetName, "    end")
	}
	klog.Info("Create JSON END")
}

// 格式化成[]string
func toStringSlice(str string) []string {
	if str == "" {
		return make([]string, 0)
	}
	reg := regexp.MustCompile(`[{}]`)
	tmp := reg.ReplaceAllString(str, ``)

	return strings.Split(tmp, ",")
}

// 将[]string 转化成 []int32
func toIntSlice(str []string) []int32 {
	var ret []int32
	if len(str) == 0 {
		return make([]int32, 0)
	}

	for _, i := range str {
		j, err := strconv.ParseInt(i, 10, 32)
		if err != nil {
			klog.Fatal(err)
		}
		ret = append(ret, int32(j))
	}
	return ret
}

// 将[]string 转化成 []float32
func toFloatSlice(str []string) []float32 {
	var ret []float32
	if len(str) == 0 {
		return make([]float32, 0)
	}

	for _, i := range str {
		j, err := strconv.ParseFloat(i, 32)
		if err != nil {
			klog.Fatal(err)
		}
		ret = append(ret, float32(j))
	}
	return ret
}

// 转化成map[int32]int32
func toStringMap(str string) map[int32]string {
	var ret = make(map[int32]string)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		tmp1 := strings.Split(val, ":")

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = tmp1[1]

	}

	return ret
}

// 转化成map[int32]int32
func toIntMap(str string) map[int32]int32 {
	var ret = make(map[int32]int32)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		if val == "" {
			continue
		}
		tmp1 := strings.Split(val, ":")

		j, err := strconv.ParseInt(tmp1[1], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = int32(j)

	}

	return ret
}

// 转化成map[int32][]string
func toListStringMap(str string) map[int32][]string {
	var ret = make(map[int32][]string)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		tmp1 := strings.Split(val, ":")

		value := strings.Split(tmp1[1], ",")

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = value

	}

	return ret
}

// 转化成map[int32][]int32
func toListIntMap(str string) map[int32][]int32 {
	var ret = make(map[int32][]int32)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		tmp1 := strings.Split(val, ":")

		value := strings.Split(tmp1[1], ",")
		var arr []int32
		for _, val := range value {
			j, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				klog.Fatal(err)
			}
			arr = append(arr, int32(j))
		}

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = arr

	}

	return ret
}

// 转化成map[int32][]float32
func toListFloatMap(str string) map[int32][]float32 {
	var ret = make(map[int32][]float32)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		tmp1 := strings.Split(val, ":")

		value := strings.Split(tmp1[1], ",")
		var arr []float32
		for _, val := range value {
			j, err := strconv.ParseFloat(val, 10)
			if err != nil {
				klog.Fatal(err)
			}
			arr = append(arr, float32(j))
		}

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = arr

	}

	return ret
}
