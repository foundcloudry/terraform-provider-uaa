---
page_title: "Cloud Foundry UAA: uaa_user"
---

# User Data Source

Gets information on a Cloud Foundry UAA user.

## Example Usage

The following example looks up an user named 'myuser'.

```
data "uaa_user" "myuser" {
    name = "myuser"    
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user to look up
* `zone_id` - (Optional) The identity zone that the user belongs to

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the user
