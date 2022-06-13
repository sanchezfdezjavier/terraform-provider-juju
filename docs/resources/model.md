---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "juju_model Resource - terraform-provider-juju"
subcategory: ""
description: |-
  A resource that represent a Juju Model.
---

# juju_model (Resource)

A resource that represent a Juju Model.

## Example Usage

```terraform
resource "juju_model" "this" {
  name = "development" # Model name. Required.

  controller = "overlord" # Controller to operate in. Optional
  cloud {                 # Deploy model to different cloud/region to the controller model. Optional
    name   = "aws"
    region = "eu-west-1"
  }

  logging_config = "<root>=INFO" # Specify log levels. Optional.

  config = { # Override default model configuration. Optional.
    development                 = true
    no-proxy                    = "jujucharms.com"
    update-status-hook-interval = "5m"
    # etc...
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name to be assigned to the model.

### Optional

- `cloud` (String) Cloud where the model will operate.
- `cloud_config` (Map of String)
- `cloud_region` (String) Cloud Region where the model will operate.

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String) Type of the model. Set by the Juju's API server

