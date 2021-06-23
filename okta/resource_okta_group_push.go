package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupPush() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupPushCreate,
		ReadContext:   resourceGroupPushRead,
		UpdateContext: resourceGroupPushUpdate,
		DeleteContext: resourceGroupPushDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": statusSchema,
		},
	}
}

func resourceGroupPushRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gp, resp, err := getSupplementFromMetadata(m).GetGroupPushMapping(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
	)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group push mapping: %v", err)
	}

	if gp == nil {
		d.SetId("")
	} else {
		d.SetId(gp.MappingID)
		d.Set("status", gp.Status)
	}
	return nil
}

func resourceGroupPushCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gp, _, err := getSupplementFromMetadata(m).CreateGroupPushMapping(
		ctx,
		d.Get("app_id").(string),
		d.Get("group_id").(string),
		d.Get("status").(string),
	)
	if err != nil {
		return diag.Errorf("failed to create group push: %v", err)
	}

	d.SetId(gp.MappingID)
	return nil
}

func resourceGroupPushUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).UpdateGroupPushMapping(
		ctx,
		d.Get("app_id").(string),
		d.Id(),
		d.Get("status").(string),
	)
	if err != nil {
		return diag.Errorf("failed to update group push: %v", err)
	}
	return nil
}

func resourceGroupPushDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteGroupPushMapping(
		ctx,
		d.Get("app_id").(string),
		d.Id(),
	)
	if err != nil {
		return diag.Errorf("failed to delete group push: %v", err)
	}
	return nil
}
