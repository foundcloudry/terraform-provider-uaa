package identityzone

import (
	"context"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/api"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/fields"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataSource = &schema.Resource{
	Schema:      dataSourceSchema,
	ReadContext: readDataSource,
}

func readDataSource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	izm := session.IdentityZoneManager()

	name := data.Get(fields.Name.String()).(string)

	identityZone, err := izm.FindByName(name)
	if err != nil {
		return diag.FromErr(err)
	}

	MapIdentityZoneToResource(identityZone, data)

	return nil
}
