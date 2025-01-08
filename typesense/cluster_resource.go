package typesense

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &clusterResource{}
	_ resource.ResourceWithConfigure = &clusterResource{}

	clusterResourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Autogenerated ID assigned by the Typesense engine.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "A string to identify the cluster for your reference in the Typesense Cloud Web console.",
				Computed:    true,
				Optional:    true,
			},
			"memory": schema.StringAttribute{
				Description: "How much RAM this cluster should have. Available options here: https://typesense.org/docs/cloud-management-api/v1/cluster-management.html#memory",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vcpu": schema.StringAttribute{
				Description: "How many CPU cores this cluster should have. Available options here: https://typesense.org/docs/cloud-management-api/v1/cluster-management.html#vcpu",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"high_performance_disk": schema.StringAttribute{
				Description: "When set to yes, the provisioned hard disk will be co-located on the same physical server that runs the node.",
				Computed:    true,
				Default:     stringdefault.StaticString("no"),
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"typesense_server_version": schema.StringAttribute{
				Description: "Cluster Typesense server version at creation time.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"high_availability": schema.StringAttribute{
				Description: "When set to yes, at least 3 nodes are provisioned in 3 different data centers to form a highly available (HA) cluster and your data is automatically replicated between all nodes.",
				Computed:    true,
				Default:     stringdefault.StaticString("no"),
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_delivery_network": schema.StringAttribute{
				Description: "When not off, nodes are provisioned in different regions and the node that's closest to it's originating location serves the traffic.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancing": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Region where the nodes should be geographically placed. Available options here: https://typesense.org/docs/cloud-management-api/v1/cluster-management.html#regions",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_upgrade_capacity": schema.BoolAttribute{
				Description: "When set to true, your cluster will be automatically upgraded when best-practice RAM/CPU thresholds are exceeded in a 12-hour rolling window.",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Current status of your cluster.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostnames": schema.SingleNestedAttribute{
				Description: "Hostnames for the cluster.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"load_balanced": schema.StringAttribute{
						Description: "Indicates if load balancing is enabled.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(), // Handle unknown values
						},
					},
					"nodes": schema.ListAttribute{
						Description: "List of nodes in the cluster.",
						ElementType: types.StringType,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(), // Handle unknown lists
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(), // Handle unknown for the entire nested object
				},
			},
		},
	}
)

// NewClusterResource is a helper function to simplify the provider implementation.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// clusterResource is the resource implementation.
type clusterResource struct {
	client *typesenseClient
}

// Configure adds the provider configured client to the resource.
func (cr *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cr.client = req.ProviderData.(*typesenseClient)
}

