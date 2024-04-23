# traveltime

## Set-up

To use the application you need to set-up Application Default Credentials (ADC) on Google Cloud.
In short:
```sh
gcloud auth application-default login
```
Read more about the process [here](https://cloud.google.com/docs/authentication/provide-credentials-adc#local-dev).

## Install


```
nix build github:DeterminateSystems/zero-to-nix
```

## Examples

```
traveltime calculate london paris --arriveBy 2024-04-25T09:00:00+01:00
```
