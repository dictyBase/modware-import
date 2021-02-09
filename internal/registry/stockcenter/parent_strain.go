package stockcenter

import (
	"strings"
	"time"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

const (
	AX3ParentID   = "DBS0237700"
	AX4ParentID   = "DBS0351471"
	ParentSpecies = "Dictyostelium discoideum"
)

func AX3Summary() string {
	var b strings.Builder
	b.WriteString("generic axenic strain, used for curation of parental")
	b.WriteString("strains when the specific AX3 is not available")
	return b.String()
}

func AX4Summary() string {
	var b strings.Builder
	b.WriteString("AX4 strain from the Thompson lab that is the parent of GWDI strains.")
	b.WriteString("A different isolate of that strain is being sequenced")
	return b.String()
}

func AX4ParentStrain() *pb.ExistingStrain {
	return &pb.ExistingStrain{
		Data: &pb.ExistingStrain_Data{
			Type: "strain",
			Id:   AX4ParentID,
			Attributes: &pb.ExistingStrainAttributes{
				CreatedAt:    aphgrpc.TimestampProto(time.Now()),
				UpdatedAt:    aphgrpc.TimestampProto(time.Now()),
				CreatedBy:    DEFAULT_USER,
				UpdatedBy:    DEFAULT_USER,
				Summary:      AX4Summary(),
				Species:      ParentSpecies,
				Label:        "AX4",
				Parent:       AX3ParentID,
				Names:        []string{"AX-4"},
				Publications: []string{"doi:10.1101/582072"},
			},
		},
	}
}

func AX3ParentStrain() *pb.ExistingStrain {
	return &pb.ExistingStrain{
		Data: &pb.ExistingStrain_Data{
			Type: "strain",
			Id:   AX3ParentID,
			Attributes: &pb.ExistingStrainAttributes{
				CreatedAt:    aphgrpc.TimestampProto(time.Now()),
				UpdatedAt:    aphgrpc.TimestampProto(time.Now()),
				CreatedBy:    DEFAULT_USER,
				UpdatedBy:    DEFAULT_USER,
				Summary:      AX3Summary(),
				Species:      ParentSpecies,
				Label:        "AX3",
				Names:        []string{"AX-3"},
				Publications: []string{"5542656"},
			},
		},
	}
}
