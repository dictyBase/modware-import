# Refactor deleteInventoryIfExists with fp-go

## Summary

Refactor the inventory deletion logic in `goldenbraid_inventory.go` to use fp-go patterns, separating annotation and group deletion into smaller composable functions with proper ordering (annotations first, then groups).

## File to Modify

- `internal/load/stockcenter/goldenbraid_inventory.go`

## New Functions to Add

### 1. `extractAnnotationIDs` (pure)

Extracts all annotation IDs from all groups using `A.FlatMap` and `A.Map`:

```go
func extractAnnotationIDs(gc *pb.TaggedAnnotationGroupCollection) []string {
    return F.Pipe2(
        gc.Data,
        A.FlatMap(func(gcd *pb.TaggedAnnotationGroup) []string {
            return F.Pipe1(
                gcd.Group.Data,
                A.Map(func(anno *pb.TaggedAnnotation_Data) string {
                    return anno.Id
                }),
            )
        }),
        F.Identity[[]string],
    )
}
```

### 2. `extractGroupIDs` (pure)

Extracts all group IDs using `A.Map`:

```go
func extractGroupIDs(gc *pb.TaggedAnnotationGroupCollection) []string {
    return F.Pipe1(
        gc.Data,
        A.Map(func(gcd *pb.TaggedAnnotationGroup) string {
            return gcd.Group.GroupId
        }),
    )
}
```

### 3. `deleteAnnotationIOE` (curried, effectful)

Deletes a single annotation by ID:

```go
func deleteAnnotationIOE(
    client pb.TaggedAnnotationServiceClient,
) func(string) IOE.IOEither[error, string] {
    return func(id string) IOE.IOEither[error, string] {
        return F.Pipe2(
            IOE.TryCatchError(func() (*pb.Empty, error) {
                return client.DeleteAnnotation(
                    context.Background(),
                    &pb.DeleteAnnotationRequest{Id: id, Purge: true},
                )
            }),
            IOE.MapLeft[*pb.Empty](func(err error) error {
                return fmt.Errorf("error deleting annotation %s: %w", id, err)
            }),
            IOE.Map[error](func(_ *pb.Empty) string { return id }),
        )
    }
}
```

### 4. `deleteGroupIOE` (curried, effectful)

Deletes a single annotation group by ID:

```go
func deleteGroupIOE(
    client pb.TaggedAnnotationServiceClient,
) func(string) IOE.IOEither[error, string] {
    return func(groupID string) IOE.IOEither[error, string] {
        return F.Pipe2(
            IOE.TryCatchError(func() (*pb.Empty, error) {
                return client.DeleteAnnotationGroup(
                    context.Background(),
                    &pb.GroupEntryId{GroupId: groupID},
                )
            }),
            IOE.MapLeft[*pb.Empty](func(err error) error {
                return fmt.Errorf("error deleting annotation group %s: %w", groupID, err)
            }),
            IOE.Map[error](func(_ *pb.Empty) string { return groupID }),
        )
    }
}
```

## Updated `deleteInventoryIfExists`

Fully inlined orchestration using point-free style:

```go
func deleteInventoryIfExists(
    ctx WithInventory,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
    client := regsc.GetAnnotationAPIClient()

    return F.Pipe3(
        // Step 1: Extract annotation IDs and delete them
        F.Pipe2(
            ctx.Inventory,
            extractAnnotationIDs,
            A.Map(deleteAnnotationIOE(client)),
        ),
        IOE.SequenceArray[error, string],
        // Step 2: Extract group IDs and delete them
        IOE.Chain(func(_ []string) IOE.IOEither[error, []string] {
            return F.Pipe3(
                ctx.Inventory,
                extractGroupIDs,
                A.Map(deleteGroupIOE(client)),
                IOE.SequenceArray[error, string],
            )
        }),
        // Step 3: Return original collection
        IOE.Map[error](func(_ []string) *pb.TaggedAnnotationGroupCollection {
            return ctx.Inventory
        }),
    )
}
```

## Import to Add

```go
A "github.com/IBM/fp-go/array"
```

## Design Decisions

1. **Ordering**: Annotations deleted first, then groups (clean up contents before container)
2. **Error handling**: Fail-fast via `IOE.SequenceArray`
3. **No special empty case**: Empty slices flow naturally through the pipeline
4. **Location**: All new functions in `goldenbraid_inventory.go` (not shared)
5. **No intermediate batch functions**: Logic inlined directly in `deleteInventoryIfExists`

## Implementation Steps

1. Add `A` import for `github.com/IBM/fp-go/array`
2. Add `extractAnnotationIDs` function
3. Add `extractGroupIDs` function
4. Add `deleteAnnotationIOE` function
5. Add `deleteGroupIOE` function
6. Replace `deleteInventoryIfExists` with new implementation
7. Run tests to verify behavior
