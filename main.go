package main

import (
	"fmt"
	"./Scene"
	"flag"
)

var (
	Dir = flag.String("d","","文件夹路径")
)

func main()  {
	flag.Parse()
	if Dir == nil || len(*Dir) == 0 {
		flag.Usage()

		return
	}
	scene := Scene.Scene{JsonDir:*Dir}
	err := scene.Dir2Obj()
	if err != nil {
		fmt.Println(err)
		return
	}

}

