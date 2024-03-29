---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maas_device Resource - terraform-provider-maas"
subcategory: ""
description: |-
  Provides a resource to manage MAAS devices.
---

# maas_device (Resource)

Provides a resource to manage MAAS devices.

## Example Usage

```terraform
resource "maas_dns_domain" "test_domain" {
  name          = "domain"
  ttl           = 3600
  authoritative = true
}

resource "maas_device" "test_device" {
  description = "Test description"
  domain      = maas_dns_domain.test_domain.name
  hostname    = "test-device"
  zone        = "default"
  network_interfaces {
    mac_address = "12:23:45:67:89:de"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_interfaces` (Block Set, Min: 1) A set of network interfaces attached to the device. (see [below for nested schema](#nestedblock--network_interfaces))

### Optional

- `description` (String) The description of the device.
- `domain` (String) The domain of the device.
- `hostname` (String) The device hostname.
- `zone` (String) The zone of the device.

### Read-Only

- `fqdn` (String) The device FQDN.
- `id` (String) The ID of this resource.
- `ip_addresses` (Set of String) A set of IP addressed assigned to the device.
- `owner` (String) The owner of the device.

<a id="nestedblock--network_interfaces"></a>
### Nested Schema for `network_interfaces`

Required:

- `mac_address` (String) MAC address of the network interface.

Read-Only:

- `id` (Number) The id of the network interface.
- `name` (String) The name of the network interface.

## Import

Import is supported using the following syntax:

```shell
# Devices can be imported using their ID or hostname. e.g.
$ terraform import maas_device.test_device test-device
```
