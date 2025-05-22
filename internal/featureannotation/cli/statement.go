package cli

const (
	// AQL query to update documents based on featureprop_id
	updateAQLQuery = `
	FOR idx IN RANGE(0,COUNT(@featureprop_ids) - 1)
		FOR prop IN @@collection
			FILTER prop.featureprop_id == @featureprop_ids[idx]
			UPDATE prop WITH { value: @values[idx] } IN @@collection
	`

	ListActiveGenesQ = `
	FOR cvt IN cvterm 
		FOR feat IN feature 
			FOR dbx IN dbxref
				FILTER cvt.cvterm_id == feat.type_id
				FILTER feat.dbxref_id == dbx.dbxref_id
				FILTER cvt.name == 'gene'
				FILTER feat.is_obsolete == 0
				RETURN {
					name: feat.uniquename,
					gene_id: dbx.accession,
					feature_id: feat.feature_id,
					created_by: feat.created_by
				}
`
	ListPubmedsByFeature = `
		FOR fpub IN feature_pub
			FOR pb IN pub 
				FILTER fpub.pub_id == pb.pub_id
				FILTER pb.pubplace == 'PUBMED'
				FILTER fpub.feature_id == @feature_id
				RETURN pb.uniquename
	`
)
