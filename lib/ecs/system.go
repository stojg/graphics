package ecs

import (
	"fmt"
	"reflect"
)

func (e *ECS) AddSystem(system interface{}, components ...Component) {
	method := reflect.ValueOf(system)
	for _, comp := range components {
		tid := e.addComponentType(comp)
		e.systemComponents[method] = append(e.systemComponents[method], tid)
	}

	num := method.Type().NumIn()
	e.systemToIn[method] = make([]reflect.Type, num)
	for i := 1; i < num; i++ {
		e.systemToIn[method][i] = method.Type().In(i)
	}
}

func (e *ECS) Update(elapsed float64) {
	in := make([]reflect.Value, 8)
	in[0] = reflect.ValueOf(elapsed)
	for method, args := range e.systemToIn {
		componentList := make([]reflect.Value, len(args)-1)

		entities := make([][]Component, 0, 0)
		for entity := range e.allEntityComponents {
			components, ok := e.canEntityBeUpdated(entity, e.systemComponents[method])
			if !ok {
				continue
			}
			entities = append(entities, components)
		}

		for i := 1; i < len(args); i++ {
			_, ok := e.allComponentTypes[args[i].Elem()]
			if !ok {
				panic(fmt.Sprintf("Can't find component type for %s", args[i].Elem()))
			}
			componentList[i-1] = reflect.MakeSlice(args[i], len(entities), len(entities))
		}

		count := 0
		for _, components := range entities {
			for i, component := range components {
				v := reflect.ValueOf(component)
				componentList[i].Index(count).Set(v)
			}
			count++
		}
		for i := 1; i < len(args); i++ {
			v, ok := e.allComponentTypes[args[i].Elem()]
			if !ok {
				panic("oh crappers")
			}
			in[i] = componentList[v]
		}
		method.Call(in[:len(args)])
	}
}

func (d *ECS) canEntityBeUpdated(entity int, componentTypes []int) ([]Component, bool) {
	result := make([]Component, len(componentTypes))
	e := Entity(entity)
	count := 0
	for i, typeID := range componentTypes {
		for _, component := range d.allEntityComponents[e] {
			if typeID == component.TID() {
				result[i] = component
				count++
				break
			}
		}
	}
	return result, count == len(componentTypes)
}
