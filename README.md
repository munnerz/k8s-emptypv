# k8s-emptypv

This is a PoC semi-sticky emptyDir flex volume plugin for Kubernetes.

It will create simple hostpath persistent volumes with a volume ID at a specified
path on the hosts filesystem.

It currently supports the following options:

* `volumeID`: a unique identifier for this volume - this will be used as the directory name for the volume
* `path`: the base path to mount under, for example `/var/lib/kubelet/volumes`

As this is a persistent volume, the data will not be deleted until the persistent volume itself is deleted.

See this issue for more information: https://github.com/kubernetes/kubernetes/issues/7562
(specifically, this comment: https://github.com/kubernetes/kubernetes/issues/7562#issuecomment-261795862)