---
page_title: "Cloud Foundry UAA: uaa_group"
---

# Group Data Source

Gets information on a Cloud Foundry UAA group.

## Example Usage

The following example looks up a group named 'mygroup'.

```
data uaa_group "mygroup" {
    display_name = "mygroup"    
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The name of the group to look up
* `zone_id` - (Optional) The identity zone that the group belongs to

## Attributes Reference

The following attributes are exported:

* `id` - The GUID of the group
* `display_name` - The name of the group
* `description` - The description of the group
