package normalize

import (
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

// Generated lenses for each record, and composed lenses that reach a
// nested field from the record that the pipeline carries.
var (
	originL  = MakeOriginLenses()
	cellsL   = MakeCellsLenses()
	rowL     = MakeSourceRowLenses()
	recordL  = MakeRecordLenses()
	diagL    = MakeDiagnosticLenses()
	outcomeL = MakeRowOutcomeLenses()
	batchL   = MakeBatchLenses()
	optionsL = MakeOptionsLenses()
	stateL   = MakeStateLenses()
	ruleL    = MakeruleLenses()
	lineL    = MakelineLenses()
	requestL = MakesheetRequestLenses()

	// State to options.
	stateFolder = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Folder))
	stateSheet  = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Sheet))
	stateOut    = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Out))

	// Source row to origin.
	recordRow = F.Pipe1(recordL.Origin, L.Compose[Record](originL.Row))
	diagFile  = F.Pipe1(diagL.Origin, L.Compose[Diagnostic](originL.File))
	diagRow   = F.Pipe1(diagL.Origin, L.Compose[Diagnostic](originL.Row))
	recFile   = F.Pipe1(recordL.Origin, L.Compose[Record](originL.File))

	// Source row to trimmed cells. Rules compare these values.
	cleanStrain     = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Strain))
	cleanDescriptor = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Descriptor))
	cleanPheno      = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Pheno))
	cleanNote       = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Note))
	cleanAssay      = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Assay))
	cleanAssayTerm  = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.AssayTerm))
	cleanEnv        = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Env))
	cleanEnvTerm    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.EnvTerm))
	cleanReference  = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Reference))
	cleanDeleted    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Deleted))
	cleanCurator    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Curator))

	// Source row to raw cells. Diagnostics report these values.
	rawStrain     = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Strain))
	rawDescriptor = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Descriptor))
	rawPheno      = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Pheno))
	rawAssay      = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Assay))
	rawAssayTerm  = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.AssayTerm))
	rawEnv        = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Env))
	rawEnvTerm    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.EnvTerm))
	rawReference  = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Reference))
	rawDeleted    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Deleted))
	rawCurator    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Curator))
)
