# Crowdbeat

Provides a [Beat](https://www.elastic.co/products/beats) that processes
audit logs from an [Atlassian Crowd](https://www.atlassian.com/software/crowd) instance.

## Configuration

See [crowdbeat.yml](https://github.com/adamstauffer/crowdbeat/tree/master/crowdbeat.yml) for a sample configuration file.

```yaml
crowdbeat:
  # Defines how after audit logs are polled
  period: 10m
  #crowd_url: https://crowd.company.com/crowd
  #crowd_username: crowdadmin
  #crowd_password: password
```

You'll need to provide the URL to your Crowd instance, along with an
administrator username and password that can access the audit logs.

Audit logs are queried for the time defined in the `period`. This is based on
the current system time, so it is recommended that `Crowdbeat` runs on the
same server as Crowd.

**DISCLAIMER** - The Crowd Audit log API is currently experimental and subject
to change. This beat was written for [Crowd version 3.2.1](https://docs.atlassian.com/atlassian-crowd/3.2.1/REST/#admin/1.0/auditlog-searchAuditLog). Later versions of Crowd may alter the API, which
will lead to incompatibility with `Crowdbeat`.


## Run

To run Crowdbeat with debugging output enabled, run:

```
./crowdbeat -c crowdbeat.yml -e -d "*"
```

## Compiling from scratch

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Crowdbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Crowdbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/adamstauffer/crowdbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Crowdbeat run the command below. This will generate a binary
in the same directory with the name crowdbeat.

```
make
```
