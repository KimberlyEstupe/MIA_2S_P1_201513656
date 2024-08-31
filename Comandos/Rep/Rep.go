package rep

import (
	"MIA_2S_P1_201513656/Herramientas"
	"MIA_2S_P1_201513656/Structs"
	"fmt"
	"path/filepath"
	"strings"
)

func Rep(entrada []string) string{
	var respuesta string
	var name string //obligatorio Nombre del reporte a generar
	var path string //obligatorio Nombre que tendrá el reporte
	var id string   //obligatorio sera el del disco o el de la particion
	//var rutaFile string	//nombre del archivo o carpeta reporte file/IS
	Valido := true 

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro,"")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			break
		}

		if strings.ToLower(valores[0]) == "name" {
			name = strings.ToLower(valores[1])
		} else if strings.ToLower(valores[0]) == "path" {
			// Eliminar comillas
			path = strings.ReplaceAll(valores[1], "\"", "")
		} else if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "ruta" {
			//rutaFile = strings.ToLower(tmp[1])
		} else {
			fmt.Println("REP Error: Parametro desconocido: ", valores[0])
			respuesta+="REP Error: Parametro desconocido: " + valores[0]
			Valido = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if Valido{
		if name != "" && id != "" && path != "" {			
			switch name{
			case "mbr":
				fmt.Println("reporte mbr")
				respuesta+= Rmbr(path, id)
			case "disk":
				fmt.Println("reporte disk")
				respuesta+= disk(path, id)
			default:
				fmt.Println("REP Error: Reporte ", name, " desconocido")
				respuesta+="REP Error: Reporte "+ name+" desconocido"
			}
		}else{
			fmt.Println("REP Error: Faltan parametros")
			respuesta+= "REP Error: Faltan parametros"
		}
	}
	return respuesta
}

//---------------------- MBR ---------------------
func Rmbr (path string, id string) string{
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			fmt.Println("Encotrado ", montado.PathM)
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]	
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP MBR Open "+ err.Error()		
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			Respuesta += "ERROR REP MBR Read "+ err.Error()		
		}

		// Close bin file
		defer file.Close()

		//Crea reporte
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MbrSize)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.FechaC[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.RepGraphviz(mbr, file)
		cad += "</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)

		Respuesta += "Reporte de MBR/EBR ejecutado"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}

	
	return Respuesta
}


//---------------- DISK -------------------------
func disk(path string, id string)string{
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			fmt.Println("Encotrado ", montado.PathM)
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]	
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP DISK Open "+ err.Error()	
			return Respuesta	
		}

		var TempMBR Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
			Respuesta += "ERROR REP READ Open "+ err.Error()
			return Respuesta	
		}

		defer file.Close()

		//inicia contenido del reporte graphviz del disco
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr> \n"
		cad += " <td bgcolor='SlateBlue'  ROWSPAN='3'> MBR </td>\n"
		cad += Structs.RepDiskGraphviz(TempMBR, file)
		cad += "\n</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		fmt.Println("RP ", rutaReporte," name ",nombre)

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte de Disk ejecutado"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}
	
	return Respuesta

}