package user

import (
	"context"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/api"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/user/fields"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataSource = &schema.Resource{
	Schema:      dataSourceSchema,
	ReadContext: readDataSource,
}

var dataSourceSchema = map[string]*schema.Schema{
	fields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.ZoneId.String(): {
		Type:     schema.TypeString,
		ForceNew: true,
		Optional: true,
		Computed: true,
	},
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	um := session.UserManager()
	name := data.Get(fields.Name.String()).(string)
	zoneId := data.Get(fields.ZoneId.String()).(string)

	user, err := um.FindByUsername(name, zoneId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(user.Id)

	return nil
}
