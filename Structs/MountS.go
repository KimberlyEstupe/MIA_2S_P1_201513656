package Structs

var Pmontaje []Mount
type Mount struct{
	MPath  string
	Letter byte
	Cont   int 
}

func AddPathM (path string, L byte, cont int){
	Pmontaje = append(Pmontaje, Mount{MPath: path ,Letter: L,Cont: cont})
}

var Montadas []mountAlready
type mountAlready struct{
	Id		 string
	PathM	 string
}

func AddMontadas(id string, path string){
	Montadas = append(Montadas, mountAlready{Id: id, PathM: path})
}

