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
        "characteristics")
                ${APP} stockcenter \
                    strainchar \
                    --access-key ${ACCESS_KEY} \
                    --secret-key  ${SECRET_KEY} \
                    --log-level debug \
                    -i strain_characteristics.tsv
            ;;
        "*")
                echo unknown command $scmd
            ;;
    esac
done
