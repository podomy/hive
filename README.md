                              HIVE
Hive is a runtime and control layer for running software across machines
in environments with poor connectivity and conditions that degrade
electronics and communication, such as space, remote terrain,
underground sites, or the sea.

It is designed for mathematical consistency. Each segment can keep
operating from local knowledge, and when segments meet again their state
is reconciled by explicit rules instead of a hidden central truth.

Formal verification:
- The design of important system parts is fully formally verified.
- Selected critical implementation paths are formally verified at the
  code level.

Hive is not standard Kubernetes. It may support Kubernetes-like
deployment workflows, but it does not promise a coherent cluster
network, a central control plane, or one always-current source of truth.
If your software needs one live global truth, Hive is the wrong place to
run it. Cluster segmentation is expected. Local operation and later
reconciliation are part of the model.

Supported IDEs we work with:
https://zed.dev/
We don't support any other IDEs at the moment. We highly
recommend for you to use Zed.
