package group

import (
	"github.com/foundcloudry/terraform-provider-uaa/uaa/group/fields"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var groupSchema = map[string]*schema.Schema{
	fields.DisplayName.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.Description.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.ZoneId.String(): {
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Computed: true,
	},
}
