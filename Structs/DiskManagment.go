package Structs

// NOTA: Recordar que los atributos de los struct deben iniciar con mayuscula
type MBR struct {
	MbrSize    int32        //mbr_tamano
	FechaC     [16]byte     //mbr_fecha_creacion
	Id         int32        //mbr_dsk_signature (random de forma unica)
	Fit        [1]byte      //dsk_fit
	Partitions [4]Partition //mbr_partitions
}

type Partition struct {
	Status      [1]byte  //part_status
	Type        [1]byte  //part_type
	Fit         [1]byte  //part_fit
	Start       int32    //part_start
	Size        int32    //part_s
	Name        [16]byte //part_name
	Correlative int32    //part_correlative
	Id          [4]byte  //part_id
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}
