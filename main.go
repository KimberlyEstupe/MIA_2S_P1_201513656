package main

import (
	AD "MIA_2S_P1_201513656/Comandos/AdminDiscos"
	//"admindiscos"
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

	// ------------------------------  Administracion de discos---------------------
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{			
			//DISK.mkdisk(parametros)
			AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK")
		}

	}else if strings.ToLower(parametros[0])=="rmdisk"{		
		if len(parametros)>1{			
			fmt.Println("rmdisk")
		}else{
			fmt.Println("rmdisk") 
		}

	}else if strings.ToLower(parametros[0])=="fdisk"{		
		if len(parametros)>1{			
			fmt.Println("fdisk")
		}else{
			fmt.Println("ERROR EN MKDISK")
		}

	}else if strings.ToLower(parametros[0])=="mount"{		
		if len(parametros)>1{			
			fmt.Println("mount")
		}else{
			fmt.Println("ERROR EN MKDISK")
		}

	}else if strings.ToLower(parametros[0])=="unmount"{		
		if len(parametros)>1{			
			fmt.Println("unmount")
		}else{
			fmt.Println("ERROR EN MKDISK")
		}

	}
}