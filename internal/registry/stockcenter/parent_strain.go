package stockcenter

import (
	"time"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

const (
	AX3ParentId   = "DBS0237700"
	AX4ParentId   = "DBS0351471"
	ParentSpecies = "Dictyostelium discoideum"
	AX3summary    = "generic axenic strain, used for curation of parental strains when the specific AX3 is not available"
	AX4summary    = "AX4 strain from the Thompson lab that is the parent of GWDI strains. A different isolate of that strain is being sequenced"
)

func AX4ParentStrain() *pb.ExistingStrain {
	return &pb.ExistingStrain{
		Data: &pb.ExistingStrain_Data{
			Type: "strain",
			Id:   AX4ParentId,
			Attributes: &pb.ExistingStrainAttributes{
				CreatedAt:    aphgrpc.TimestampProto(time.Now()),
				UpdatedAt:    aphgrpc.TimestampProto(time.Now()),
				CreatedBy:    DEFAULT_USER,
				UpdatedBy:    DEFAULT_USER,
				Summary:      AX4summary,
				Species:      ParentSpecies,
				Label:        "AX4",
				Publications: []string{"doi:10.1101/582072"},
			},
		},
	}
}

func AX3ParentStrain() *pb.ExistingStrain {
	return &pb.ExistingStrain{
		Data: &pb.ExistingStrain_Data{
			Type: "strain",
			Id:   AX3ParentId,
			Attributes: &pb.ExistingStrainAttributes{
				CreatedAt:    aphgrpc.TimestampProto(time.Now()),
				UpdatedAt:    aphgrpc.TimestampProto(time.Now()),
				CreatedBy:    DEFAULT_USER,
				UpdatedBy:    DEFAULT_USER,
				Summary:      AX3summary,
				Species:      ParentSpecies,
				Label:        "AX3",
				Publications: []string{"5542656"},
			},
		},
	}
}
