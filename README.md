# k8s-unused-secret-detector

Detect unused Kubernetes Secrets

## Build

```
git clone git@github.com:dtan4/k8s-unused-secret-detector.git
cd k8s-unused-secret-detector
make
```

## Usage

```bash
k8s-unused-secret-detector [-n NAMESPACE] [--context CONTEXT]
```

Example: Delete unused Secrets in Namespace `awesome`

```bash
k8s-unused-secret-detector -n awesome | kubectl delete secret -n awesome
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4/))
