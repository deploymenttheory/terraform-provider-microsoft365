package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccMacOSPlatformScript_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test macOS Script ` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
						description      = "Test description for macOS platform script"
						script_content   = "#!/bin/bash\necho 'Hello World'"
						run_as_account   = "system"
						file_name        = "test-script.sh"
						execution_frequency = "P1D"
						retry_count      = 3
						block_execution_notifications = true
						assignments = {
							all_devices   = false
							all_users     = true
							include_group_ids = []
							exclude_group_ids = []
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script "+mocks.TestName()+strconv.Itoa(rand.IntN(9999))),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1D"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "false"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "id", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScript_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateMicrosoftGraphMocks()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create/post_device_shell_scripts.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/post_assign.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001?%24expand=assignments",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_device_shell_script_with_assignments.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test macOS Script"
						description      = "Test description for macOS platform script"
						script_content   = "#!/bin/bash\necho 'Hello World'"
						run_as_account   = "system"
						file_name        = "test-script.sh"
						execution_frequency = "P1D"
						retry_count      = 3
						block_execution_notifications = true
						assignments = {
							all_devices   = false
							all_users     = true
							include_group_ids = []
							exclude_group_ids = []
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Test description for macOS platform script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1D"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScript_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateMicrosoftGraphMocks()

	patch_device_shell_scripts_inx := 0
	post_assign_inx := 0
	get_device_shell_scripts_inx := 0

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Update/post_device_shell_scripts.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			patch_device_shell_scripts_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/patch_device_shell_scripts_%d.json", patch_device_shell_scripts_inx)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		func(req *http.Request) (*http.Response, error) {
			post_assign_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/post_assign_%d.json", post_assign_inx)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001?%24expand=assignments",
		func(req *http.Request) (*http.Response, error) {
			get_device_shell_scripts_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/get_device_shell_script_with_assignments_%d.json", get_device_shell_scripts_inx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test macOS Script"
						description      = "Test description for macOS platform script"
						script_content   = "#!/bin/bash\necho 'Hello World'"
						run_as_account   = "system"
						file_name        = "test-script.sh"
						execution_frequency = "P1D"
						retry_count      = 3
						block_execution_notifications = true
						assignments = {
							all_devices   = false
							all_users     = true
							include_group_ids = []
							exclude_group_ids = []
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Test description for macOS platform script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1D"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "false"),
				),
			},
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Updated macOS Script"
						description      = "Updated description for macOS platform script"
						script_content   = "#!/bin/bash\necho 'Hello Updated World'\necho 'Second line'"
						run_as_account   = "user"
						file_name        = "updated-script.sh"
						execution_frequency = "P7D"
						retry_count      = 5
						block_execution_notifications = false
						assignments = {
							all_devices   = true
							all_users     = false
							include_group_ids = ["11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"]
							exclude_group_ids = ["33333333-3333-3333-3333-333333333333"]
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Updated macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Updated description for macOS platform script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "updated-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P7D"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScript_Validate_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateMicrosoftGraphMocks()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Delete/post_device_shell_scripts.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Delete/post_assign.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001?%24expand=assignments",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Delete/get_device_shell_script_with_assignments.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test macOS Script for Delete"
						description      = "Test description for delete scenario"
						script_content   = "#!/bin/bash\necho 'To be deleted'"
						run_as_account   = "system"
						file_name        = "delete-test-script.sh"
						assignments = {
							all_devices   = false
							all_users     = true
							include_group_ids = []
							exclude_group_ids = []
						}
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script for Delete"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScript_Validate_Assignment_Validation_Errors(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateMicrosoftGraphMocks()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test Assignment Validation"
						script_content   = "#!/bin/bash\necho 'test'"
						run_as_account   = "system"
						file_name        = "test.sh"
						assignments = {
							all_devices   = true
							all_users     = false
							include_group_ids = ["11111111-1111-1111-1111-111111111111"]
							exclude_group_ids = []
						}
					}
				`,
				ExpectError: regexp.MustCompile("cannot assign to All Devices and Include Groups at the same time"),
			},
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test Assignment Validation"
						script_content   = "#!/bin/bash\necho 'test'"
						run_as_account   = "system"
						file_name        = "test.sh"
						assignments = {
							all_devices   = false
							all_users     = true
							include_group_ids = ["11111111-1111-1111-1111-111111111111"]
							exclude_group_ids = []
						}
					}
				`,
				ExpectError: regexp.MustCompile("cannot assign to All Users and Include Groups at the same time"),
			},
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test Assignment Validation"
						script_content   = "#!/bin/bash\necho 'test'"
						run_as_account   = "system"
						file_name        = "test.sh"
						assignments = {
							all_devices   = false
							all_users     = false
							include_group_ids = ["22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111"]
							exclude_group_ids = []
						}
					}
				`,
				ExpectError: regexp.MustCompile("include_group_ids must be in alphanumeric order"),
			},
			{
				Config: mocks.ProviderConfigMinimal() + `
					resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
						display_name     = "Test Assignment Validation"
						script_content   = "#!/bin/bash\necho 'test'"
						run_as_account   = "system"
						file_name        = "test.sh"
						assignments = {
							all_devices   = false
							all_users     = false
							include_group_ids = ["11111111-1111-1111-1111-111111111111"]
							exclude_group_ids = ["11111111-1111-1111-1111-111111111111"]
						}
					}
				`,
				ExpectError: regexp.MustCompile("group .* is used in both include and exclude assignments"),
			},
		},
	})
}
