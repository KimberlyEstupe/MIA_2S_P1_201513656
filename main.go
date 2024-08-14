package main

import (
	Comandos "MIA_2S_P1_201513656/Comandos"
	DISK 	"MIA_2S_P1_201513656/Comandos/AdminDiscos"
	"bufio"
	"fmt"
	"os"
	"strings"
	
)

//DISK 	"./Comandos/AdminDiscos"

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
			
			DISK.mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK")
		}
	}else if strings.ToLower(parametros[0])=="rmdisk"{
		fmt.Println("rmdisk") 
	}else if strings.ToLower(parametros[0])=="fdisk"{
		fmt.Println("fdisk")
	}else if strings.ToLower(parametros[0])=="rmdisk"{
		fmt.Println("rmdisk")
	}else if strings.ToLower(parametros[0])=="mount"{
		fmt.Println("rmdisk")
	}else if strings.ToLower(parametros[0])=="unmount"{
		fmt.Println("rmdisk")
	}
}