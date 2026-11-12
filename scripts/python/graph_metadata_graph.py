"""Graph SDK metadata exports."""

import asyncio
import os
import shutil

from colorama import Fore, Style
from graph_metadata_core import GRAPH_EXPORT_ERRORS, DataExtractor
from kiota_abstractions.base_request_configuration import RequestConfiguration
from kiota_abstractions.native_response_handler import NativeResponseHandler
from kiota_http.middleware.options import ResponseHandlerOption


class GraphExtractor(DataExtractor):
    """Extracts data from Microsoft Graph API using the official SDK"""

    async def extract_service_principals_and_endpoints(self, export_path: str) -> None:
        """Extract Service Principals and Endpoints using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Service Principals and Endpoints...{Style.RESET_ALL}")

        # Create tasks for concurrent execution
        tasks = [
            self._extract_endpoints(export_path),
            self._extract_service_principals(export_path),
        ]

        await asyncio.gather(*tasks, return_exceptions=True)
        print(
            f"{Fore.GREEN}✅ Service Principals and Endpoints extraction completed{Style.RESET_ALL}"
        )

    async def _extract_endpoints(self, export_path: str):
        """Extract endpoints using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing endpoints...{Style.RESET_ALL}")

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Get the specific service principal for Microsoft Graph
            sp = await beta_client.service_principals.by_service_principal_id(
                "0000000a-0000-0000-c000-000000000000"
            ).get()
            if sp and sp.endpoints:
                endpoints_data = list(sp.endpoints)
                endpoints_data.sort(key=lambda x: x.capability or "")
                await self.save_json(endpoints_data, os.path.join(export_path, "Endpoints.json"))
            else:
                print(f"{Fore.YELLOW}    ⚠️ No endpoints found{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract endpoints: {str(e)}{Style.RESET_ALL}")

    async def _extract_service_principals(self, export_path: str):
        """Extract service principals using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing service principals...{Style.RESET_ALL}")

        service_principals_path = os.path.join(export_path, "ServicePrincipals")
        if os.path.exists(service_principals_path):
            shutil.rmtree(service_principals_path)
        os.makedirs(service_principals_path)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Get all service principals with automatic pagination
            service_principals = await beta_client.service_principals.get()

            if service_principals and service_principals.value:
                # Save each service principal concurrently
                tasks = []
                for sp in service_principals.value:
                    if sp.app_id:
                        file_path = os.path.join(service_principals_path, f"{sp.app_id}.json")
                        tasks.append(self.save_json(sp, file_path))

                # Process in batches to avoid overwhelming the file system
                batch_size = 20
                for i in range(0, len(tasks), batch_size):
                    batch = tasks[i : i + batch_size]
                    await asyncio.gather(*batch, return_exceptions=True)
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.YELLOW}    ⚠️ Failed to extract service principals:"
                f" {str(e)}{Style.RESET_ALL}"
            )

    async def extract_role_definitions(self, export_path: str) -> None:
        """Extract Role Definitions using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Role Definitions...{Style.RESET_ALL}")

        role_definitions_path = os.path.join(export_path, "RoleDefinitions")
        if os.path.exists(role_definitions_path):
            shutil.rmtree(role_definitions_path)
        os.makedirs(role_definitions_path)

        providers = [
            "cloudPC",
            "deviceManagement",
            "directory",
            "entitlementManagement",
            "exchange",
        ]

        # Process providers concurrently
        tasks = [
            self._extract_role_definitions_for_provider(provider, role_definitions_path)
            for provider in providers
        ]
        results = await asyncio.gather(*tasks, return_exceptions=True)

        successful_providers = []
        failed_providers = []

        for provider, result in zip(providers, results):
            if isinstance(result, Exception):
                error_message = str(result)
                if (
                    "Authorization_RequestDenied" in error_message
                    or "Insufficient privileges" in error_message
                ):
                    print(
                        f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {provider} role"
                        f" definitions - skipping{Style.RESET_ALL}"
                    )
                else:
                    print(
                        f"{Fore.RED}    ❌ Failed to extract {provider} role definitions:"
                        f" {error_message}{Style.RESET_ALL}"
                    )
                failed_providers.append(provider)
            else:
                successful_providers.append(provider)

        if successful_providers:
            print(
                f"{Fore.GREEN}✅ Role Definitions extraction completed for:"
                f" {', '.join(successful_providers)}{Style.RESET_ALL}"
            )

        if failed_providers:
            print(
                f"{Fore.YELLOW}⚠️ Role Definitions extraction failed for:"
                f" {', '.join(failed_providers)} (insufficient permissions){Style.RESET_ALL}"
            )

    async def _extract_role_definitions_for_provider(
        self, provider: str, role_definitions_path: str
    ):
        """Extract role definitions for a specific provider using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing {provider} role definitions...{Style.RESET_ALL}")

        provider_path = os.path.join(role_definitions_path, provider)
        os.makedirs(provider_path)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Get role definitions for the specific provider
            role_definitions = await beta_client.role_management.by_rbac_application_id(
                provider
            ).role_definitions.get()

            if role_definitions and role_definitions.value:
                # Save role definitions concurrently
                tasks = []
                for role_def in role_definitions.value:
                    if role_def.id:
                        file_path = os.path.join(provider_path, f"{role_def.id}.json")
                        tasks.append(self.save_json(role_def, file_path))

                await asyncio.gather(*tasks, return_exceptions=True)
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ No role definitions found for {provider}{Style.RESET_ALL}"
                )
        except GRAPH_EXPORT_ERRORS as e:
            raise RuntimeError(f"Failed to extract {provider} role definitions: {str(e)}") from e

    async def extract_resource_operations(self, export_path: str) -> None:
        """Extract Resource Operations using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Resource Operations...{Style.RESET_ALL}")

        resource_operations_path = os.path.join(export_path, "ResourceOperations")
        if os.path.exists(resource_operations_path):
            shutil.rmtree(resource_operations_path)
        os.makedirs(resource_operations_path)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Get device management resource operations
            resource_operations = await beta_client.device_management.resource_operations.get()

            if resource_operations and resource_operations.value:
                # Save operations concurrently
                tasks = []
                for operation in resource_operations.value:
                    if operation.id:
                        file_path = os.path.join(resource_operations_path, f"{operation.id}.json")
                        tasks.append(self.save_json(operation, file_path))

                await asyncio.gather(*tasks, return_exceptions=True)

            print(f"{Fore.GREEN}✅ Resource Operations extraction completed{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(f"{Fore.RED}❌ Failed to extract Resource Operations: {str(e)}{Style.RESET_ALL}")
            raise

    async def extract_defender_hunting_tables(self, export_path: str) -> None:
        """Extract Defender Hunting Tables using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Defender Hunting Table Schemas...{Style.RESET_ALL}")

        defender_path = os.path.join(export_path, "Defender")
        if os.path.exists(defender_path):
            shutil.rmtree(defender_path)
        os.makedirs(defender_path)

        hunting_tables = [
            "AlertEvidence",
            "AlertInfo",
            "BehaviorEntities",
            "BehaviorInfo",
            "AADSignInEventsBeta",
            "AADSpnSignInEventsBeta",
            "CloudAppEvents",
            "IdentityInfo",
            "IdentityLogonEvents",
            "EmailAttachmentInfo",
            "EmailEvents",
            "EmailPostDeliveryEvents",
            "EmailUrlInfo",
            "UrlClickEvents",
            "ExposureGraphEdges",
            "ExposureGraphNodes",
        ]

        # Process tables concurrently but with smaller batches for hunting queries
        tasks = []
        for table in hunting_tables:
            tasks.append(self._extract_hunting_table_schema(table, defender_path))

        results = await asyncio.gather(*tasks, return_exceptions=True)

        successful_tables = []
        failed_tables = []

        for table, result in zip(hunting_tables, results):
            if isinstance(result, Exception):
                failed_tables.append(table)
            else:
                successful_tables.append(table)

        if successful_tables:
            print(
                f"{Fore.GREEN}✅ Defender Hunting Tables extraction completed for:"
                f" {len(successful_tables)}/{len(hunting_tables)} tables{Style.RESET_ALL}"
            )

        if failed_tables:
            print(
                f"{Fore.YELLOW}⚠️ Some Defender tables could not be accessed: {len(failed_tables)}"
                f" tables{Style.RESET_ALL}"
            )

    async def _extract_hunting_table_schema(self, table: str, defender_path: str):
        """Extract schema for a single hunting table using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing {table} schema...{Style.RESET_ALL}")

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Create hunting query to get schema
            query = (
                f'{table} | getschema | project Description="", Type=split(DataType, ".")[1],'
                ' Entity="", Name=ColumnName'
            )

            request_body = {"query": query}

            response = await beta_client.security.run_hunting_query.post(body=request_body)

            if response and response.results:
                await self.save_json(response.results, os.path.join(defender_path, f"{table}.json"))
                return True
            print(f"{Fore.YELLOW}    ⚠️ No results returned for {table}{Style.RESET_ALL}")
            return False
        except GRAPH_EXPORT_ERRORS as e:
            error_message = str(e)
            if (
                "Authorization_RequestDenied" in error_message
                or "Insufficient privileges" in error_message
            ):
                print(
                    f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {table} -"
                    f" skipping{Style.RESET_ALL}"
                )
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ Failed to get schema for {table}:"
                    f" {error_message}{Style.RESET_ALL}"
                )
            return False

    async def extract_dcv2_configuration_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Settings using Graph SDK with direct URLs"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Settings...{Style.RESET_ALL}")

        dcv2_path = os.path.join(export_path, "DCv2")
        if os.path.exists(dcv2_path):
            shutil.rmtree(dcv2_path)

        print(f"{Fore.YELLOW}  📋 Processing configuration settings...{Style.RESET_ALL}")
        settings_path = os.path.join(dcv2_path, "Settings")
        os.makedirs(settings_path, exist_ok=True)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Use the working pattern with with_url() and proper request config
            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )

            # Get device management configuration settings using direct URL
            config_settings_url = (
                "https://graph.microsoft.com/beta/deviceManagement/configurationSettings"
            )
            config_settings_response = await beta_client.device_management.with_url(
                config_settings_url
            ).get(request_configuration=request_config)
            config_settings_data = config_settings_response.json()

            if config_settings_data and config_settings_data.get("value"):
                # Save settings concurrently
                tasks = []
                for setting in config_settings_data["value"]:
                    if setting.get("id"):
                        # Remove version field if present
                        if "version" in setting:
                            del setting["version"]
                        if "riskLevel" in setting:
                            del setting["riskLevel"]

                        file_path = os.path.join(settings_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))

                # Process in batches
                batch_size = 20
                for i in range(0, len(tasks), batch_size):
                    batch = tasks[i : i + batch_size]
                    await asyncio.gather(*batch, return_exceptions=True)

            # Create backwards compatibility folder
            backwards_compat_path = os.path.join(export_path, "settings")
            if os.path.exists(backwards_compat_path):
                shutil.rmtree(backwards_compat_path)
            shutil.copytree(settings_path, backwards_compat_path)

            print(
                f"{Fore.GREEN}✅ DCv2 Configuration Settings extraction completed{Style.RESET_ALL}"
            )
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}❌ Failed to extract DCv2 Configuration Settings:"
                f" {str(e)}{Style.RESET_ALL}"
            )
            raise

    async def extract_dcv2_compliance_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Compliance Settings using Graph SDK with direct URLs"""
        print(
            f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Compliance"
            f" Settings...{Style.RESET_ALL}"
        )

        compliance_path = os.path.join(export_path, "DCv2", "Compliance")
        os.makedirs(compliance_path, exist_ok=True)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )

            # Get device management compliance settings using direct URL
            compliance_settings_url = (
                "https://graph.microsoft.com/beta/deviceManagement/complianceSettings"
            )
            compliance_settings_response = await beta_client.device_management.with_url(
                compliance_settings_url
            ).get(request_configuration=request_config)
            compliance_settings_data = compliance_settings_response.json()

            if compliance_settings_data and compliance_settings_data.get("value"):
                tasks = []
                for setting in compliance_settings_data["value"]:
                    if setting.get("id"):
                        # Remove version field if present
                        if "version" in setting:
                            del setting["version"]
                        if "riskLevel" in setting:
                            del setting["riskLevel"]

                        file_path = os.path.join(compliance_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))

                await asyncio.gather(*tasks, return_exceptions=True)

            print(f"{Fore.GREEN}✅ DCv2 Compliance Settings extraction completed{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}❌ Failed to extract DCv2 Compliance Settings:"
                f" {str(e)}{Style.RESET_ALL}"
            )
            raise

    async def extract_dcv2_policy_templates(self, export_path: str) -> None:
        """Extract Device Configuration v2 Policy Templates using Graph SDK with direct URLs"""
        print(
            f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Policy Templates...{Style.RESET_ALL}"
        )

        templates_path = os.path.join(export_path, "DCv2", "Templates")
        os.makedirs(templates_path, exist_ok=True)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )

            # Get device management configuration policy templates using direct URL
            templates_url = (
                "https://graph.microsoft.com/beta/deviceManagement/configurationPolicyTemplates"
            )
            templates_response = await beta_client.device_management.with_url(templates_url).get(
                request_configuration=request_config
            )
            templates_data = templates_response.json()

            if templates_data and templates_data.get("value"):
                tasks = []
                for template in templates_data["value"]:
                    if template.get("id"):
                        file_path = os.path.join(templates_path, f"{template['id']}.json")
                        tasks.append(self.save_json(template, file_path))

                await asyncio.gather(*tasks, return_exceptions=True)

            print(f"{Fore.GREEN}✅ DCv2 Policy Templates extraction completed{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}❌ Failed to extract DCv2 Policy Templates: {str(e)}{Style.RESET_ALL}"
            )
            raise

    async def extract_dcv2_inventory_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Inventory Settings using Graph SDK with direct URLs"""
        print(
            f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Inventory"
            f" Settings...{Style.RESET_ALL}"
        )

        inventory_path = os.path.join(export_path, "DCv2", "Inventory")
        os.makedirs(inventory_path, exist_ok=True)

        try:
            beta_client = self.graph_sdk.get_beta_client()

            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )

            # Get device management inventory settings using direct URL
            inventory_url = "https://graph.microsoft.com/beta/deviceManagement/inventorySettings"
            inventory_response = await beta_client.device_management.with_url(inventory_url).get(
                request_configuration=request_config
            )
            inventory_data = inventory_response.json()

            if inventory_data and inventory_data.get("value"):
                tasks = []
                for setting in inventory_data["value"]:
                    if setting.get("id"):
                        # Remove version field if present
                        if "version" in setting:
                            del setting["version"]
                        if "riskLevel" in setting:
                            del setting["riskLevel"]

                        file_path = os.path.join(inventory_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))

                await asyncio.gather(*tasks, return_exceptions=True)

            print(f"{Fore.GREEN}✅ DCv2 Inventory Settings extraction completed{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}❌ Failed to extract DCv2 Inventory Settings: {str(e)}{Style.RESET_ALL}"
            )
            raise
