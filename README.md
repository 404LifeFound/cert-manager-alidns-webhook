# cert-manager-alidns-webhook
support alicloud default credential auth chain, like AK,ENV,OIDC(RRSA in ACK)

# configuration

```
change values.yaml in helm chart for
groupName: acme.yourcompany.com
extraEnvs:
  GROUP_NAME: acme.yourcompany.com
  ALIDNS_REGION: ap-northeast-1

```
`

```
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    email: your_email@email.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key-stag
    solvers:
    - dns01:
        webhook:
          groupName: acme.yourcompany.com
          solverName: alidns
```
