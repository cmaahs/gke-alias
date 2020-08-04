# gke-alias

"gcloud container clusters get-credentials" currently doesn't have an "alias" feature.
This tool is to simplify my scripts to shorten the cluster name.

## Use

`gke-alias get` will display the Current Context with the Alias and Full Cluster Name.

`gke-alias set --alias short-alias` will set the alias of the Current Context.

## KubeConfig

The tool will first look at the first PATH entry in the KUBECONFIG file.
If that does not exist, it will look for the `~/.kube/config` file as default.

## Install

`brew install cmaahs/admin-scripts/gke-alias`

## Other OS/Architectures

I am currently only building x64 architectures, and only MacOS/linux.  Feel free to
clone and build from the source.

[github project](https://github.com/cmaahs/gke-alias)
