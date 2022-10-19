package identityzone

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/fields"
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
