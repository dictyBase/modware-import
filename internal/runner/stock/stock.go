package stock

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	logLevel = "info"
)

// LoadPlasmid load all plasmid and related data
func LoadPlasmid() error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	if err := env.CheckWithoutDB(); err != nil {
		return fmt.Errorf("error in checking for env vars %s", err)
	}
	mg.SerialDeps(
		mg.F(plasmid, bin),
		mg.F(plasmidInv, bin),
	)
	return nil
}

// LoadStrain load all strain and related data
func LoadStrain() error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	if err := env.CheckWithoutDB(); err != nil {
		return fmt.Errorf("error in checking for env vars %s", err)
	}
	mg.SerialDeps(
		mg.F(strain, bin),
		mg.F(characteristics, bin),
		mg.F(strainProp, bin),
		mg.F(strainSyn, bin),
		mg.F(strainInv, bin),
		mg.F(phenotype, bin),
		mg.F(genotype, bin),
		mg.F(Gwdi, bin),
	)
	return nil
}

// Strain loads strain data including curator assignment
func strain(bin string) error {
	runner.ConsoleLog("Loading strain data ...")
	defer runner.ConsoleLog("Done loading strain data ...")
	cmd := collection.Extend(append(baseCmd(), "strain"), minioCmd(), []string{
		"-g", "strain_genes.tsv",
		"-i", "strain_strain.tsv",
		"--strain-annotator-input", "strain_user_annotations.csv",
		"-p", "strain_publications.tsv"})
	return sh.Run(bin, cmd...)
}

// Plasmid loads plasmid data including curator assignment
func plasmid(bin string) error {
	runner.ConsoleLog("Loading plasmid data ...")
	defer runner.ConsoleLog("Done loading plasmid data ...")
	cmd := collection.Extend(
		append(baseCmd(), "plasmid"),
		minioCmd(), []string{
			"-a", "plasmid_user_annotations.csv",
			"-p", "plasmid_publications.tsv",
			"-i", "plasmid_plasmid.tsv",
			"-g", "plasmid_genes.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// PlasmidInv loads plasmid inventory data
func plasmidInv(bin string) error {
	runner.ConsoleLog("Loading plasmid inventory ...")
	defer runner.ConsoleLog("Done loading plasmid inventory ...")
	cmd := collection.Extend(
		append(baseCmd(), "plasmid-inventory"),
		minioCmd(), []string{"-i", "plasmid_inventory.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// Characteristics loads strain characteristics
func characteristics(bin string) error {
	runner.ConsoleLog("Loading strain characteristics ...")
	defer runner.ConsoleLog("Done loading strain characteristics ...")
	cmd := collection.Extend(
		append(baseCmd(), "strainchar"),
		minioCmd(), []string{"-i", "strain_characteristics.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// StrainProp loads strain property data
func strainProp(bin string) error {
	runner.ConsoleLog("Loading strain properties ...")
	defer runner.ConsoleLog("Done loading strain properties ....")
	cmd := collection.Extend(
		append(baseCmd(), "strainprop"),
		minioCmd(),
		[]string{"-i", "strain_props.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// Genotype load strain genotype data
func genotype(bin string) error {
	runner.ConsoleLog("Loading strain genotype ...")
	defer runner.ConsoleLog("Done loading strain genotype ...")
	cmd := collection.Extend(
		append(baseCmd(), "genotype"),
		minioCmd(),
		[]string{"-i", "strain_genotype.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// StrainSyn loads strain synonym data
func strainSyn(bin string) error {
	runner.ConsoleLog("Loading strain synonym ...")
	defer runner.ConsoleLog("Done loading strain synonym ...")
	cmd := collection.Extend(
		append(baseCmd(), "strainsyn"),
		minioCmd(),
		[]string{"-i", "strain_props.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// StrainInv loads strain inventory data
func strainInv(bin string) error {
	runner.ConsoleLog("Loading strain inventory ...")
	defer runner.ConsoleLog("Done loading strain inventory ...")
	cmd := collection.Extend(
		append(baseCmd(), "strain-inventory"),
		minioCmd(),
		[]string{"-i", "strain_inventory.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// Phenotype loads strain phenotype data
func phenotype(bin string) error {
	runner.ConsoleLog("Loading strain phenotype ...")
	defer runner.ConsoleLog("Done loading strain phenotype ...")
	cmd := collection.Extend(
		append(baseCmd(), "phenotype"),
		minioCmd(),
		[]string{"-i", "strain_phenotype.tsv"},
	)
	return sh.Run(bin, cmd...)
}

// Gwdi loads GWDI strain mutant data
func Gwdi(bin string) error {
	runner.ConsoleLog("Loading gwdi strain ...")
	defer runner.ConsoleLog("Done loading gwdi strain ...")
	cmd := collection.Extend(
		append(baseCmd(), "gwdi"),
		minioCmd(),
		[]string{"-i", "gwdi_strain.csv"},
	)
	return sh.Run(bin, cmd...)
}

func baseCmd() []string {
	return []string{
		"stockcenter",
		"--log-level", logLevel,
	}
}

func minioCmd() []string {
	return []string{
		"--access-key", env.MinioAccessKey(),
		"--secret-key", env.MinioSecretKey(),
		"--s3-bucket-path", "import/data/stockcenter",
	}
}
