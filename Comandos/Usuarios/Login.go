package Usuarios

import (
	"fmt"
	"strings"
)

func Login(entrada []string) string{
	var respuesta string
	var user string //obligatorio. Nombre 
	var pass string //obligatorio
	var id string   //obligatorio. Id de la particion en la que quiero iniciar sesion
	Valido := true
	//var pathDico string

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

		//********************  ID *****************
		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])

		//********************  USER *****************
		} else if strings.ToLower(valores[0]) == "user" {
			user = valores[1]

		//******************** PASS *****************
		} else if strings.ToLower(valores[0]) == "pass" {
			pass = valores[1]

		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("LOGIN ERROR: Parametro desconocido: ", valores[0])
			Valido = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	// Se valida que se haya ingresado los diferentes parametros
	if id==""{
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO ID ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO ID "
	}

	if pass==""{
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO PASS ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO PASS "
	}

	if user==""{
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO USER ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO USER "
	}

	if Valido{
		fmt.Println("VALIDACIONES LISTAS")
	}

	return respuesta
}