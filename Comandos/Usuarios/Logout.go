package Usuarios

import (
	"fmt"
	"strings"
)

func Logout(entrada []string) string{
	var respuesta string
	Valido := true

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro,"")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR MKDIS, valor desconocido de parametros " + valores[1]+ "\n"
			Valido = false
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}
	}

	if Valido{
		fmt.Println("Validado")
	}
	return respuesta
}