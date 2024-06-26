# traveltime

## Set-up

To use the application you need to set-up Application Default Credentials (ADC) on Google Cloud.
In short:
```sh
gcloud auth application-default login
```
Read more about the process [here](https://cloud.google.com/docs/authentication/provide-credentials-adc#local-dev).

`gcloud` is available in the nix development shell via:
```sh
nix develop -c gcloud auth application-default login
```

## Install

```sh
nix run github:toffernator/traveltime
```

Alternatively, build from source and place the executable in your `$PATH`:

```sh
go build 
mv ./traveltime <SOMEWHERE/IN/YOUR/$PATH>
```

## Examples

```sh
traveltime calculate london paris --arriveBy 2024-04-25T09:00:00+01:00
```

