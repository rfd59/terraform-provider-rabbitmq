# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.5.0 (April 21, 2025)

FEATURE:

* Add `internal` and `alternate_exchange` fields for `rabbitmq_exchange` resource - @rfavreau

* Modify the default value for `durable` field for `rabbitmq_exchange` resource - @rfavreau

* Add an input check to `type` field for `rabbitmq_exchange` resource and set its optional - @rfavreau

* Force the `rabbitmq_exchange` resource to be recreated when a field is updated - @rfavreau

* Add `internal` and `alternate_exchange` fields for `rabbitmq_exchange` datasource - @rfavreau

BUILD / DEV:

* Update _rabbit-hole_ dependency - @rfavreau

* Refactor `rabbitmq_queue` resource - @rfavreau

* Add _RabbitMQ 4.1_ into unit tests - @rfavreau

* Upgrade _resource_exchange_ unit tests - @rfavreau


## 2.4.0 (March 16, 2025)

FEATURE:

* Add `type` read only field for `rabbitmq_queue` resource - @rfavreau

* Add `rabbitmq_queue` datasource - @rfavreau

FIX:

* `x-queue-type` argument causing always queue recreation [#32](https://github.com/rfd59/terraform-provider-rabbitmq/issues/32) - @pnowy
  ([#35](https://github.com/rfd59/terraform-provider-rabbitmq/pull/35))

* Datasource _vhost_ returns no error message when the vhost doesn't exist [#36](https://github.com/rfd59/terraform-provider-rabbitmq/issues/36) - @rfavreau

BUILD / DEV:

* Update dependencies - @rfavreau
  ([#29](https://github.com/rfd59/terraform-provider-rabbitmq/pull/29), [#31](https://github.com/rfd59/terraform-provider-rabbitmq/pull/31), [#33](https://github.com/rfd59/terraform-provider-rabbitmq/pull/33), [#34](https://github.com/rfd59/terraform-provider-rabbitmq/pull/34))

* Upgrade _resource_queue_ unit tests - @rfavreau

* Fix golangci-lint - @rfavreau

## 2.3.0 (December 06, 2024)

FEATURES:

* Add custom headers to _RabbitMQ_ Api - @Dbzman
  ([#27](https://github.com/rfd59/terraform-provider-rabbitmq/pull/27))

BUILD / DEV:

* Upgrade to Golang 1.23 - @rfavreau

* Update dependencies - @rfavreau
  ([#23](https://github.com/rfd59/terraform-provider-rabbitmq/pull/23), [#25](https://github.com/rfd59/terraform-provider-rabbitmq/pull/25), [#26](https://github.com/rfd59/terraform-provider-rabbitmq/pull/26))


## 2.2.0 (September 29, 2024)

FIX:

* Fix vhost resource - @rfavreau
  ([#16](https://github.com/rfd59/terraform-provider-rabbitmq/pull/16))
  > - Import `default_queue_type` when reading vhost resource ([#14](https://github.com/rfd59/terraform-provider-rabbitmq/pull/14))
  > - Fix `default_queue_type` for _RabbitMQ 3.10_
  > - Validate function for `default_queue_type` attribute
  > - Set default value for `default_queue_type`
  > - Update the Acceptance Tests

BUILD / DEV:

* Update GitHub Actions - @rfavreau
  ([#15](https://github.com/rfd59/terraform-provider-rabbitmq/pull/15))

* GitHub settings - @rfavreau
  ([#18](https://github.com/rfd59/terraform-provider-rabbitmq/pull/18))

## 2.1.0 (May 26, 2024)

FEATURES:

* Manage the user limits [#10](https://github.com/rfd59/terraform-provider-rabbitmq/issues/10) - @rfavreau
  ([#10](https://github.com/rfd59/terraform-provider-rabbitmq/pull/10))

BUILD / DEV:

* New repository structure - @rfavreau
  ([#11](https://github.com/rfd59/terraform-provider-rabbitmq/pull/11))

* Build task for Provider documentation - @rfavreau
  ([#12](https://github.com/rfd59/terraform-provider-rabbitmq/pull/12))

## 2.0.0 (May 20, 2024)

FEATURES:

* Added vhost options/limits and Added shovel parameters - @Galvill
  ([#3](https://github.com/rfd59/terraform-provider-rabbitmq/pull/3))

FIX:

* Add already exists exceptions and enhance errors - @Chahed
  ([#6](https://github.com/rfd59/terraform-provider-rabbitmq/pull/6))

BUILD / DEV:

* Update project - @rfavreau
  ([#2](https://github.com/rfd59/terraform-provider-rabbitmq/pull/2))

* Update GitHub Actions - @rfavreau
  ([#7](https://github.com/rfd59/terraform-provider-rabbitmq/pull/7))

## 1.8.0 (March 15, 2023)

FEATURES:

* Support for RabbitMq versions 3.9 and 3.10 - @mryan43 
  ([#40](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/40))

FIX:

* `rabbitmq_federation`: Fix problematic default `message_ttl` - @ahmadalli
  ([#47](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/47))

* `rabbitmq_binding`: Get only bindings related to source/destination to be faster - @avitsidis 
  ([#43](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/43))

## 1.7.0 (August 20, 2022)

FEATURES:

* `rabbitmq_operator_policy`: new resource - @MrLuje
  ([#8](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/8))

* Add `rabbitmq_vhost`, `rabbitmq_user` and `rabbitmq_exchange` datasources - @Skeen
  ([#37](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/37))

FIX:

* `rabbitmq_shovel`: ForceNew on every parameters - @akurz
  ([#27](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/27))

BUILD / DEV:

* Update Terraform SDK to v2 and Go to 1.19 - @cyrilgdn
  ([#21](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/21))
  ([#38](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/38))

## 1.6.0 (September 01, 2021)

FEATURES:

* Allow configuration of a RabbitMQ-specific proxy - @haines
  ([#16](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/16))

* Update rabbit-hole to 2.10.0 - @MrLuje
  ([#14](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/14))

DOCUMENTATION:

* `provider`: Add `clientcert_file` and `clientkey_file` documentation - @nico2610
  ([#10](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/10))

DEV IMPROVEMENTS:

* Configure Github actions to run acceptance tests - @cyrilgdn
  ([#11](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/11))

* Run golangci-lint in Github actions - @cyrilgdn
  ([#12](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/12))

* Update to go1.15 and remove vendor directory - @cyrilgdn
  ([#13](https://github.com/cyrilgdn/terraform-provider-rabbitmq/pull/13))

## 1.5.1 (November 11, 2020)

FEATURES:

* `rabbitmq_shovel`: Add more parameters and allow to import.
  ([#60](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/60))

DEV IMPROVEMENTS:

* Add goreleaser config
* Pusblish on Terraform registry: https://registry.terraform.io/providers/cyrilgdn/rabbitmq/latest

## 1.5.0:

Replaced by 1.5.1.

## 1.4.0 (July 17, 2020)

FEATURES:

* `rabbitmq_federation_upstream`: New resource to manage federation upstreams.
  ([#55](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/55))

* `rabbitmq_shovel`: New resource to manage shovels.
  ([#48](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/48))

* `provider`: Adding client certificate authentication
  ([#29](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/29))

* `rabbitmq_binding`: Allow to specify arguments directly as JSON with `arguments_json`.
  ([#59](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/59))

DEV IMPROVEMENTS:

* Upgrade rabbithole to v2.2.
  ([#54](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/54)) and ([#57](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/57))

* Remove official support of RabbitMQ 3.6.
  ([#58](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/58))

* Upgrade to Go 1.14

## 1.3.0 (February 23, 2020)

FEATURES:

* New resource: ``rabbitmq_topic_permissions``. This allows to manage permissions on topic exchanges.
  This is compatible with RabbitMQ 3.7 and newer.
  ([#49](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/49))

FIXES:

* ``rabbitmq_queue``: Set ForceNew on all attributes. Queues cannot be changed after creation.
  ([#38](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/38))
  ([#53](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/53))

* ``rabbitmq_permissions``: Fix error when setting empty permissions.
  ([#52](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/52))

IMPROVEMENTS:

* Allow to use the provider behind a proxy.
  It reads HTTPS_PROXY / HTTP_PROXY environment variables to configure the HTTP client (cf [net/http documentation](https://godoc.org/net/http#ProxyFromEnvironment))
  ([#39](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/39))

* Document the configuration of the provider with environment variables.
  ([#50](https://github.com/terraform-providers/terraform-provider-rabbitmq/pull/50))

## 1.2.0 (January 08, 2020)

FIXES:

* rabbitmq_user: Fix tags/password update.
  ([#31](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/31))

* Correctly handle "not found" errors
  ([#45](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/45))

DEV IMPROVEMENTS:

* Upgrade to Go 1.13
  ([#46](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/46))

* Terraform SDK migrated to new standalone Terraform plugin SDK.
  ([#46](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/46))

* Execute acceptance tests in Travis.
  ([#47](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/47))

## 1.1.0 (June 21, 2019)

FIXES:

* Fixed issue preventing policies from updating ([#18](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/18))
* Policy: rename user variable to name ([#19](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/19))
* Fixed `arguments_json` in the queue resource, unfortunately it never worked and failed silently. A queue that receives arguments outside of terraform, where said arguments are not of type string, and was originally configured via `arguments` will be saved to `arguments_json`. This will present the user a diff but avoids a permanent error. ([#26](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/26))

DEV IMPROVEMENTS:

* Upgrade to Go 1.11 ([#23](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/23))
* Provider has been switched to use go modules and bumps the Terraform SDK to v0.11 ([#26](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/26))
* Makefile: add `website` and `website-test` targets ([#15](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/15))
* Upgrade `hashicorp/terraform` to v0.12.2 for latest Terraform 0.12 SDK ([#34](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/34))

## 1.0.0 (April 27, 2018)

IMPROVEMENTS:

* Allow vhost names to contain slashes ([#11](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/11))

FIXES:

* Allow integer values for policy definitions ([#13](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/13))

## 0.2.0 (September 26, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Due to a bug discovered where bindings were not being correctly stored in state, `rabbitmq_bindings.properties_key` is now a read-only, computed field.

IMPROVEMENTS:

* Added `arguments_json` to `rabbitmq_queue`. This argument can accept a nested JSON string which can contain additional settings for the queue. This is useful for queue settings which have non-string values. ([#6](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/6))

FIXES:

* Fix bindings not being saved to state ([#8](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/8))
* Fix issue in `rabbitmq_user` where tags were removed when a password was changed ([#7](https://github.com/terraform-providers/terraform-provider-rabbitmq/issues/7))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
