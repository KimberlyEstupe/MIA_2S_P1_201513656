package main

import (
	AD "MIA_2S_P1_201513656/Comandos/AdminDiscos"
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {
	reader :=bufio.NewScanner(os.Stdin)

	fmt.Println("Ingresar comando: ")
	reader.Scan()

	Analizar(reader.Text())
}

func Analizar(entrada string){
	//Recibe una linea y la descompone entre el comando y sus parametros
	parametros:= strings.Split(entrada, " -")

	// *============================* ADMINISTRACION DE DISCOS *============================*
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{			
			//DISK.mkdisk(parametros)
			AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
		}

	}else if strings.ToLower(parametros[0])=="rmdisk"{		
		if len(parametros)>1{			
			AD.Rmdisk(parametros)
		}else{
			fmt.Println("ERROR RMDISK EN RMDISK, FALTAN PARAMETROS EN RMDISK") 
		}

	}else if strings.ToLower(parametros[0])=="fdisk"{		
		if len(parametros)>1{			
			AD.Fdisk(parametros)
		}else{
			fmt.Println("ERROR EN FDISK, FALTAN PARAMETROS EN FDISK")
		}

	}else if strings.ToLower(parametros[0])=="mount"{		
		if len(parametros)>1{			
			AD.Mount(parametros)
		}else{
			fmt.Println("ERROR EN MOUNT, FALTAN PARAMETROS EN MOUNT")
		}

	}else if strings.ToLower(parametros[0])=="unmount"{		
		if len(parametros)>1{			
			fmt.Println("unmount")
		}else{
			fmt.Println("ERROR EN UNMOUNT, FALTAN PARAMETROS EN UNMOUNT")
		}

	// *============================* OTROS *============================*

	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
	} else {
		fmt.Println("Comando no reconocible")
	}
	
}