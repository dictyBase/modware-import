# Justfile for goldenbraid

name := "goldenbraid"
namespace := "dictybase"
github_user := "sba964"
platform := "linux/amd64"

image := namespace + "/" + name
ghcr_image := "ghcr.io/" + image

# Build the docker image for the target platform
build tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{image}}:{{tag}} .

# Build and push the docker image
push tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{image}}:{{tag}} --push .

# Build for GitHub Container Registry
build-ghcr tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{ghcr_image}}:{{tag}} .

# Push to GitHub Container Registry
push-ghcr tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{github_user}} --password-stdin
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{ghcr_image}}:{{tag}} --push .

# List images
list:
    docker images | grep {{image}}
