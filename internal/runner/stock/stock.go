package stock

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	logLevel = "info"
)

// LoadAll loads all stock data
func LoadAll(branch string) error {
	mg.Deps(mg.F(gwdi, branch))
	return nil
}

// Strain loads strain data including curator assignment
func strain(branch string) error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	if err := env.ServiceEnvs(); err != nil {
		return err
	}
	mg.Deps(mg.F(runner.BuildBranch, branch))
	s := runner.TermSpinner("Loading strain data ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"strain",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-a", "strain_user_annotations.csv",
		"-g", "strain_genes.tsv",
		"-i", "strain_strain.tsv",
		"-p", "strain_publications.tsv",
	)
}

// Plasmid loads plasmid data including curator assignment
func plasmid(branch string) error {
	mg.Deps(mg.F(strainSyn, branch))
	s := runner.TermSpinner("Loading plasmid data ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"plasmid",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-a", "plasmid_user_annotations.csv",
		"-g", "plasmid_genes.tsv",
		"-i", "plasmid_strain.tsv",
		"-p", "plasmid_publications.tsv",
	)
}

// Characteristics loads strain characteristics
func characteristics(branch string) error {
	mg.Deps(mg.F(strain, branch))
	s := runner.TermSpinner("Loading strain characteristics ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"strainchar",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_characteristics.tsv",
	)
}

// StrainProp loads strain property data
func strainProp(branch string) error {
	mg.Deps(mg.F(strainInv, branch))
	s := runner.TermSpinner("Loading strain properties ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"strainprop",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_props.tsv",
	)
}

// Genotype load strain genotype data
func genotype(branch string) error {
	mg.Deps(mg.F(characteristics, branch))
	s := runner.TermSpinner("Loading strain genotype ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"genotype",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_genotype.tsv",
	)
}

// StrainSyn loads strain synonym data
func strainSyn(branch string) error {
	mg.Deps(mg.F(strainProp, branch))
	s := runner.TermSpinner("Loading strain synonym ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"strainsyn",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_props.tsv",
	)
}

// StrainInv loads strain inventory data
func strainInv(branch string) error {
	mg.Deps(mg.F(phenotype, branch))
	s := runner.TermSpinner("Loading strain inventory ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"strain-inventory",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_inventory.tsv",
	)
}

// Phenotype loads strain phenotype data
func phenotype(branch string) error {
	mg.Deps(mg.F(genotype, branch))
	s := runner.TermSpinner("Loading strain phenotype ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"phenotype",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "strain_phenotype.tsv",
	)
}

// PlasmidInv loads plasmid inventory data
func plasmidInv(branch string) error {
	mg.Deps(mg.F(plasmid, branch))
	s := runner.TermSpinner("Loading plasmid inventory ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"plasmid-inventory",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "plasmid-inventory.tsv",
	)
}

// Gwdi loads GWDI strain mutant data
func gwdi(branch string) error {
	mg.Deps(mg.F(plasmidInv, branch))
	s := runner.TermSpinner("Loading gwdi strain ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		logLevel,
		"stockcenter",
		"gwdi",
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"-i", "gwdi_strain.csv",
	)
}
