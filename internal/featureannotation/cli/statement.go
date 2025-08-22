package cli

const (

	// GeneProductQuery is the AQL query for fetching gene products
	GeneProductQuery = `
    	FOR locus IN locus_gp
		FOR gp IN gene_product
        		FILTER gp.is_automated == 0
        		FILTER locus.gene_product_no == gp.gene_product_no
        		FILTER locus.locus_no == @feature_id
        		RETURN {
            			gene_product: gp.gene_product,
            			created_by: gp.created_by,
				created_on: gp.date_created
        		}
`
	// AQL query to update documents based on featureprop_id
	updateAQLQuery = `
	FOR idx IN RANGE(0,COUNT(@featureprop_ids) - 1)
		FOR prop IN @@collection
			FILTER prop.featureprop_id == @featureprop_ids[idx]
			UPDATE prop WITH { value: @values[idx] } IN @@collection
	`

	ListActiveGenesQ = `
	LET gene_type_id = FIRST(
    		FOR t IN cvterm 
		FILTER t.name == "gene" 
		LIMIT 1 
		RETURN t.cvterm_id
	)

	LET feat_org_id = FIRST(
    		FOR org IN organism 
        	FILTER org.common_name == "dicty" 
        	LIMIT 1 
        	RETURN org.organism_id
	)

	FOR feat IN feature
    		FOR dbx IN dbxref
        		FILTER feat.type_id == gene_type_id
        		FILTER feat.dbxref_id == dbx.dbxref_id
        		FILTER feat.organism_id == feat_org_id
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
	ListSynonyms = `
	// Step 1: Pre-calculate the constant cvterm IDs to avoid repeated lookups.
	LET gene_type_id = FIRST(
    		FOR t IN cvterm FILTER t.name == "gene" LIMIT 1 RETURN t.cvterm_id
	)
	LET synonym_type_id = FIRST(
		FOR cv IN cv
			FOR t IN cvterm 
				FILTER cv.cv_id == t.cv_id
				FILTER cv.name == 'null'
				FILTER t.name == "synonym" 
				LIMIT 1 
				RETURN t.cvterm_id
	)

	// Step 2: Start from the features that are genes. This is our primary loop.
	FOR feat IN feature
    		FOR dbx IN dbxref
        // Use a subquery to efficiently gather all synonyms for the current feature.
        // With the indexes, this subquery is extremely fast.
        LET synonyms = (
            FOR fsyn IN feature_synonym
                FOR syn IN synonym_
                    FILTER fsyn.feature_id == feat.feature_id 
                    FILTER syn.synonym_id == fsyn.synonym_id 
                    FILTER syn.type_id == synonym_type_id
                    RETURN syn.name
        )

		FILTER feat.is_obsolete == 0
		FILTER feat.is_deleted == 0
		FILTER feat.type_id == gene_type_id

	//Join to the dbxref to get the gene's accession ID.
      	// This join is done once per gene.
        	FILTER feat.dbxref_id == dbx.dbxref_id
      	// Only genes is synonyms
        	FILTER LENGTH(synonyms) > 0
        // Return the final composed document.
		RETURN {
		    gene_id: dbx.accession,
		    name: feat.uniquename,
		    synonyms: synonyms
		}
`
)
