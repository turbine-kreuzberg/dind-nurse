# dind-nurse

Docker in docker nurse is a sidecar for dind containers. \
Its task is to keep dind healthy.

## Tasks

### Memory creep

__Issue__

As dind is running longer and longer is uses more and more memory. The memory reserved to the container stays the same. As a result, the memory available for the buildsteps deminishes over time.

__Solution__

Check the memory usage during idle phases and restart dind in case a limit was reached.

### Limit paralelizm

__Issue__

The resources requested from Kubernetes for the dind container are limited. Builds running in parallel need to share these resources. The dind container will OOM terminate on case the dind-daemon together with all parallel build steps use more memory then the limit.

__Solution__

Add memory limits to the build reqeusts and keep track of the current use. Delay forwarding build requests to dind until enough resources are available.

### Garbage collection

__Issue__

Dind is used to keep local state and cache builds. This state grows over time. The used persistent volume will be full at some point. At that point builds will fail.

__Solution__

Check the free space and delete old data until enough space is free.
