package stripe

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func resourceStripeFile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeFileRead,
		CreateContext: resourceStripeFileCreate,
		DeleteContext: resourceStripeFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The returned file type (for example, csv, pdf, jpg, or png).",
			},
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The suitable name for saving the file to a filesystem.",
			},
			"base64content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A content file to upload encoded by base64.",
			},
			"purpose": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The purpose of the uploaded file.",
			},
			"link_data": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional parameters that automatically create a file link for the newly created file.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create": {
							Type:     schema.TypeBool,
							Required: true,
							Description: "Set this to true to create a file link for the newly created file. " +
								"Creating a link is only possible when the file’s purpose is one of the following: " +
								"business_icon, business_logo, customer_signature, dispute_evidence, pci_document, " +
								"tax_document_user_upload, or terminal_reader_splashscreen.",
						},
						"expires_at": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The link isn’t available after this future timestamp.",
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Set of key-value pairs that you can attach to an object. " +
								"This can be useful for storing additional information about the object " +
								"in a structured format.",
						},
					},
				},
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique identifier for the object.",
						},
						"object": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "String representing the object’s type. " +
								"Objects of the same type share the same value.",
						},
						"created": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time at which the object was created. Measured in seconds since the Unix epoch.",
						},
						"expired": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Returns if the link is already expired.",
						},
						"expires_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time that the link expires",
						},
						"livemode": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: "Has the value true if the object exists in live mode or the value false " +
								"if the object exists in test mode.",
						},
						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "Set of key-value pairs that you can attach to an object. " +
								"This can be useful for storing additional information about the object in a structured format.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The publicly accessible URL to download the file.",
						},
					},
				},
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String representing the object’s type. Objects of the same type share the same value.",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time at which the object was created. Measured in seconds since the Unix epoch.",
			},
			"expires_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The file expires and isn’t available at this time in epoch seconds.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file object in bytes.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Use your live secret API key to download the file from this URL.",
			},
		},
	}
}

func resourceStripeFileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var file *stripe.File
	var err error

	err = retryWithBackOff(func() error {
		file, err = c.Files.Get(d.Id(), nil)
		return err
	})
	switch {
	case isNotFoundErr(err):
		d.SetId("") // remove when resource does not exist
		return nil
	case err != nil:
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("type", file.Type),
		d.Set("filename", file.Filename),
		d.Set("purpose", file.Purpose),
		func() error {
			var linkMapSlice []map[string]interface{}
			if file.Links != nil && len(file.Links.Data) > 0 {
				for _, linkData := range file.Links.Data {
					linkMap := map[string]interface{}{
						"id":         linkData.ID,
						"object":     linkData.Object,
						"created":    linkData.Created,
						"expired":    linkData.Expired,
						"expires_at": linkData.ExpiresAt,
						"livemode":   linkData.Livemode,
						"metadata":   linkData.Metadata,
						"url":        linkData.URL,
					}
					linkMapSlice = append(linkMapSlice, linkMap)
				}
			}
			return d.Set("links", linkMapSlice)
		}(),
		d.Set("object", file.Object),
		d.Set("created", file.Created),
		d.Set("expires_at", file.ExpiresAt),
		d.Set("size", file.Size),
		d.Set("url", file.URL),
	)
}

func resourceStripeFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var file *stripe.File
	var err error

	params := &stripe.FileParams{
		Filename: stripe.String(ExtractString(d, "filename")),
		Purpose:  stripe.String(ExtractString(d, "purpose")),
	}

	content, err := base64.StdEncoding.DecodeString(ExtractString(d, "base64content"))
	if err != nil {
		return diag.FromErr(err)
	}

	params.FileReader = bytes.NewReader(content)

	if linkData, set := d.GetOk("link_data"); set {
		params.FileLinkData = &stripe.FileFileLinkDataParams{}
		for k, v := range ToMap(linkData) {
			switch k {
			case "create":
				params.FileLinkData.Create = stripe.Bool(ToBool(v))
			case "expires_at":
				expires := ToInt64(v)
				if expires > 0 {
					params.FileLinkData.ExpiresAt = stripe.Int64(expires)
				}
			case "metadata":
				for key, value := range ToMap(v) {
					params.FileLinkData.AddMetadata(key, ToString(value))
				}
			}
		}
	}

	err = retryWithBackOff(func() error {
		file, err = c.Files.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(file.ID)
	return resourceStripeFileRead(ctx, d, m)
}

func resourceStripeFileDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "[WARN] Stripe API doesn't support deletion of file")
	d.SetId("")
	return nil
}
