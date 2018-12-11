package Scene

import (

	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
)

type Position struct {
	Data [] float64 `json:"data"`
	Type string `json:"type"`
	Components int `json:"components"`
} 

type Vertice struct {
	Position *Position `json:"position"`
	TexCoord0 *Position `json:"texCoord0"`
	Normal *Position `json:"normal"`
}

type Mesh struct {
	Indices []int `json:"indices"`
	Vertices int `json:"vertices"`
}

type Node struct {
	Position []float64 `json:"position"`
	Rotation []float64 `json:"rotation"`
	Scale []float64 `json:"scale"`
	Name string `json:"name"`
}

type Model struct {
	Vertices []*Vertice `json:"vertices"`
	Nodes []Node `json:"nodes"`
	Meshes []*Mesh `json:"meshes"`
	MeshInstances []*MeshInstance `json:"meshInstances"`
}

type MeshInstance struct {
	Node int `json:"node"`
	Mesh int `json:"mesh"`
}

type Scene struct {
	Model Model `json:"model"`
	JsonDir string
	JsonFile string
}

func (this *Scene)Dir2Obj()error {

	filepath.Walk(this.JsonDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				if ok, _ := filepath.Match("*.json", filepath.Base(path)); !ok {
					return nil
				}
				scene := &Scene{JsonFile:path,JsonDir:this.JsonDir}

				datas,err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				err = json.Unmarshal(datas,scene)
				if err != nil {
					return err
				}
				scene.ToObj()

			}
		} else {
			fmt.Println(err)
		}
		return err
	})



	return nil
}

func (this *Scene)Json2Scene()error {


	datas,err := ioutil.ReadFile(this.JsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(datas,this)
	if err != nil {
		return err
	}

	return nil
}


func (this *Scene)ToObj() error  {
	os.MkdirAll(this.JsonDir+"/objs",os.ModePerm)
	for _, chunk := range this.Model.MeshInstances {
		f,err := os.Create(this.JsonDir+"/objs/"+filepath.Base(this.JsonFile)+"_"+this.Model.Nodes[chunk.Node].Name+".obj")

		if err != nil {
			fmt.Println(err)
			continue
		}
		defer f.Close()

		f.WriteString("# Created by wangxiaohua, a free and open source\r\n")
		f.WriteString("# OBJ serializer for JavaScript\r\n")
		/*f.WriteString("mtllib "+chunk.ChunkName+".mtl\r\n")
		f.WriteString("o g "+chunk.ChunkName+"\r\n")
		f.WriteString("usemtl "+chunk.ChunkName+"\r\n")*/

		mesh := this.Model.Meshes[chunk.Mesh]
		vertice := this.Model.Vertices[mesh.Vertices]
		xyz := vertice.Position.Data
		for i := 0; i < len(xyz) / 3; i++ {
			f.WriteString(fmt.Sprintf("v %f %f %f\r\n",xyz[i*3 + 0],xyz[i*3 + 1],xyz[i*3 + 2]))
		}

		uv := vertice.TexCoord0.Data
		for i := 0; i < len(uv) / 2; i++ {
			f.WriteString(fmt.Sprintf("vt %f %f\r\n",uv[i*2 + 0],uv[i*2 + 1]))
		}


		vn := vertice.Normal.Data
		for i := 0; i < len(vn) / 3; i++ {
			f.WriteString(fmt.Sprintf("vn %f %f %f\r\n",uv[i*2 + 0],uv[i*2 + 1]))
		}


		fmt.Println(len(xyz),len(uv),len(vn))
		indexs := mesh.Indices
		for i := 0; i < len(indexs) / 3; i++ {
			f.WriteString(fmt.Sprintf("f %d/%d/%d %d/%d/%d %d/%d/%d\r\n",
				indexs[i*3 + 0]+1,indexs[i*3 + 0]+1,indexs[i*3 + 0]+1,
					indexs[i*3 + 1]+1,indexs[i*3 + 1]+1,indexs[i*3 + 1]+1,
				indexs[i*3 + 2]+1,indexs[i*3 + 2]+1,indexs[i*3 + 2]+1,))
		}

		//saveMtl(chunk.ChunkName,chunk.MaterialName)
	}
	return nil
}