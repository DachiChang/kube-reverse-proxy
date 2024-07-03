# Kube reverse proxy

## build and run in bastion

1. in bastion server
    a. cat ~/.kube/config
    b. echo -n CERTIFICATE-AUTHORITY-DATA | base64 -d > authority.cert
    c. echo -n CLIENT-CERTIFICATE-DATA | base64 -d > user.cert
    d. echo -n CLIENT-KEY-DATA | base64 -d > user.key
2. modify targetURI if needed
3. go build .
4. ./run.sh

## Local to connect to k8s control plane via bastion

1. ssh -N -L :4433:127.0.0.1:3000 devops@BASTION_SERVER (enable local 127.0.0.1:4433 to kube-reverse-proxy)
2. vim ~/.kube/config
```yaml=
- cluster:
    certificate-authority-data: XXXXX
    server: 127.0.0.1:4433 (replace to local tunnel to bastion)
  name: pae-uat-aks
```
