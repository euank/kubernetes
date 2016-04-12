# Rktnetes Jenkins Tests

This directory contains jenkins-job-builder formatted yaml files describing the jenkens jobs we run to test Kubernetes with the rkt container runtime.

## Jenkins configuration

Unfortunately, the job-builder plugin can only manage jobs, not large configuration which might be necessary. Below, the process I used to configure a jenkins master and agents is described. This might not match your process, but should provide a guide for what plugins you need (!important!), etc.


### Master

<!-- TODO this should be formatted more nicely -->

Jenkins version: 2.0-rc-1

'configuration' changes: Shell executable -> `/bin/bash`

Plugins (note, plugins in the default set are not listed; it's assumed you have them. Feel free to improve this list):

* [conditional-buildstep](https://wiki.jenkins-ci.org/display/JENKINS/Conditional+BuildStep+Plugin) - v1.7
* [Environment Injector Plugin](https://wiki.jenkins-ci.org/display/JENKINS/EnvInject+Plugin) - v1.92.1
* [Flaky Test Handler plugin](https://wiki.jenkins-ci.org/display/JENKINS/Flaky+Test+Handler+Plugin) - v1.0.4
* [Flexible Publish Plugin](http://wiki.jenkins-ci.org/display/JENKINS/Flexible+Publish+Plugin) - v0.15.2
* [Post-Build Script Plug-in](http://wiki.jenkins-ci.org/display/JENKINS/PostBuildScript+Plugin) - v0.17
* [Google Cloud Storage plugin](https://wiki.jenkins-ci.org/display/JENKINS/Google+Cloud+Storage+Plugin) - v0.10
* [Log Parser Plugin](https://wiki.jenkins-ci.org/display/JENKINS/Log+Parser+Plugin) - v2.0

You'll also have to setup some agents. That one's up to you
