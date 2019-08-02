#!/bin/sh 

APP=$(which app)
${APP} stockcenter \
    strain \
    --access-key ${ACCESS_KEY} \
    --secret-key  ${SECRET_KEY} \
    --log-level debug \
    -a strain_user_annotations.csv \
    -g strain_genes.tsv \
    -i strain_strain.tsv \
    -p strain_publications.tsv

${APP} stockcenter \
    strainchar \
    --access-key ${ACCESS_KEY} \
    --secret-key  ${SECRET_KEY} \
    --log-level debug \
    -i strain_characteristics.tsv
