# Universal service broker Redis Driver

A docker image with management api is needed.

Supported image:

https://hub.docker.com/_/redis/

tag: 3.0.6

Docker version > 1.5

If running on stackato, make sure to update docker to 1.8 via `kato path install`

Configuration example: 

```sh
{
    "docker_endpoint": "unix:///var/run/docker.sock",
    "docker_image": "redis",
    "docker_image_version": "3.0.6"
}
```
