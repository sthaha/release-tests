Pipeline Spec
=============

Specifications are defined via H1 tag, you could either use the syntax above or "# Pipeline Spec"

Any context step defined under spec section, gets executed before every scenario. Every unordered list is a step.

* Operator should be installed

Run output pipeline 
-------------------
tags: e2e, integration
 Define a pipeline which has 2 Tasks
 1. To create file 
 2. To read file content created by above task

* Create Namespace
* Create output pipeline
* Run pipeline
* Validate pipelinerun for success status

________________________
These are teardown steps

* Delete Namespace