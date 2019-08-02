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
                --log-level debug \
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
                --log-level debug \
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
                    --log-level debug \
                    -i strain_characteristics.tsv
            ;;
        "strainprop")
                ${APP} stockcenter \
                    strainprop \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level debug \
                    -i strain_props.tsv
            ;;
        "*")
                echo unknown command $scmd
            ;;
    esac
done
