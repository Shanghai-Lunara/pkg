package data_generator

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
	"k8s.io/klog/v2"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	jsonDir   string
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

	for _, row := range arr[4:] {
		if len(row) == 0 {
			break
		}
		if row[0] == "" {
			break
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
				parseInt := ToInt(row[i])
				if i == 0 {
					id = int32(parseInt)
				}
				newValue.FieldByName(field.Name).SetInt(parseInt)
			case "string":
				newValue.FieldByName(field.Name).SetString(row[i])
			case "float32":
				newValue.FieldByName(field.Name).SetFloat(ToFloat(row[i]))
			case "[]int32":
				tmp := ToStringSlice(row[i])
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToIntSlice(tmp)))
			case "[]float32":
				tmp := ToStringSlice(row[i])
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToFloatSlice(tmp)))
			case "[]string":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToStringSlice(row[i])))
			case "map[int32][]int32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToListIntMap(row[i])))
			case "map[int32][]string":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToListStringMap(row[i])))
			case "map[int32]int32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToIntMap(row[i])))
			case "map[int32]string":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToStringMap(row[i])))
			case "map[int32][]float32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToListFloatMap(row[i])))
			case "map[int32]float32":
				newValue.FieldByName(field.Name).Set(reflect.ValueOf(ToFloatMap(row[i])))
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
	wg := sync.WaitGroup{}
	wg.Add(len(fileMap))
	for k, v := range fileMap {
		go func(sheetName string, f *excelize.File) {
			defer wg.Done()
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

			ff := CreateFile(sheetName, "json", jsonDir)
			_, err = io.WriteString(ff, string(jsons))
			if err != nil {
				klog.Fatal("file: ", sheetName, "err: ", err.Error())
			}
			ff.Close()
			klog.Info(sheetName, "    end")
		}(k, v)

	}
	wg.Wait()
	klog.Info("Create JSON END")
}

func ToInt(str string) int64 {
	parseInt, _ := strconv.ParseInt(str, 10, 64)
	return parseInt
}

func ToFloat(str string) float64 {
	parseFloat, _ := strconv.ParseFloat(str, 32)
	return parseFloat
}

// 格式化成[]string
func ToStringSlice(str string) []string {
	if str == "" {
		return make([]string, 0)
	}
	reg := regexp.MustCompile(`[{}]`)
	tmp := reg.ReplaceAllString(str, ``)

	return strings.Split(tmp, ",")
}

// 将[]string 转化成 []int32
func ToIntSlice(str []string) []int32 {
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
func ToFloatSlice(str []string) []float32 {
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
func ToStringMap(str string) map[int32]string {
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
func ToIntMap(str string) map[int32]int32 {
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
func ToListStringMap(str string) map[int32][]string {
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

// 转化成map[int32]float32
func ToFloatMap(str string) map[int32]float32 {
	var ret = make(map[int32]float32)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		if val == "" {
			continue
		}
		tmp1 := strings.Split(val, ":")

		j, err := strconv.ParseFloat(tmp1[1], 10)
		if err != nil {
			klog.Fatal(err)
		}

		index, err := strconv.ParseInt(tmp1[0], 10, 32)
		if err != nil {
			klog.Fatal(err)
		}

		ret[int32(index)] = float32(j)

	}

	return ret
}

// 转化成map[int32][]int32
func ToListIntMap(str string) map[int32][]int32 {
	var ret = make(map[int32][]int32)
	if str == "" {
		return ret
	}

	tmp := strings.Split(str, ";")

	for _, val := range tmp {
		if val == "" {
			continue
		}
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
func ToListFloatMap(str string) map[int32][]float32 {
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
