package main

/* Instalar cors:
 go get -u github.com/rs/cors
*Instalar mod:
init go.mud

y agregarlos a las importaciones

*/

import (
	AD "MIA_2S_P1_201513656/Comandos/AdminDiscos"
	Rep "MIA_2S_P1_201513656/Comandos/Rep"
	AA 	"MIA_2S_P1_201513656/Comandos/AdminSisArchivos"
	Us  "MIA_2S_P1_201513656/Comandos/Usuarios"
	PCA	"MIA_2S_P1_201513656/Comandos/AdminPermisos"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type Entrada struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}


func main() {
	//EndPoint 
	http.HandleFunc("/analizar", getCadenaAnalizar)

	// Configurar CORS con opciones predeterminadas
	//Permisos para enviar y recir informacion
	c := cors.Default()

	// Configurar el manejador HTTP con CORS
	handler := c.Handler(http.DefaultServeMux)

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", handler)

}


func getCadenaAnalizar(w http.ResponseWriter, r *http.Request) {
	var respuesta string
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")
	
	var status StatusResponse
	//verificar que sea un post
	if r.Method == http.MethodPost {
		var entrada Entrada
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			status = StatusResponse{Message: "Error al decodificar JSON", Type: "unsucces"}
			json.NewEncoder(w).Encode(status)
			return
		}
		
		//creo un lector de bufer para el archivo
		lector := bufio.NewScanner(strings.NewReader(entrada.Text))
		//leer el archivo linea por linea
		for lector.Scan() {
			//Elimina los saltos de linea
			if lector.Text() != ""{
				//Divido por # para ignorar todo lo que este a la derecha del mismo
				linea := strings.Split(lector.Text(), "#") //lector.Text() retorna la linea leida
				if len(linea[0]) != 0 {
					fmt.Println("\n*********************************************************************************************")
					fmt.Println("Comando en ejecucion: ", linea[0])
					respuesta += "***************************************************************************************************************************\n"
					respuesta += "Comando en ejecucion: " + linea[0] + "\n"
					respuesta += Analizar(linea[0])  + "\n"
				}	
				//Comentarios			
				if len(linea) > 1 && linea[1] != "" {
					fmt.Println("#"+linea[1] +"\n")
					respuesta += "#"+linea[1] +"\n"
				}
			}
			
		}

		//fmt.Println("Cadena recibida ", entrada.Text)
		w.WriteHeader(http.StatusOK)

		status = StatusResponse{Message: respuesta, Type: "succes"}
		json.NewEncoder(w).Encode(status)

	} else {
		//http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		status = StatusResponse{Message: "Metodo no permitido", Type: "unsucces"}
		json.NewEncoder(w).Encode(status)
	}
}


func Analizar(entrada string) (string){
	respuesta := ""
	//Quitar espacios en blanco del final
	tmp := strings.TrimRight(entrada," ")
	//Recibe una linea y la descompone entre el comando y sus parametros
	parametros:= strings.Split(tmp, " -")

	// *============================* ADMINISTRACION DE DISCOS *============================*
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{					
			respuesta = AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			respuesta = "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}

	}else if strings.ToLower(parametros[0])=="rmdisk"{		
		if len(parametros)>1{			
			respuesta = AD.Rmdisk(parametros)
		}else{
			fmt.Println("ERROR RMDISK EN RMDISK, FALTAN PARAMETROS EN RMDISK") 
			respuesta = "ERROR EN RMDISK, FALTAN PARAMETROS EN RMDISK"
		}

	}else if strings.ToLower(parametros[0])=="fdisk"{		
		if len(parametros)>1{			
			respuesta = AD.Fdisk(parametros)
		}else{
			fmt.Println("ERROR EN FDISK, FALTAN PARAMETROS EN FDISK")
			respuesta = "ERROR EN FDISK, FALTAN PARAMETROS EN FDISK"
		}

	}else if strings.ToLower(parametros[0])=="mount"{		
		if len(parametros)>1{			
			respuesta = AD.Mount(parametros)
			
		}else{
			fmt.Println("ERROR EN MOUNT, FALTAN PARAMETROS EN MOUNT")
			respuesta = "ERROR EN MOUNT, FALTAN PARAMETROS EN MOUNT"
		}
	}else if strings.ToLower(parametros[0])=="unmount"{		
		if len(parametros)>1{			
			respuesta = AD.Unmoun(parametros)
			
		}else{
			fmt.Println("ERROR EN MOUNT, FALTAN PARAMETROS EN MOUNT")
			respuesta = "ERROR EN MOUNT, FALTAN PARAMETROS EN MOUNT"
		}
	
	// *===================* ADMINISTRACION DE SISTEMA DE ARCHIVOS *======================*

	}else if strings.ToLower(parametros[0])=="mkfs"{		
		if len(parametros)>1{			
			respuesta = AA.MKfs(parametros)
		}else{
			fmt.Println("ERROR EN MKFS, FALTAN PARAMETROS")
		}
	
	// *===================* ADMINISTRACION DE USUARIOS Y CARPETAS *======================*

	}else if strings.ToLower(parametros[0])=="login"{		
		if len(parametros)>1{			
			respuesta = Us.Login(parametros)
		}else{
			fmt.Println("ERROR EN LOGIN, FALTAN PARAMETROS")
		}
	
	}else if strings.ToLower(parametros[0])=="logout"{		
		respuesta = Us.Logout()

	}else if strings.ToLower(parametros[0])=="mkgrp"{		
		if len(parametros)>1{			
			respuesta = Us.Mkgrp(parametros)
		}else{
			fmt.Println("ERROR EN MKGRP, FALTAN PARAMETROS")
		}
	
	}else if strings.ToLower(parametros[0])=="rmgrp"{		
		if len(parametros)>1{			
			respuesta = Us.Rmgrp(parametros)
		}else{
			fmt.Println("ERROR EN RMGRP, FALTAN PARAMETROS")
		}
	
	}else if strings.ToLower(parametros[0])=="mkusr"{		
		if len(parametros)>1{			
			respuesta = Us.Mkusr(parametros)
		}else{
			fmt.Println("ERROR EN RMGRP, FALTAN PARAMETROS")
		}	

	}else if strings.ToLower(parametros[0])=="rmusr"{		
		if len(parametros)>1{			
			respuesta = Us.Rmusr(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
		}

	}else if strings.ToLower(parametros[0])=="chgrp"{		
		if len(parametros)>1{			
			respuesta = Us.Chgrp(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
		}
	
	// *=======================* PERMISOS DE CARPETAS Y ARHICVOS *============================*
	}else if strings.ToLower(parametros[0])=="mkfile"{		
		if len(parametros)>1{			
			respuesta = PCA.MKfile(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
		}
	
	}else if strings.ToLower(parametros[0])=="cat"{		
		if len(parametros)>1{			
			respuesta = PCA.Cat(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
		}
	
	}else if strings.ToLower(parametros[0])=="mkdir"{		
		if len(parametros)>1{			
			respuesta = PCA.Mkdir(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
		}

	// *============================* OTROS *============================*
	} else if strings.ToLower(parametros[0]) == "rep" {
		//REP
		if len(parametros) > 1 {
			respuesta = Rep.Rep(parametros)
		} else {
			fmt.Println("REP ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
		return ""
	} else {
		fmt.Println("Comando no reconocible")
		return "ERROR: COMANDO NO RECONOCIBLE"
	}

	return respuesta	
}