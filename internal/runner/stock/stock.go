package stock

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	command  = "importer"
	logLevel = "info"
)

// Clean deletes all order data from arangodb database
func Clean() error {
	mg.SerialDeps(
		runner.Build,
		mg.F(runner.CleanDB, "stock"),
	)
	return nil
}

// LoadAll loads all stock data
func LoadAll() error {
	mg.Deps(Gwdi)
	return nil
}

// Strain loads strain data including curator assignment
func Strain() error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	if err := env.ServiceEnvs(); err != nil {
		return err
	}
	mg.Deps(runner.Build)
	s := runner.TermSpinner("Loading strain data ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func Plasmid() error {
	mg.Deps(StrainSyn)
	s := runner.TermSpinner("Loading plasmid data ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func Characteristics() error {
	mg.Deps(Strain)
	s := runner.TermSpinner("Loading strain characteristics ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func StrainProp() error {
	mg.Deps(StrainInv)
	s := runner.TermSpinner("Loading strain properties ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func Genotype() error {
	mg.Deps(Characteristics)
	s := runner.TermSpinner("Loading strain genotype ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func StrainSyn() error {
	mg.Deps(StrainProp)
	s := runner.TermSpinner("Loading strain synonym ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func StrainInv() error {
	mg.Deps(Phenotype)
	s := runner.TermSpinner("Loading strain inventory ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func Phenotype() error {
	mg.Deps(Genotype)
	s := runner.TermSpinner("Loading strain phenotype ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func PlasmidInv() error {
	mg.Deps(Plasmid)
	s := runner.TermSpinner("Loading plasmid inventory ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
func Gwdi() error {
	mg.Deps(PlasmidInv)
	s := runner.TermSpinner("Loading gwdi strain ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", command),
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
