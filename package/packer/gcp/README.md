# Bhojpur OS - GCP packer

## Setup

Configure a Google Compute Engine service account (`account.json`):
- https://www.packer.io/docs/builders/googlecompute/#running-without-a-compute-engine-service-account

Configure `${GCP_PROJECT_ID}`.

## Build AMD64

```shell script
packer build template.json
```
