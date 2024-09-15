package adminpermisos

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	TI "MIA_2S_P1_201513656/ToolsInodos"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func MKfile(entrada []string) string{
	respuesta := "Comando mkfile"
	parametrosDesconocidos := false
	var path string
	var cont string	//path del archivo que esta en nuestra maquina y se copiara en el usuario utilizado
	size := 0 		//opcional, si no viene toma valor 0
	r := false
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		fmt.Println("ERROR MKFILE: SESION NO INICIADA")
		respuesta += "ERROR MKFILE: NO HAY SECION INICIADA" + "\n"
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
			//-------------- CONT ---------------------
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
	

	if path ==""{
		fmt.Println("MKFIEL ERROR NO SE INGRESO PARAMETRO PATH")
		return "MKFIEL ERROR NO SE INGRESO PARAMETRO PATH"
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "MKFILE ERROR OPEN FILE "+err.Error()+ "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "MKFILE ERROR READ FILE "+err.Error()+ "\n"
	}

	// Close bin file
	defer Disco.Close()

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
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("MKFILE ERROR. Particion sin formato")
			return "MKFILE ERROR. Particion sin formato"+ "\n"
		}

		//Validar que exista la ruta
		stepPath := strings.Split(path, "/")
		finRuta := len(stepPath) - 1 //es el archivo -> stepPath[finRuta] = archivoNuevo.txt
		idInicial := int32(0)
		idActual := int32(0)
		crear := -1
		//No incluye a finRuta, es decir, se queda en el aterior. EJ: Tamaño=5, finRuta=4. El ultimo que evalua es stepPath[3]
		for i, itemPath := range stepPath[1:finRuta] {
			idActual = TI.BuscarInodo(idInicial, "/"+itemPath, superBloque, Disco)
			//si el actual y el inicial son iguales significa que no existe la carpeta
			if idInicial != idActual {
				idInicial = idActual
			} else {
				crear = i + 1 //porque estoy iniciando desde 1 e i inicia en 0
				break
			}
		}

		//crear carpetas padre si se tiene permiso
		if crear != -1 {
			if r {
				for _, item := range stepPath[crear:finRuta] {
					idInicial = TI.CreaCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), Disco)
					if idInicial == 0 {
						fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
						return "MKFILE ERROR: No se pudo crear carpeta"
					}
				}
			} else {
				fmt.Println("MKDIR ERROR: Carpeta ", stepPath[crear], " no existe. Sin permiso de crear carpetas padre")
				return "MKFILE ERROR: Carpeta "+ stepPath[crear]+ " no existe. Sin permiso de crear carpetas padre"
			}

		}

		//verificar que no exista el archivo (recordar que BuscarInodo busca de la forma /nombreBuscar)
		idNuevo := TI.BuscarInodo(idInicial, "/"+stepPath[finRuta], superBloque, Disco)
		if idNuevo == idInicial {
			fmt.Println("Crear el archivo")
			if cont == "" {
				fmt.Println("No hay cont")
			}else{
				archivoC, err := Herramientas.OpenFile(cont)
				if err != nil {
					return "MKFILE ERROR OPEN FILE "+err.Error()+ "\n"
				}

				//lee el contenido del archivo
				content, err := ioutil.ReadFile(cont)
				if err != nil {
					fmt.Println(err)
					return "ERROR MKFILE "+err.Error()
				}

				fmt.Println(string(content))
			
				// Close bin file
				defer archivoC.Close()
			}
		}else{
			fmt.Println("El archivo ya existe")
			return "El archivo ya existe"
		}
	}
	return respuesta
}