// Metadata returns the resource type name.
func (cr *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the resource.
func (cr *clusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = clusterResourceSchema
}

// Create creates the resource and sets the initial Terraform state.
func (cr *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan typesenseClusterModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new cluster
	cluster, err := cr.client.CreateCluster(typesenseCluster{
		Memory:                plan.Memory.ValueString(),
		VCPU:                  plan.VCPU.ValueString(),
		Regions:               []string{plan.Region.ValueString()},
		HighAvailability:      plan.HighAvailability.ValueString(),
		SearchDeliveryNetwork: plan.SearchDeliveryNetwork.ValueString(),
		HighPerformanceDisk:   plan.HighPerformanceDisk.ValueString(),
		Name:                  plan.Name.ValueString(),
		AutoUpgradeCapacity:   plan.AutoUpgradeCapacity.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	// Waiting until state is not provisioning.
	clusterId := cluster.ID
	for {
		time.Sleep(8 * time.Second)
		cluster, err = cr.client.GetCluster(clusterId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for cluster state",
				"Cluster created, but could not reach expected state: "+err.Error(),
			)
			return
		}
		if cluster.Status != "in_service" {
			continue
		}
		break
	}

	plan.ID = types.StringValue(cluster.ID)
	plan.Name = types.StringValue(cluster.Name)
	plan.Memory = types.StringValue(cluster.Memory)
	plan.VCPU = types.StringValue(cluster.VCPU)
	plan.HighPerformanceDisk = types.StringValue(cluster.HighPerformanceDisk)
	plan.TypesenseServerVersion = types.StringValue(cluster.TypesenseServerVersion)
	plan.HighAvailability = types.StringValue(cluster.HighAvailability)
	plan.SearchDeliveryNetwork = types.StringValue(cluster.SearchDeliveryNetwork)
	plan.LoadBalancing = types.StringValue(cluster.LoadBalancing)
	plan.Region = types.StringValue(cluster.Regions[0])
	plan.AutoUpgradeCapacity = types.BoolValue(cluster.AutoUpgradeCapacity)
	plan.Status = types.StringValue(cluster.Status)

	nodes := make([]types.String, len(cluster.Hostnames.Nodes))
	for i, elem := range cluster.Hostnames.Nodes {
		nodes[i] = types.StringValue(elem)
	}

	plan.Hostnames = typesenseHostnamesModel{
		LoadBalanced: types.StringValue(cluster.Hostnames.LoadBalanced),
		Nodes:        nodes,
	}
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (cr *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state typesenseClusterModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed cluster value from Typesense
	cluster, err := cr.client.GetCluster(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Typesense Cluster",
			"Could not read Typesense Cluster ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	state.ID = types.StringValue(cluster.ID)
	state.Name = types.StringValue(cluster.Name)
	state.Memory = types.StringValue(cluster.Memory)
	state.VCPU = types.StringValue(cluster.VCPU)
	state.HighPerformanceDisk = types.StringValue(cluster.HighPerformanceDisk)
	state.TypesenseServerVersion = types.StringValue(cluster.TypesenseServerVersion)
	state.HighAvailability = types.StringValue(cluster.HighAvailability)
	state.SearchDeliveryNetwork = types.StringValue(cluster.SearchDeliveryNetwork)
	state.LoadBalancing = types.StringValue(cluster.LoadBalancing)
	state.Region = types.StringValue(cluster.Regions[0])
	state.AutoUpgradeCapacity = types.BoolValue(cluster.AutoUpgradeCapacity)
	state.Status = types.StringValue(cluster.Status)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (cr *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan typesenseClusterModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing cluster
	err := cr.client.UpdateCluster(typesenseCluster{
		ID:                  plan.ID.ValueString(),
		Name:                plan.Name.ValueString(),
		AutoUpgradeCapacity: plan.AutoUpgradeCapacity.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Typesense Cluster",
			"Could not update cluster, unexpected error: "+err.Error(),
		)
		return
	}
	// Get refreshed cluster value from Typesense
	cluster, err := cr.client.GetCluster(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Typesense Cluster",
			"Could not read Typesense Cluster ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(cluster.ID)
	plan.Name = types.StringValue(cluster.Name)
	plan.Memory = types.StringValue(cluster.Memory)
	plan.VCPU = types.StringValue(cluster.VCPU)
	plan.HighPerformanceDisk = types.StringValue(cluster.HighPerformanceDisk)
	plan.TypesenseServerVersion = types.StringValue(cluster.TypesenseServerVersion)
	plan.HighAvailability = types.StringValue(cluster.HighAvailability)
	plan.SearchDeliveryNetwork = types.StringValue(cluster.SearchDeliveryNetwork)
	plan.LoadBalancing = types.StringValue(cluster.LoadBalancing)
	plan.Region = types.StringValue(cluster.Regions[0])
	plan.AutoUpgradeCapacity = types.BoolValue(cluster.AutoUpgradeCapacity)
	plan.Status = types.StringValue(cluster.Status)

	nodes := make([]types.String, len(cluster.Hostnames.Nodes))
	for i, elem := range cluster.Hostnames.Nodes {
		nodes[i] = types.StringValue(elem)
	}

	plan.Hostnames = typesenseHostnamesModel{
		LoadBalanced: types.StringValue(cluster.Hostnames.LoadBalanced),
		Nodes:        nodes,
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (cr *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state typesenseClusterModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Terminate cluster
	err := cr.client.TerminateCluster(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Typesense Cluster",
			"Could not delete cluster, unexpected error: "+err.Error(),
		)
		return
	}
}

func (cr *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
