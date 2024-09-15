package adminpermisos

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"fmt"
	"strconv"
	"strings"
)

func MKfile(entrada []string) string{
	respuesta := "Comando mkfile"
	parametrosDesconocidos := false
	var path string
	var cont string
	size := 0 //opcional, si no viene toma valor 0
	r := false
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		fmt.Println("ERROR CAT: SESION NO INICIADA")
		respuesta += "ERROR CAT: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores) == 2{
			// --------------- PAHT ------------------
			if strings.ToLower(valores[0]) == "path"  {
				path = strings.ReplaceAll(valores[1],"\"","")
			//-------------- SIZE ---------------------
			}else if strings.ToLower(valores[0]) == "size"{
				//convierto a tipo int
				var err error
				size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
				if err != nil {
					fmt.Println("MKFILE Error: Size solo acepta valores enteros. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores enteros. Ingreso: " + valores[1]
				}

				//valido que sea mayor a 0
				if size < 0 {
					fmt.Println("MKFILE Error: Size solo acepta valores positivos. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores positivos. Ingreso: "+ valores[1]
				}
			}else if strings.ToLower(valores[0]) == "cont"{
				cont = strings.ReplaceAll(valores[1], "\"", "")
			}else{
				parametrosDesconocidos = true
			}
		}else if len(valores) == 1{
			if strings.ToLower(valores[0]) == "r"{
				r = true
			}else{
				parametrosDesconocidos = true
			}
		}else{
			parametrosDesconocidos = true
		}

		if parametrosDesconocidos{
			fmt.Println("MKFILE Error: Parametro desconocido: ", valores[0])
			respuesta += "MKFILE Error: Parametro desconocido: "+ valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}
	fmt.Println(cont," ", r)

	if path ==""{
		fmt.Println("MKFIEL ERROR NO SE INGRESO PARAMETRO PATH")
		return "MKFIEL ERROR NO SE INGRESO PARAMETRO PATH"
	}

	//Abrimos el disco
	file, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "CAR ERROR OPEN FILE "+err.Error()+ "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return "CAR ERROR READ FILE "+err.Error()+ "\n"
	}

	// Close bin file
	defer file.Close()

	//Encontrar la particion correcta
	agregar := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {		
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			agregar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if agregar{
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("CAT ERROR. Particion sin formato")
			return "CAT ERROR. Particion sin formato"+ "\n"
		}
	}
	
	return respuesta
}