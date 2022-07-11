package main

import (
	"io/ioutil"
	"os"
	"fmt"

	"github.com/valyala/fastjson"

	// "encoding/json" doesn't work due to https://github.com/tinygo-org/tinygo/issues/447
	// "github.com/json-iterator/tinygo" doesn't work due to not supporting `interface{}``.
)

var p fastjson.Parser

func main() {
	err := run()
	checkErr(err)
}

func run() error {
	jsonInput, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

    v, err := p.Parse(string(jsonInput))
	if err != nil {
		return err
	}

	nsBytes := v.GetStringBytes("functionConfig", "data", "namespace")
	ns := fmt.Sprintf(`"%v"`, string(nsBytes))

	items := v.GetArray("items")

	for _, item := range items {
		metadata := item.Get("metadata")
		metadata.Set("namespace", fastjson.MustParse(ns))
	}

	mj := v.MarshalTo(nil)
	fmt.Println(string(mj))

	// can't use unstructured, since k8s.io/apimachinery depends on klog which depends on os/user that is not supported by tinygo compiler.
	// var rl map[string]interface{}
	// err = json.Unmarshal(jsonInput, &rl)
	// if err != nil {
	// 	return err
	// }

	// ns, found, err := NestedString(rl, "functionConfig", "data", "namespace")
	// if err != nil {
	// 	return err
	// }
	// if !found {
	// 	return fmt.Errorf("unable to find namespace in functionConfig")
	// }
	
	// items, found, err := NestedSlice(rl, "items")
	// if err != nil {
	// 	return err
	// }
	// if !found {
	// 	return fmt.Errorf("unable to find items")
	// }

	// for _, obj := range items {
	// 	typedObj := obj.(map[string]interface{})
	// 	err = SetNestedField(typedObj, ns, "metadata", "namespace")
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// mj, err := json.Marshal(rl)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(mj))
	
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
