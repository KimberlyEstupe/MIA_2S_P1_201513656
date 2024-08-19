package Admindiscos

import (
	"fmt"
	"os"
	"strings"
)


func Rmdisk(entrada []string){
	//Quitar espacios en blanco
	tmp := strings.TrimRight(entrada[1],"")
	valores := strings.Split(tmp,"=")
	var path string


	if len(valores)!=2{
		fmt.Println("ERROR RMDISK, valor desconocido de parametros ",valores[1])
		return
	}else{		
		path = strings.ReplaceAll(valores[1],"\"","")
	}

	//validar si existe el archivo a eliminar
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("RMDISK Error: El disco ", path, " no existe")
		return
	}

	//Eliminar disco
	err2 := os.Remove(path)
			if err2 != nil {
				fmt.Println("RMDISK Error: No pudo removerse el disco ")
				return
			}
			fmt.Println("Disco ", path, "eliminado correctamente:")

}
