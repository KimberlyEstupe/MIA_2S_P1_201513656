package admindiscos

import (
	"fmt"	
	"strconv"
	"strings"
)


//recibe los parametros de mkdisk
func Mkdisk(entrada []string) {

	var size int	//Obligatorio	
	fit :="FF"		//Puede ser FF, BF, WF
	unit := 1048576	//PUede ser megas(1048576) o kilos (1024)
	Valido := true	//Valida los parametros correcto
	InitSize := false	//Valida el ingreso del parametro size


	/*
	Se recorren todos los parametros
	_ seria el indice, pero se omite. 
	El [1:] indica que se inicializa en el primer parametro de mkdisk
	*/
	for _,parametro :=range entrada[1:]{
		//Quitar espacios en blanco
		tmp := strings.TrimRight(parametro,"")

		//Dividir los parametros entre parametro y valor
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			Valido = false
			break
		}
		
		//********************  SIZE *****************
		if strings.ToLower(valores[0])=="size"{
			InitSize = true
			var err error
			size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
			//if err != nil || size <= 0 { //Se manejaria como un solo error
			if err != nil {
				fmt.Println("MKDISK Error: -size debe ser un valor numerico. se leyo ", valores[1])
				Valido = false
				break
			} else if size <= 0 { //se valida que sea mayor a 0 (positivo)
				fmt.Println("MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", valores[1])
				Valido = false
				break
			}

		

		//********************  Fit *****************
		}else if strings.ToLower(valores[0])=="fit"{
			fmt.Println("Fit: ", valores[1])
		
		//*************** UNIT ***********************
		} else if strings.ToLower(valores[0]) == "unit" {
			//si la unidad es k
			if strings.ToLower(valores[1]) == "k" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1024
				//si la unidad no es k ni m es error (si fuera m toma el valor con el que se inicializo unit al inicio del metodo)
			} else if strings.ToLower(valores[1]) != "m" {
				fmt.Println("MKDISK Error en -unit. Valores aceptados: k, m. ingreso: ", valores[1])
				Valido = false
				break
			}
		//******************* PATH *************
		} else if strings.ToLower(valores[0]) == "path" {
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", valores[0])
			Valido = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	fmt.Println(fit,", ",unit)

	

	if Valido{
		if InitSize{
			fmt.Println("Esta inicializado")
		}else{
			fmt.Println("ERROR: Debe ingresar el parametro Size")
		}
	}

	
}