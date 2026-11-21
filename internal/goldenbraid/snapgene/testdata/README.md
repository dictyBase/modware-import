# SnapGene test fixtures

These files are real GoldenBraid assembly maps from `assemblies.zip`. The source is the dictyBase download mirror at
https://dicty.labs.duke.edu/downloads#goldenbraid-plasmid-files. The files keep the dictyBase terms of use.

| File | Bases | DNA flags | Features | Gap segments | Primers |
|------|-------|-----------|----------|--------------|---------|
| `pDGB_A1_act15p_mNeonGreen_act8t.dna` | 3666 | `0x0f` | 8 | 0 | 1 |
| `pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t.dna` | 4932 | `0x0f` | 13 | 0 | 1 |
| `pDGB_A1_act15p_DDB_G0278295_mNeonGreen_act8t.dna` | 4458 | `0x0f` | 11 | 2 | 1 |

Two file names are short forms. The original names are
`pDGB_A1_act15p_GST_GGGSx4Linker_H2Bv3_GGGSx4Linker_mNeonGreen_act8t.dna` and
`pDGB_A1_act15p_DDB_G0278295_GGGSx4Linker_mNeonGreen_act8t.dna`.

The archive `assemblies.zip` holds 156 `.dna` files in one folder, `GoldenBraid Assemblies/`.
All 156 files are circular. No feature segment crosses the origin. 14 files have gap segments.
