// resource_domain_template.go
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDomainTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainTemplateCreate,
		ReadContext:   resourceDomainTemplateRead,
		UpdateContext: resourceDomainTemplateUpdate,
		DeleteContext: resourceDomainTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(slugRegex, "must only contain alphanumeric characters, hyphens, and underscores"),
			},
			"domain": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"allow_copy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"editable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"langs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lang": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Required: true,
						},
						"body": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

var slugRegex = `^[-a-zA-Z0-9_]+$`

func resourceDomainTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainTemplate := &DomainTemplate{
		Name:      d.Get("name").(string),
		Slug:      d.Get("slug").(string),
		Domain:    d.Get("domain").(int),
		AllowCopy: d.Get("allow_copy").(bool),
		Editable:  d.Get("editable").(bool),
	}

	// Create the domain template
	createdTemplate, err := client.CreateDomainTemplate(domainTemplate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTemplate.UUID)

	// Add langs if defined
	if langs := d.Get("langs").([]interface{}); len(langs) > 0 {
		// We need to create a new template with the langs
		templateWithLangs := &DomainTemplate{
			Name:      createdTemplate.Name,
			Slug:      createdTemplate.Slug,
			Domain:    createdTemplate.Domain,
			AllowCopy: createdTemplate.AllowCopy,
			Editable:  createdTemplate.Editable,
			Langs:     make([]DomainTemplateLang, 0, len(langs)),
		}

		for _, lang := range langs {
			langMap := lang.(map[string]interface{})
			templateWithLangs.Langs = append(templateWithLangs.Langs, DomainTemplateLang{
				Lang:     langMap["lang"].(string),
				Subject:  langMap["subject"].(string),
				Body:     langMap["body"].(string),
				Default:  langMap["default"].(bool),
				Template: createdTemplate.ID,
			})
		}

		_, err = client.UpdateDomainTemplate(createdTemplate.UUID, templateWithLangs)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDomainTemplateRead(ctx, d, m)
}

func resourceDomainTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	domainTemplate, err := client.GetDomainTemplate(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("uuid", domainTemplate.UUID)
	d.Set("name", domainTemplate.Name)
	d.Set("slug", domainTemplate.Slug)
	d.Set("domain", domainTemplate.Domain)
	d.Set("allow_copy", domainTemplate.AllowCopy)
	d.Set("editable", domainTemplate.Editable)
	d.Set("created_at", domainTemplate.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", domainTemplate.UpdatedAt.Format(time.RFC3339))

	if domainTemplate.Langs != nil && len(domainTemplate.Langs) > 0 {
		langs := make([]interface{}, len(domainTemplate.Langs))
		for i, lang := range domainTemplate.Langs {
			langs[i] = map[string]interface{}{
				"lang":    lang.Lang,
				"subject": lang.Subject,
				"body":    lang.Body,
				"default": lang.Default,
			}
		}
		d.Set("langs", langs)
	}

	return diags
}

func resourceDomainTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainTemplate := &DomainTemplate{
		Name:      d.Get("name").(string),
		Slug:      d.Get("slug").(string),
		AllowCopy: d.Get("allow_copy").(bool),
		Editable:  d.Get("editable").(bool),
		Domain:    d.Get("domain").(int),
	}

	// Add langs if defined
	if d.HasChange("langs") || d.HasChange("name") || d.HasChange("slug") || d.HasChange("allow_copy") || d.HasChange("editable") {
		langs := d.Get("langs").([]interface{})
		domainTemplate.Langs = make([]DomainTemplateLang, 0, len(langs))

		for _, lang := range langs {
			langMap := lang.(map[string]interface{})
			domainTemplate.Langs = append(domainTemplate.Langs, DomainTemplateLang{
				Lang:    langMap["lang"].(string),
				Subject: langMap["subject"].(string),
				Body:    langMap["body"].(string),
				Default: langMap["default"].(bool),
			})
		}

		_, err := client.UpdateDomainTemplate(d.Id(), domainTemplate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDomainTemplateRead(ctx, d, m)
}

func resourceDomainTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	err := client.DeleteDomainTemplate(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func dataSourceDomainTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainTemplateRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allow_copy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"langs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics
	uuid := d.Get("uuid").(string)

	domainTemplate, err := client.GetDomainTemplate(uuid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainTemplate.UUID)
	d.Set("uuid", domainTemplate.UUID)
	d.Set("name", domainTemplate.Name)
	d.Set("slug", domainTemplate.Slug)
	d.Set("domain", domainTemplate.Domain)
	d.Set("allow_copy", domainTemplate.AllowCopy)
	d.Set("editable", domainTemplate.Editable)
	d.Set("created_at", domainTemplate.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", domainTemplate.UpdatedAt.Format(time.RFC3339))

	if domainTemplate.Langs != nil && len(domainTemplate.Langs) > 0 {
		langs := make([]interface{}, len(domainTemplate.Langs))
		for i, lang := range domainTemplate.Langs {
			langs[i] = map[string]interface{}{
				"lang":    lang.Lang,
				"subject": lang.Subject,
				"body":    lang.Body,
				"default": lang.Default,
			}
		}
		d.Set("langs", langs)
	}

	return diags
}
