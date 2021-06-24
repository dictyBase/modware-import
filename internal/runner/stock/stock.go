package stock

import (
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	logLevel = "info"
)

// LoadAll loads all stock data
func LoadAll() error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	mg.Deps(mg.F(gwdi, bin))
	return nil
}

// LoadStrain load all strain and related data
func LoadStrain() error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	mg.SerialDeps(
		mg.F(strain, bin),
		mg.F(characteristics, bin),
		mg.F(strainProp, bin),
		mg.F(strainSyn, bin),
		mg.F(strainInv, bin),
		mg.F(phenotype, bin),
		mg.F(genotype, bin),
		mg.F(gwdi, bin),
	)
	return nil
}

// Strain loads strain data including curator assignment
func strain(bin string) error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	if err := env.ServiceEnvs(); err != nil {
		return err
	}
	s := runner.TermSpinner("Loading strain data ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "strain")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{
		"-g", "strain_genes.tsv",
		"-i", "strain_strain.tsv",
		"--strain-annotator-input", "strain_user_annotations.csv",
		"-p", "strain_publications.tsv"}...)
	return sh.Run(bin, cmd...)
}

// Plasmid loads plasmid data including curator assignment
func plasmid(bin string) error {
	mg.Deps(mg.F(strainSyn, bin))
	s := runner.TermSpinner("Loading plasmid data ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "plasmid")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{
		"-a", "plasmid_user_annotations.csv",
		"-g", "plasmid_genes.tsv",
		"-i", "plasmid_strain.tsv",
		"-g", "plasmid_genes.tsv"}...)
	return sh.Run(bin, cmd...)
}

// Characteristics loads strain characteristics
func characteristics(bin string) error {
	s := runner.TermSpinner("Loading strain characteristics ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "strainchar")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_characteristics.tsv"}...)
	return sh.Run(bin, cmd...)
}

// StrainProp loads strain property data
func strainProp(bin string) error {
	s := runner.TermSpinner("Loading strain properties ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "strainprop")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_props.tsv"}...)
	return sh.Run(bin, cmd...)
}

// Genotype load strain genotype data
func genotype(bin string) error {
	s := runner.TermSpinner("Loading strain genotype ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "genotype")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_genotype.tsv"}...)
	return sh.Run(bin, cmd...)
}

// StrainSyn loads strain synonym data
func strainSyn(bin string) error {
	s := runner.TermSpinner("Loading strain synonym ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "strainsyn")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_props.tsv"}...)
	return sh.Run(bin, cmd...)
}

// StrainInv loads strain inventory data
func strainInv(bin string) error {
	s := runner.TermSpinner("Loading strain inventory ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "strain-inventory")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_inventory.tsv"}...)
	return sh.Run(bin, cmd...)
}

// Phenotype loads strain phenotype data
func phenotype(bin string) error {
	s := runner.TermSpinner("Loading strain phenotype ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "phenotype")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "strain_phenotype.tsv"}...)
	return sh.Run(bin, cmd...)
}

// PlasmidInv loads plasmid inventory data
func plasmidInv(bin string) error {
	mg.Deps(mg.F(plasmid, bin))
	s := runner.TermSpinner("Loading plasmid inventory ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "plasmid-inventory")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "plasmid-inventory.csv"}...)
	return sh.Run(bin, cmd...)
}

// Gwdi loads GWDI strain mutant data
func gwdi(bin string) error {
	s := runner.TermSpinner("Loading gwdi strain ...")
	defer s.Stop()
	s.Start()
	cmd := append(baseCmd(), "gwdi")
	cmd = append(cmd, minioCmd()...)
	cmd = append(cmd, []string{"-i", "gwdi_strain.csv"}...)
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
