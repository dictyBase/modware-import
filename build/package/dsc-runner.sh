#!/bin/sh 

APP=$(which app)

for scmd in "$@"
do
    case $scmd in
        "strain")
            ${APP} stockcenter \
                strain \
                --access-key ${ACCESS_KEY} \
                --secret-key  ${SECRET_KEY} \
                --log-level ${LOG_LEVEL} \
                -a strain_user_annotations.csv \
                -g strain_genes.tsv \
                -i strain_strain.tsv \
                -p strain_publications.tsv
            ;;
        "plasmid")
            ${APP} stockcenter \
                plasmid \
                --access-key ${ACCESS_KEY} \
                --secret-key  ${SECRET_KEY} \
                --log-level ${LOG_LEVEL} \
                -a plasmid_user_annotations.csv \
                -g plasmid_genes.tsv \
                -i plasmid_plasmid.tsv \
                -p plasmid_publications.tsv
            ;;
        "characteristics")
                ${APP} stockcenter \
                    strainchar \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_characteristics.tsv
            ;;
        "strainprop")
                ${APP} stockcenter \
                    strainprop \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_props.tsv
            ;;
        "genotype")
                ${APP} stockcenter \
                    genotype \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_genotype.tsv
            ;;
        "strainsyn")
                ${APP} stockcenter \
                    strainsyn \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_props.tsv
            ;;
        "straininv")
                ${APP} stockcenter \
                    strain-inventory \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_inventory.tsv
            ;;
        "phenotype")
                ${APP} stockcenter \
                    phenotype \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i strain_phenotype.tsv \
            ;;
        "gwdi")
                ${APP} stockcenter \
                    gwdi \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -c 5 \
                    -i gwdi_strain_lite.csv \
            ;;
        "plasmidinv")
                ${APP} stockcenter \
                    plasmid-inventory \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level ${LOG_LEVEL} \
                    -i plasmid_inventory.tsv \
            ;;
        "*")
                echo unknown command ${scmd}
            ;;
    esac
done
