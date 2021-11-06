package main

import (
	"errors"
	"reflect"
)

func main() {
	// var x float64 = 3.4
	// fmt.Println("type: ", reflect.TypeOf(x))
	// fmt.Println("value: ", reflect.ValueOf(x).String())

	// v := reflect.ValueOf(x)
	// fmt.Println("type: ", v.Type())
	// fmt.Println("kind: ", v.Kind())
	// fmt.Println("value: ", v.Float())

	// var y float64 = 4.7
	// p := reflect.ValueOf(&y)
	// fmt.Println("type of p: ", p.Type())
	// fmt.Println("settability of p: ", p.CanSet())

	// vp := p.Elem()
	// fmt.Println("settability of vp: ", vp.CanSet())

	// vp.SetFloat(7.1)
	// fmt.Println(vp.Interface())
	// fmt.Println(y)

	// type T struct {
	// 	A int
	// 	B string
	// }

	// t := T{23, "skidoo"}
	// s := reflect.ValueOf(&t).Elem()

	// typeOfT := s.Type()

	// for i := 0; i < s.NumField(); i++ {
	// 	f := s.Field(i)
	// 	fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	// }
}

func i2s(data interface{}, out interface{}) error {
	//Проверим, что результат это указатель, чтобы можно было в него писать данные
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.New("Result must be of type pointer")
	}

	//На входе у нас всегда либо Slice, либо map из интерфейсов
	//Получается, что если приходит мап, то на выходе должна быть структура
	//Но если приходит слайс, то нужно на выходе так же отдать слайс
	typeOfData := reflect.ValueOf(data).Kind()
	typeOfOut := reflect.ValueOf(out).Elem().Kind()

	isSlice := typeOfData == reflect.Slice
	isMap := typeOfData == reflect.Map

	if isSlice && typeOfOut != reflect.Slice {
		return errors.New("Input and output types are mismatch")
	}

	if isMap && typeOfOut != reflect.Struct {
		return errors.New("Input and output types are mismatch")
	}

	if isMap {
		err := mapExtractor(data, out)
		if err != nil {
			return err
		}
	}

	if isSlice {
		sliceExtractor(data, out)
	}

	return nil
}

// Пока не понятно, как заполнять значения когда в рекурсию уходишь
// Похоже через указатель не работает совсем, нужно подебажить
func mapExtractor(input interface{}, out interface{}) error {
	inputMap := reflect.ValueOf(input)
	result := reflect.ValueOf(out).Elem()
	typeOfResult := result.Type()

	for i := 0; i < result.NumField(); i++ {
		fieldName := typeOfResult.Field(i).Name
		inputValue := inputMap.MapIndex(reflect.ValueOf(fieldName))

		//Предполагаем, что если нет во входных данных поля, то оно просто не заполняется
		if !inputValue.IsValid() {
			continue
		}

		outValue := result.Field(i)

		//Если у нас поле является мап, то его нужно снова экстрактить
		fieldStructType := reflect.TypeOf(inputValue.Interface()).Kind()
		isComplexFiled := fieldStructType == reflect.Map || fieldStructType == reflect.Slice

		if isComplexFiled {
			err := i2s(inputValue.Interface(), outValue.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		err := setActualValue(inputValue, outValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func sliceExtractor(in interface{}, out interface{}) error {
	sliceValueIn := reflect.ValueOf(in)
	sliceValueOut := reflect.ValueOf(out)

	sliceOfStruct := reflect.MakeSlice(sliceValueOut.Type().Elem(), 0, sliceValueIn.Cap())

	for i := 0; i < sliceValueIn.Len(); i++ {
		structInSlice := reflect.New(sliceValueOut.Type().Elem().Elem())
		err := i2s(sliceValueIn.Index(i).Interface(), structInSlice.Interface())
		if err != nil {
			return err
		}
		sliceOfStruct = reflect.Append(sliceOfStruct, structInSlice.Elem())
	}

	sliceValueOut.Elem().Set(sliceOfStruct)
	return nil
}

func setActualValue(in reflect.Value, out reflect.Value) error {
	outputType := out.Type().Kind()

	i, ok := in.Interface().(int64)
	if ok && outputType == reflect.Int {
		out.SetInt(i)
		return nil
	}

	isNumericType := outputType == reflect.Float64 || outputType == reflect.Int

	f, ok := in.Interface().(float64)
	if ok && isNumericType {
		if outputType == reflect.Int {
			i := int64(f)
			out.SetInt(i)
			return nil
		}

		if outputType == reflect.Float64 {
			out.SetFloat(f)
			return nil
		}
	}

	s, ok := in.Interface().(string)
	if ok && outputType == reflect.String {
		out.SetString(s)
		return nil
	}

	b, ok := in.Interface().(bool)
	if ok && outputType == reflect.Bool {
		out.SetBool(b)
		return nil
	}

	return errors.New("Input type: " + in.Elem().Type().String() + ", output type: " + outputType.String())
}
