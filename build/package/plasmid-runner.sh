#!/bin/sh 

APP=$(which app)

for scmd in "$@"
do
    case $scmd in
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
        *)
                echo unknown command $scmd
            ;;
    esac
done
