package identityzone

import (
	"context"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Resource = &schema.Resource{
	Schema:        identityZoneSchema,
	CreateContext: createResource,
	ReadContext:   readResource,
	UpdateContext: updateResource,
	DeleteContext: deleteResource,
}

func createResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	izm := session.IdentityZoneManager()

	identityZone := MapResourceToIdentityZone(data)
	response, err := izm.Create(identityZone)
	if err != nil {
		return diag.FromErr(err)
	}

	MapIdentityZoneToResource(response, data)

	return nil
}

func readResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	izm := session.IdentityZoneManager()

	response, err := izm.FindById(data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	MapIdentityZoneToResource(response, data)

	return nil
}

func updateResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	izm := session.IdentityZoneManager()

	identityZone := MapResourceToIdentityZone(data)
	response, err := izm.Update(data.Id(), identityZone)
	if err != nil {
		return diag.FromErr(err)
	}

	MapIdentityZoneToResource(response, data)

	return nil
}

func deleteResource(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {

	session := i.(*api.Session)
	if session == nil {
		return diag.Errorf("client is nil")
	}

	izm := session.IdentityZoneManager()

	if err := izm.Delete(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
