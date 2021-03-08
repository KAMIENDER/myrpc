package service

import "reflect"

type ServiceInfo struct {
	name2handler map[string]*reflect.Value
	name2params  map[string][]reflect.Type
	name2results map[string][]reflect.Type
}

func NewServiceInfo(
	name2handler map[string]*reflect.Value,
	name2params map[string][]reflect.Type,
	name2result map[string][]reflect.Type,
) *ServiceInfo {
	return &ServiceInfo{
		name2handler: name2handler,
		name2params:  name2params,
		name2results: name2result,
	}
}

func (i ServiceInfo) Handler(methodName string) (*reflect.Value, bool) {
	handler, ok := i.name2handler[methodName]
	return handler, ok
}

func (i ServiceInfo) ParamsTypes(methodName string) ([]reflect.Type, bool) {
	paramsTypes, ok := i.name2params[methodName]
	return paramsTypes, ok
}

func NewServiceInfoByService(service RPCService) *ServiceInfo {
	name2handler := make(map[string]*reflect.Value)
	name2params := make(map[string][]reflect.Type)
	name2result := make(map[string][]reflect.Type)

	serviceType := reflect.ValueOf(service)
	serviceType.NumMethod()

	elemV := serviceType.Elem()
	elemT := elemV.Type()
	fieldNum := elemV.NumMethod()
	for i := 0; i < fieldNum; i++ {
		t := elemT.Method(i)
		v := elemV.Method(i)
		if v.Kind() == reflect.Func {
			vType := v.Type()

			args := make([]reflect.Type, 0)
			result := make([]reflect.Type, 0)

			for i := 0; i < vType.NumIn(); i++ {
				arg := v.Type().In(i)
				args = append(args, arg)
			}
			for i := 0; i < vType.NumOut(); i++ {
				result = append(result, vType.Out(i))
			}
			name2handler[t.Name] = &v
			name2result[t.Name] = result
			name2params[t.Name] = args
		}
	}
	return NewServiceInfo(name2handler, name2params, name2result)
}
