# Intro to Containers

## Tools

- [docker](https://www.docker.com/get-started/)
- [dive](https://github.com/wagoodman/dive)

## Command Snippets

```bash
# Building a container
docker build . -t name:tag-version

# Inspect the image
docker inspect <image>

# View layers and image fs
dive <image>

# Running a container
docker run [args] <image> [commands]
```
