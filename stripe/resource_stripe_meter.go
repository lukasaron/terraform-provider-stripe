package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func resourceStripeMeter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeMeterRead,
		CreateContext: resourceStripeMeterCreate,
		UpdateContext: resourceStripeMeterUpdate,
		DeleteContext: resourceStripeMeterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"default_aggregation": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The default settings to aggregate a meter’s events with",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"formula": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Specifies how events are aggregated. Allowed values " +
								"are count to count the number of events and sum to sum each " +
								"event’s value.",
						},
					},
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The meter’s name.",
			},
			"event_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The name of the meter event to record usage for. " +
					"Corresponds with the event_name field on meter events",
			},
			"customer_mapping": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Fields that specify how to map a meter event to a customer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_payload_key": {
							Type:     schema.TypeString,
							Required: true,
							Description: "The key in the usage event payload to use for mapping " +
								"the event to a customer.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							Description: "The method for mapping a meter event to a customer. " +
								"Must be by_id",
						},
					},
				},
			},
			"event_time_window": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time window to pre-aggregate meter events for, if any.",
			},
			"value_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Fields that specify how to calculate a meter event’s value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_payload_key": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The key in the usage event payload to use as the ` +
								`value for this meter. For example, if the event payload  ` +
								`contains usage on a bytes_used field, then set the ` +
								`event_payload_key to “bytes_used”`,
						},
					},
				},
			},
		},
	}
}

func resourceStripeMeterRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var meter *stripe.BillingMeter
	var err error

	err = retryWithBackOff(func() error {
		params := &stripe.BillingMeterParams{}

		meter, err = c.BillingMeters.Get(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("display_name", meter.DisplayName),
		d.Set("event_name", meter.EventName),
		d.Set("event_time_window", meter.EventTimeWindow),
		func() error {
			if meter.DefaultAggregation != nil {
				return d.Set("default_aggregation", []map[string]interface{}{
					{
						"formula": meter.DefaultAggregation.Formula,
					},
				})
			}
			return nil
		}(),
		func() error {
			if meter.CustomerMapping != nil {
				return d.Set("customer_mapping", []map[string]interface{}{
					{
						"event_payload_key": meter.CustomerMapping.EventPayloadKey,
						"type":              meter.CustomerMapping.Type,
					},
				})
			}
			return nil
		}(),
		func() error {
			if meter.ValueSettings != nil {
				return d.Set("value_settings", []map[string]interface{}{
					{
						"event_payload_key": meter.ValueSettings.EventPayloadKey,
					},
				})
			}
			return nil
		}(),
	)
}

func resourceStripeMeterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var meter *stripe.BillingMeter
	var err error

	params := &stripe.BillingMeterParams{
		DisplayName:     stripe.String(ExtractString(d, "display_name")),
		EventName:       stripe.String(ExtractString(d, "event_name")),
		EventTimeWindow: stripe.String(ExtractString(d, "event_time_window")),
		ValueSettings:   &stripe.BillingMeterValueSettingsParams{},
	}

	if defaultAggregation, set := d.GetOk("default_aggregation"); set {
		params.DefaultAggregation = &stripe.BillingMeterDefaultAggregationParams{}
		defaultAggregationMap := ToMap(defaultAggregation)
		for k, v := range defaultAggregationMap {
			switch {
			case k == "formula" && ToString(v) != "":
				params.DefaultAggregation.Formula = stripe.String(ToString(v))
			}
		}
	}

	if customerMapping, set := d.GetOk("customer_mapping"); set {
		params.CustomerMapping = &stripe.BillingMeterCustomerMappingParams{}
		customerMappingMap := ToMap(customerMapping)
		for k, v := range customerMappingMap {
			switch {
			case k == "event_payload_key" && ToString(v) != "":
				params.CustomerMapping.EventPayloadKey = stripe.String(ToString(v))
			case k == "type" && ToString(v) != "":
				params.CustomerMapping.Type = stripe.String(ToString(v))
			}
		}
	}

	if valueSettings, set := d.GetOk("value_settings"); set {
		params.ValueSettings = &stripe.BillingMeterValueSettingsParams{}
		valueSettingsMap := ToMap(valueSettings)
		for k, v := range valueSettingsMap {
			switch {
			case k == "event_payload_key" && ToString(v) != "":
				params.ValueSettings.EventPayloadKey = stripe.String(ToString(v))
			}
		}
	}

	err = retryWithBackOff(func() error {
		meter, err = c.BillingMeters.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(meter.ID)
	return resourceStripeMeterRead(ctx, d, m)
}

func resourceStripeMeterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.BillingMeterParams{}

	if d.HasChange("display_name") {
		params.DisplayName = stripe.String(ExtractString(d, "display_name"))
	}

	err = retryWithBackOff(func() error {
		_, err = c.BillingMeters.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeMeterRead(ctx, d, m)
}

func resourceStripeMeterDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := stripe.BillingMeterDeactivateParams{}

	err = retryWithBackOff(func() error {
		_, err = c.BillingMeters.Deactivate(d.Id(), &params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
