# CF-Puppeteer  

*cf plugin for hands-off, zero downtime application deploys*

## notice

This project was forked from [contraband](https://github.com/contraband/autopilot).

It was renamed to *cf-puppeteer* and is being further developed under the new name.

# changelog

To get an overview of the changes between versions, read the [changelog](CHANGELOG.md).

## version

The latest version of *CF-Puppeteer* is *0.0.14*. It works with and is based on Cloud Foundry CLI version 6.43.0.

For more details on the most recent release, check out the [changelog](CHANGELOG.md).

## cf installation

Download the latest version from the [releases][releases] page and make it executable.

```
$ cf install-plugin path/to/downloaded/binary
```

[releases]: https://github.com/happytobi/cf-puppeteer/releases

## usage

```
$ cf zero-downtime-push \
    -f path/to/new_manifest.yml \
    -p path/to/new/path
    -t 120
```

To get more information go to [CF-Puppeteer homepage](https://cf-puppeteer.happytobi.com/)

### passing an application name

To override the application name from the manifest, specify it as command line argument. For example:

```
$ cf zero-downtime-push application-to-replace \
    -f path/to/new_manifest.yml \
    -p path/to/new/path
    -t 120
```

### changing the health check settings

To have more control over the health checks of your application, *CF-Puppeteer* supports additional parameters. For example:

```
$ cf zero-downtime-push application-to-replace \
    -f path/to/new_manifest.yml \
    -p path/to/new/path
    -t 120
    --health-check-type http
    --health-check-http-endpoint /health
    --invocation-timeout 10
```

While *CF-Puppeteer* gives precedence to command line parameters, you can also specify `health-check-type` and `health-check-http-endpoint` in the application manifest. However, Cloud Foundry currently does not support `invocation-timeout` in application manifests. Therefore, if you want to set it, always use the command line.

## method

*CF-Puppeteer* takes a different approach compared to other zero-downtime plugins. It
does not perform any [complex route re-mappings][indiana-jones]. Instead it uses the manifest feature of the Cloud Foundry CLI. The method also has the advantage of treating a manifest as the source of truth and will converge the
state of the system towards that. This makes the plugin ideal for continuous
delivery environments.

1. The old application is renamed to `<APP-NAME>-venerable`. It keeps its old route
   mappings. This change is invisible to users.

2. The new application is pushed to `<APP-NAME>` (assuming that the name has
   not been changed in the manifest). It binds to the same routes as the old
   application (due to them being defined in the manifest) and traffic begins to
   be load-balanced between the two applications.

3. The old application is deleted along with its route mappings. All traffic now goes to the new application.

[indiana-jones]: https://www.youtube.com/watch?v=0gU35Tgtlmg

## local development
for local development you need to install [govendor](https://github.com/kardianos/govendor)
