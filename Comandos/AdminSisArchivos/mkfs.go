package Adminsisarchivos

import (
	"fmt"
	"strings"
)

func MKfs(entrada []string) (string){
	var respuesta string

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro,"")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			break
		}

		
	}

	return respuesta
}