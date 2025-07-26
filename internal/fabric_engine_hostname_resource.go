package provider

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
)

// FabricEngineHostnameResource implements resource.Resource.
type FabricEngineHostnameResource struct {
	client *ExtrmFabricEngineClient
}

// NewFabricEngineHostnameResource returns a new instance of the resource.
func NewFabricEngineHostnameResource() resource.Resource {
	return &FabricEngineHostnameResource{}
}

// FabricEngineHostnameModel describes the resource model used in Terraform state.
type FabricEngineHostnameModel struct {
	ID       types.String `tfsdk:"id"`
	Hostname types.String `tfsdk:"hostname"`
}

func (r *FabricEngineHostnameResource) Metadata(
	ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_hostname"
}

func (r *FabricEngineHostnameResource) Schema(
	ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true},
			"hostname": schema.StringAttribute{Required: true},
		},
	}
}

// Configure retrieves the provider data (SSH client parameters) and assigns it to the resource.
func (r *FabricEngineHostnameResource) Configure(
	ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*ExtrmFabricEngineClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected client type", "The provider did not return a valid client")
		return
	}
	r.client = c
}

// Create connects to the device and sets the hostname.
func (r *FabricEngineHostnameResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan FabricEngineHostnameModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// SSH client configuration
	cfg := &ssh.ClientConfig{
		User:            r.client.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(r.client.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", r.client.Host, r.client.Port)
	client, err := ssh.Dial("tcp", address, cfg)
	if err != nil {
		resp.Diagnostics.AddError("SSH connection error", err.Error())
		return
	}
	defer client.Close()

	// Start an interactive shell
	session, err := client.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Cannot create SSH session", err.Error())
		return
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		resp.Diagnostics.AddError("Cannot get stdin pipe", err.Error())
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		resp.Diagnostics.AddError("Cannot get stdout pipe", err.Error())
		return
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		resp.Diagnostics.AddError("Cannot get stderr pipe", err.Error())
		return
	}

	if err := session.Shell(); err != nil {
		resp.Diagnostics.AddError("Failed to start remote shell", err.Error())
		return
	}

	// Helper to send commands and collect output
	send := func(cmd string) error {
		_, err := fmt.Fprintf(stdin, "%s\n", cmd)
		return err
	}
	var output string
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			line := scanner.Text()
			output += line + "\n"
		}
	}()

	// Sequence of commands
	// Sequence of commands with error handling
	if err := send("enable"); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to send 'enable': %s\noutput:\n%s", err, output),
		)
		return
	}
	if err := send("configure terminal"); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to send 'configure terminal': %s\noutput:\n%s", err, output),
		)
		return
	}
	if err := send(fmt.Sprintf("sys name %s", plan.Hostname.ValueString())); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to set hostname: %s\noutput:\n%s", err, output),
		)
		return
	}
	// Quit configuration mode
	if err := send("exit"); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to exit configuration mode: %s\noutput:\n%s", err, output),
		)
		return
	}
	// Save the configuration:contentReference[oaicite:0]{index=0}
	if err := send("save config"); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to save configuration: %s\noutput:\n%s", err, output),
		)
		return
	}
	// Exit the CLI session
	if err := send("exit"); err != nil {
		resp.Diagnostics.AddError(
			"SSH command failed",
			fmt.Sprintf("failed to exit session: %s\noutput:\n%s", err, output),
		)
		return
	}
	if err := session.Wait(); err != nil {
		if !strings.Contains(err.Error(), "exited without exit status") {
			resp.Diagnostics.AddError(
				"SSH command sequence failed",
				fmt.Sprintf("error: %s\noutput:\n%s", err, output),
			)
			return
		}
	}

	// Enregistrer l’état
	plan.ID = plan.Hostname
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read fetches the current hostname by executing "show sys-info" and parsing the SysName field.
func (r *FabricEngineHostnameResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state FabricEngineHostnameModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Connect via SSH.
	config := &ssh.ClientConfig{
		User:            r.client.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(r.client.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", r.client.Host, r.client.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		resp.Diagnostics.AddError("SSH connection error", err.Error())
		return
	}
	defer client.Close()

	// Run "show sys-info" and capture the output.
	session, err := client.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Cannot create SSH session", err.Error())
		return
	}
	defer session.Close()

	output, err := session.Output("show sys-info")
	if err != nil {
		resp.Diagnostics.AddError("Failed to run show sys-info", err.Error())
		return
	}

	// Parse the SysName from the command output (the Ansible collection uses a similar pattern:contentReference[oaicite:3]{index=3}).
	re := regexp.MustCompile(`SysName\s+:\s+(\S+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) == 2 {
		state.Hostname = types.StringValue(matches[1])
		state.ID = state.Hostname
	} else {
		// If parsing fails, keep the previous state.
		resp.Diagnostics.AddWarning("Hostname not found", "Could not parse SysName from show sys-info output")
	}

	// Save the updated state.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update changes the hostname only if the desired value differs from the current state.
func (r *FabricEngineHostnameResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan FabricEngineHostnameModel
	var state FabricEngineHostnameModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the hostname is unchanged, no action is needed.
	if plan.Hostname.ValueString() == state.Hostname.ValueString() {
		return
	}

	// Connect and set the new hostname.
	config := &ssh.ClientConfig{
		User:            r.client.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(r.client.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", r.client.Host, r.client.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		resp.Diagnostics.AddError("SSH connection error", err.Error())
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Cannot create SSH session", err.Error())
		return
	}
	defer session.Close()

	cmd := fmt.Sprintf("sys name %s", plan.Hostname.ValueString())
	if err := session.Run(cmd); err != nil {
		resp.Diagnostics.AddError("Command execution failed", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete resets the hostname to a default model (TEST-FABRIC-ENGINE) and removes the resource from state.
func (r *FabricEngineHostnameResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Connect via SSH.
	config := &ssh.ClientConfig{
		User:            r.client.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(r.client.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", r.client.Host, r.client.Port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		resp.Diagnostics.AddError("SSH connection error", err.Error())
		return
	}
	defer client.Close()

	// Set the hostname back to a generic value. You can later replace "TEST-FABRIC-ENGINE"
	// with the actual model of your device.
	session, err := client.NewSession()
	if err != nil {
		resp.Diagnostics.AddError("Cannot create SSH session", err.Error())
		return
	}
	defer session.Close()

	if err := session.Run("sys name TEST-FABRIC-ENGINE"); err != nil {
		resp.Diagnostics.AddError("Failed to reset hostname", err.Error())
		return
	}

	// Remove the resource from Terraform state.
	resp.State.RemoveResource(ctx)
}
