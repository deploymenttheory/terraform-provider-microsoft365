#!/usr/bin/env python3
"""
Microsoft Graph Data Exporter with M365 Metadata Extraction

This script extracts metadata from Microsoft Graph API and various M365 services,
saving the data to JSON files in a specified directory.

Usage:
    python graph_metadata_extraction.py --tenant-id <tenant_id> --client-id <client_id> --client-secret <client_secret> [--export-path <export_path>]

Example:
    python graph_metadata_extraction.py --tenant-id "your-tenant-id" --client-id "your-client-id" --client-secret "your-client-secret" --export-path "./GraphMetadata"
"""

import argparse
import json
import os
import sys
import time
import re
import shutil
from datetime import datetime, UTC
from typing import Dict, List, Optional, Any, Union
from pathlib import Path

import requests
from colorama import Fore, Style, init
from azure.identity import ClientSecretCredential, CredentialUnavailableError
import urllib.parse

# Initialize colorama
init()

class GraphMetadataExtractor:
    """Class to handle Microsoft Graph metadata extraction."""

    def __init__(self, tenant_id: str, client_id: str, client_secret: str, required_permissions: Optional[List[str]] = None):
        """Initialize the Graph Metadata Extractor.

        Args:
            tenant_id: Azure AD tenant ID
            client_id: Azure AD application client ID
            client_secret: Azure AD application client secret
            required_permissions: List of required Microsoft Graph permissions to validate
        """
        self.tenant_id = tenant_id
        self.client_id = client_id
        self.client_secret = client_secret
        self.required_permissions = required_permissions or [
            # Directory (Microsoft Entra ID) provider
            "RoleManagement.Read.Directory",  # Least privileged for directory role definitions
            "Directory.Read.All",  # Required for directory operations
            
            # Cloud PC provider
            "RoleManagement.Read.CloudPC",  # Least privileged for Cloud PC role definitions
            "CloudPC.Read.All",  # Required for Cloud PC operations
            
            # Device Management (Intune) provider
            "DeviceManagementRBAC.Read.All",  # Required for device management role definitions
            "DeviceManagementServiceConfig.Read.All",
            "DeviceManagementConfiguration.Read.All",
            "Policy.Read.DeviceConfiguration", # Required for settings catalog

            "ThreatHunting.Read.All", # 

            # Exchange Online provider
            "RoleManagement.Read.Exchange",  # Least privileged for Exchange role definitions
            
            # Other required permissions
            "IdentityRiskyServicePrincipal.Read.All",
            
            "SecurityEvents.Read.All"
        ]
        self.credential = None
        self.headers = None
        self.graph_url = "https://graph.microsoft.com/v1.0"
        self.graph_beta_url = "https://graph.microsoft.com/beta"

    def authenticate(self) -> bool:
        """Authenticate with Microsoft Graph using Azure Identity SDK.
        
        Uses ClientSecretCredential from azure-identity for robust authentication
        with automatic token refresh and proper error handling. This is the recommended
        approach for service principal authentication with Microsoft Graph.
        
        Returns:
            bool: True if authentication successful, False otherwise
            
        Raises:
            CredentialUnavailableError: If authentication with Azure AD fails
            Exception: For other authentication-related errors
        """
        print(f"{Fore.CYAN}🔐 Authenticating with Microsoft Graph...{Style.RESET_ALL}")
        try:
            # Create credential using Azure Identity SDK
            self.credential = ClientSecretCredential(
                tenant_id=self.tenant_id,
                client_id=self.client_id,
                client_secret=self.client_secret
            )
            
            # Get token for Microsoft Graph scope
            token = self.credential.get_token("https://graph.microsoft.com/.default")
            
            # Set up headers for requests
            self.headers = {
                'Authorization': f'Bearer {token.token}',
                'Content-Type': 'application/json'
            }
            
            print(f"{Fore.GREEN}✅ Authentication completed successfully{Style.RESET_ALL}")
            return True
            
        except CredentialUnavailableError as e:
            print(f"{Fore.RED}❌ Authentication failed: {e}{Style.RESET_ALL}")
            return False
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to connect to Microsoft Graph: {e}{Style.RESET_ALL}")
            return False

    def _refresh_token_if_needed(self):
        """Refresh access token if needed using Azure Identity SDK.
        
        The Azure Identity SDK handles token refresh automatically,
        but this method can be called to get a fresh token explicitly.
        """
        try:
            if self.credential:
                token = self.credential.get_token("https://graph.microsoft.com/.default")
                self.headers['Authorization'] = f'Bearer {token.token}'
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Token refresh warning: {e}{Style.RESET_ALL}")

    def make_request(self, url: str, method: str = "GET", body: Optional[Dict] = None, handle_paging: bool = True) -> Union[Dict, List]:
        """Make a request to the Microsoft Graph API.

        Args:
            url: The API endpoint to call
            method: HTTP method (GET, POST, etc.)
            body: Request body for POST requests
            handle_paging: Whether to handle pagination

        Returns:
            The JSON response from the API
        """
        retry_count = 0
        max_retries = 3
        all_results = []
        current_url = url

        while True:
            while retry_count <= max_retries:
                try:
                    if method == "POST" and body:
                        response = requests.post(current_url, headers=self.headers, json=body)
                    else:
                        response = requests.get(current_url, headers=self.headers)
                    
                    # If token expired, refresh and retry once
                    if response.status_code == 401:
                        self._refresh_token_if_needed()
                        if method == "POST" and body:
                            response = requests.post(current_url, headers=self.headers, json=body)
                        else:
                            response = requests.get(current_url, headers=self.headers)
                    
                    response.raise_for_status()
                    data = response.json()

                    if "value" in data:
                        all_results.extend(data["value"])
                        if handle_paging and "@odata.nextLink" in data:
                            current_url = data["@odata.nextLink"]
                            retry_count = 0
                            break
                        else:
                            return all_results
                    else:
                        return data

                except requests.exceptions.RequestException as e:
                    status_code = e.response.status_code if hasattr(e, 'response') and e.response else None
                    if status_code in (429, 503):
                        retry_count += 1
                        retry_after = int(e.response.headers.get('Retry-After', retry_count * 10))
                        print(f"{Fore.YELLOW}⚠️ Rate limited. Retrying in {retry_after} seconds... (Attempt {retry_count}/{max_retries}){Style.RESET_ALL}")
                        time.sleep(retry_after)
                    else:
                        # Try refreshing token once more for any request errors
                        self._refresh_token_if_needed()
                        if method == "POST" and body:
                            response = requests.post(current_url, headers=self.headers, json=body)
                        else:
                            response = requests.get(current_url, headers=self.headers)
                        response.raise_for_status()
                        data = response.json()
                        
                        if "value" in data:
                            all_results.extend(data["value"])
                            if handle_paging and "@odata.nextLink" in data:
                                current_url = data["@odata.nextLink"]
                                retry_count = 0
                                break
                            else:
                                return all_results
                        else:
                            return data

            if retry_count > max_retries:
                raise Exception("Max retry attempts reached for {} request to {}".format(method, current_url))

            if not handle_paging or "@odata.nextLink" not in data:
                break

        return all_results

    def make_external_request(self, url: str, method: str = "GET", headers: Optional[Dict] = None, body: Optional[str] = None) -> Union[Dict, str]:
        """Make a request to external APIs (non-Graph).

        Args:
            url: The API endpoint to call
            method: HTTP method (GET, POST, etc.)
            headers: Custom headers for the request
            body: Request body for POST requests

        Returns:
            The response from the API
        """
        retry_count = 0
        max_retries = 3
        
        if headers is None:
            headers = {
                'Authorization': self.headers['Authorization'],
                'User-Agent': 'Microsoft Graph Python Data Exporter',
                'Content-Type': 'application/json'
            }

        while retry_count <= max_retries:
            try:
                if method == "POST":
                    response = requests.post(url, headers=headers, data=body if isinstance(body, str) else json.dumps(body))
                else:
                    response = requests.get(url, headers=headers)
                
                response.raise_for_status()
                
                # Try to parse as JSON, fall back to text
                try:
                    return response.json()
                except json.JSONDecodeError:
                    return response.text
                    
            except requests.exceptions.RequestException as e:
                status_code = e.response.status_code if hasattr(e, 'response') and e.response else None
                if status_code in (429, 503):
                    retry_count += 1
                    retry_after = int(e.response.headers.get('Retry-After', retry_count * 2))
                    print(f"{Fore.YELLOW}⚠️ Rate limited on external request. Retrying in {retry_after} seconds... (Attempt {retry_count}/{max_retries}){Style.RESET_ALL}")
                    time.sleep(retry_after)
                else:
                    raise
        
        raise Exception(f"Max retry attempts reached for {method} request to {url}")

    def save_json_data(self, data: Any, file_path: str) -> None:
        """Save data to a JSON file.

        Args:
            data: The data to save
            file_path: The output file path
        """
        try:
            if data is None:
                print(f"{Fore.YELLOW}  ⚠️ Skipping {os.path.basename(file_path)} - No data to save{Style.RESET_ALL}")
                return

            directory = os.path.dirname(file_path)
            if not os.path.exists(directory):
                os.makedirs(directory)

            with open(file_path, 'w', encoding='utf-8') as f:
                json.dump(data, f, indent=2)
            print(f"{Fore.GREEN}  ✅ Saved: {os.path.basename(file_path)}{Style.RESET_ALL}")

        except Exception as e:
            print(f"{Fore.RED}  ❌ Failed to save: {os.path.basename(file_path)} - {str(e)}{Style.RESET_ALL}")
            raise

    def get_setting_status_errors(self, export_path: str) -> None:
        """Extract Setting Status Errors."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Setting Status Errors...{Style.RESET_ALL}")
            
            # Get version information from Intune portal
            version_response = self.make_external_request("https://intune.microsoft.com/signin/idpRedirect.js", headers={})
            
            if isinstance(version_response, str):
                match = re.search(r'"extensionsPageVersion":(\{[^}]+\})', version_response)
                if match:
                    try:
                        versions = json.loads(match.group(1))
                        device_settings_version = versions['Microsoft_Intune_DeviceSettings'][0]
                        
                        root = "https://afd-v2.hosting.portal.azure.net"
                        setting_status_url = f"{root}/intunedevicesettings/Content/{device_settings_version}/Scripts/DeviceConfiguration/Blades/DevicePoliciesStatus/SettingStatus.js"
                        
                        setting_status_response = self.make_external_request(setting_status_url, headers={})
                        
                        if isinstance(setting_status_response, str):
                            # Look for the SettingStatusErrorMap with more flexible regex
                            match = re.search(r'SettingStatusErrorMap\s*=\s*(\{.*?\});', setting_status_response, re.DOTALL)
                            if match:
                                js_object = match.group(1)
                                try:
                                    # Clean the JavaScript object string for JSON parsing
                                    # Remove control characters and fix common JS->JSON issues
                                    cleaned_js = re.sub(r'[\x00-\x1f\x7f-\x9f]', '', js_object)  # Remove control chars
                                    cleaned_js = re.sub(r'([{,]\s*)([a-zA-Z_$][a-zA-Z0-9_$]*)\s*:', r'\1"\2":', cleaned_js)  # Quote unquoted keys
                                    cleaned_js = re.sub(r"'([^']*)'", r'"\1"', cleaned_js)  # Convert single quotes to double quotes
                                    
                                    error_map = json.loads(cleaned_js)
                                    self.save_json_data(error_map, os.path.join(export_path, "SettingStatusErrors.json"))
                                except json.JSONDecodeError as json_err:
                                    print(f"{Fore.YELLOW}  ⚠️ Failed to parse SettingStatusErrorMap as JSON: {json_err}{Style.RESET_ALL}")
                                    # Save the raw JavaScript object as a text file instead
                                    with open(os.path.join(export_path, "SettingStatusErrors_raw.txt"), 'w', encoding='utf-8') as f:
                                        f.write(js_object)
                                    print(f"{Fore.YELLOW}  ⚠️ Saved raw JavaScript object to SettingStatusErrors_raw.txt{Style.RESET_ALL}")
                            else:
                                print(f"{Fore.YELLOW}  ⚠️ SettingStatusErrorMap not found in response{Style.RESET_ALL}")
                        else:
                            print(f"{Fore.YELLOW}  ⚠️ Failed to fetch setting status script{Style.RESET_ALL}")
                    except (json.JSONDecodeError, KeyError) as e:
                        print(f"{Fore.YELLOW}  ⚠️ Failed to parse version information: {e}{Style.RESET_ALL}")
                else:
                    print(f"{Fore.YELLOW}  ⚠️ Version information not found in redirect script{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}  ⚠️ Failed to fetch Intune redirect script{Style.RESET_ALL}")
            
            print(f"{Fore.GREEN}✅ Setting Status Errors extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Setting Status Errors: {str(e)}{Style.RESET_ALL}")
            raise

    def get_dcv1_policies(self, export_path: str) -> None:
        """Extract Device Configuration v1 Policies."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Device Configuration v1 Policies...{Style.RESET_ALL}")
            
            # Get version information
            version_response = self.make_external_request("https://intune.microsoft.com/signin/idpRedirect.js", headers={})
            
            if isinstance(version_response, str):
                match = re.search(r'"extensionsPageVersion":(\{[^}]+\})', version_response)
                if match:
                    versions = json.loads(match.group(1))
                    device_settings_version = versions['Microsoft_Intune_DeviceSettings'][0]
                    
                    root = "https://afd-v2.hosting.portal.azure.net"
                    root_device_settings = f"{root}/intunedevicesettings/Content/{device_settings_version}/Scripts/DeviceConfiguration"
                    
                    # Clean existing DCv1 directory
                    dcv1_path = os.path.join(export_path, "DCv1")
                    if os.path.exists(dcv1_path):
                        shutil.rmtree(dcv1_path)
                    
                    for source in ['Configuration', 'Compliance']:
                        print(f"{Fore.YELLOW}  📋 Processing {source} metadata...{Style.RESET_ALL}")
                        
                        source_path = os.path.join(dcv1_path, source)
                        os.makedirs(source_path, exist_ok=True)
                        
                        metadata_url = f"{root_device_settings}/Metadata/{source}Metadata.js"
                        metadata_response = self.make_external_request(metadata_url, headers={})
                        
                        if isinstance(metadata_response, str) and 'metadata = ' in metadata_response:
                            match = re.search(r'(?s)metadata = (\{.+\});', metadata_response)
                            if match:
                                try:
                                    metadata = json.loads(match.group(1))
                                    
                                    for family_name, family_data in metadata.items():
                                        if isinstance(family_data, list):
                                            for setting in family_data:
                                                if isinstance(setting, dict) and 'id' in setting:
                                                    clean_id = '_'.join(setting['id'].split('_')[:-1])
                                                    setting['id'] = clean_id
                                                    
                                                    # Clean nested IDs
                                                    setting = self._remove_dcv1_version_suffixes(setting)
                                                    
                                                    file_path = os.path.join(source_path, f"{clean_id}.json")
                                                    self.save_json_data(setting, file_path)
                                except json.JSONDecodeError as e:
                                    print(f"{Fore.YELLOW}  ⚠️ Failed to parse {source} metadata: {e}{Style.RESET_ALL}")
            
            print(f"{Fore.GREEN}✅ Device Configuration v1 extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv1 policies: {str(e)}{Style.RESET_ALL}")
            raise

    def _remove_dcv1_version_suffixes(self, setting: Dict) -> Dict:
        """Remove version suffixes from DCv1 setting IDs."""
        # Clean child settings
        if 'childSettings' in setting and setting['childSettings']:
            for child in setting['childSettings']:
                if 'id' in child:
                    child['id'] = '_'.join(child['id'].split('_')[:-1])
                child = self._remove_dcv1_version_suffixes(child)
        
        # Clean options
        if 'options' in setting and setting['options']:
            for option in setting['options']:
                if 'children' in option and option['children']:
                    for child in option['children']:
                        if 'id' in child:
                            child['id'] = '_'.join(child['id'].split('_')[:-1])
                        child = self._remove_dcv1_version_suffixes(child)
        
        # Clean complex options
        if 'complexOptions' in setting and setting['complexOptions']:
            for complex_option in setting['complexOptions']:
                if 'id' in complex_option:
                    complex_option['id'] = '_'.join(complex_option['id'].split('_')[:-1])
                complex_option = self._remove_dcv1_version_suffixes(complex_option)
        
        # Clean columns
        if 'columns' in setting and setting['columns']:
            for column in setting['columns']:
                if 'metadata' in column and column['metadata']:
                    if 'id' in column['metadata']:
                        column['metadata']['id'] = '_'.join(column['metadata']['id'].split('_')[:-1])
                    column['metadata'] = self._remove_dcv1_version_suffixes(column['metadata'])
        
        return setting

    def get_service_principals_and_endpoints(self, export_path: str) -> None:
        """Extract Service Principals and Endpoints."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Service Principals and Endpoints...{Style.RESET_ALL}")
            
            # Get endpoints
            print(f"{Fore.YELLOW}  📋 Processing endpoints...{Style.RESET_ALL}")
            endpoints_url = f"{self.graph_beta_url}/servicePrincipals/appId=0000000a-0000-0000-c000-000000000000/endpoints"
            endpoints = self.make_request(endpoints_url, handle_paging=False)
            
            if endpoints:
                # Handle both list and dict responses
                if isinstance(endpoints, list):
                    endpoints_data = endpoints
                elif isinstance(endpoints, dict) and "value" in endpoints:
                    endpoints_data = endpoints["value"]
                else:
                    endpoints_data = [endpoints]
                
                endpoints_data.sort(key=lambda x: x.get('capability', ''))
                self.save_json_data(endpoints_data, os.path.join(export_path, "Endpoints.json"))
            else:
                print(f"{Fore.YELLOW}  ⚠️ No endpoints data returned{Style.RESET_ALL}")
            
            # Get service principals
            print(f"{Fore.YELLOW}  📋 Processing service principals...{Style.RESET_ALL}")
            service_principals_path = os.path.join(export_path, "ServicePrincipals")
            if os.path.exists(service_principals_path):
                shutil.rmtree(service_principals_path)
            os.makedirs(service_principals_path)
            
            service_principals = self.make_request(f"{self.graph_beta_url}/servicePrincipals")
            
            if service_principals:
                # Handle both list and dict responses
                if isinstance(service_principals, list):
                    sp_data = service_principals
                elif isinstance(service_principals, dict) and "value" in service_principals:
                    sp_data = service_principals["value"]
                else:
                    sp_data = [service_principals]
                
                for sp in sp_data:
                    if sp.get('appId'):
                        file_path = os.path.join(service_principals_path, f"{sp['appId']}.json")
                        self.save_json_data(sp, file_path)
            else:
                print(f"{Fore.YELLOW}  ⚠️ No service principals data returned{Style.RESET_ALL}")
            
            print(f"{Fore.GREEN}✅ Service Principals and Endpoints extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Service Principals and Endpoints: {str(e)}{Style.RESET_ALL}")
            raise

    def get_role_definitions(self, export_path: str) -> None:
        """Extract Role Definitions."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Role Definitions...{Style.RESET_ALL}")
            
            role_definitions_path = os.path.join(export_path, "RoleDefinitions")
            if os.path.exists(role_definitions_path):
                shutil.rmtree(role_definitions_path)
            os.makedirs(role_definitions_path)
            
            providers = ['cloudPC', 'deviceManagement', 'directory', 'entitlementManagement', 'exchange']
            successful_providers = []
            failed_providers = []
            
            for provider in providers:
                print(f"{Fore.YELLOW}  📋 Processing {provider} role definitions...{Style.RESET_ALL}")
                
                try:
                    provider_path = os.path.join(role_definitions_path, provider)
                    os.makedirs(provider_path)
                    
                    role_def_url = f"{self.graph_beta_url}/roleManagement/{provider}/roleDefinitions"
                    response = self.make_request(role_def_url)
                    
                    # Handle both list and dict responses
                    if isinstance(response, list):
                        role_definitions = response
                    elif isinstance(response, dict) and "value" in response:
                        role_definitions = response["value"]
                    else:
                        role_definitions = [response] if response else []
                    
                    if role_definitions:
                        for role_def in role_definitions:
                            if isinstance(role_def, dict) and "id" in role_def:
                                role_id = role_def["id"]
                                file_path = os.path.join(provider_path, f"{role_id}.json")
                                self.save_json_data(role_def, file_path)
                        successful_providers.append(provider)
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ No role definitions found for {provider}{Style.RESET_ALL}")
                        successful_providers.append(provider)
                except Exception as e:
                    error_message = str(e)
                    if "Authorization_RequestDenied" in error_message or "Insufficient privileges" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {provider} role definitions - skipping{Style.RESET_ALL}")
                        failed_providers.append(provider)
                    else:
                        print(f"{Fore.RED}    ❌ Failed to extract {provider} role definitions: {error_message}{Style.RESET_ALL}")
                        failed_providers.append(provider)
            
            if successful_providers:
                print(f"{Fore.GREEN}✅ Role Definitions extraction completed for: {', '.join(successful_providers)}{Style.RESET_ALL}")
            
            if failed_providers:
                print(f"{Fore.YELLOW}⚠️ Role Definitions extraction failed for: {', '.join(failed_providers)} (insufficient permissions){Style.RESET_ALL}")
                
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Role Definitions: {str(e)}{Style.RESET_ALL}")
            raise

    def get_resource_operations(self, export_path: str) -> None:
        """Extract Resource Operations."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Resource Operations...{Style.RESET_ALL}")
            
            resource_operations_path = os.path.join(export_path, "ResourceOperations")
            if os.path.exists(resource_operations_path):
                shutil.rmtree(resource_operations_path)
            os.makedirs(resource_operations_path)
            
            resource_ops_url = f"{self.graph_beta_url}/deviceManagement/resourceOperations"
            resource_operations = self.make_request(resource_ops_url)
            
            # Handle both list and dict responses
            if isinstance(resource_operations, list):
                operations_data = resource_operations
            elif isinstance(resource_operations, dict) and "value" in resource_operations:
                operations_data = resource_operations["value"]
            else:
                operations_data = [resource_operations] if resource_operations else []
            
            for operation in operations_data:
                operation_id = operation['id']
                file_path = os.path.join(resource_operations_path, f"{operation_id}.json")
                self.save_json_data(operation, file_path)
            
            print(f"{Fore.GREEN}✅ Resource Operations extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Resource Operations: {str(e)}{Style.RESET_ALL}")
            raise

    def get_defender_hunting_tables(self, export_path: str) -> None:
        """Extract Defender Hunting Tables."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Defender Hunting Table Schemas...{Style.RESET_ALL}")
            
            defender_path = os.path.join(export_path, "Defender")
            if os.path.exists(defender_path):
                shutil.rmtree(defender_path)
            os.makedirs(defender_path)
            
            hunting_tables = [
                'AlertEvidence', 'AlertInfo', 'BehaviorEntities', 'BehaviorInfo',
                'AADSignInEventsBeta', 'AADSpnSignInEventsBeta', 'CloudAppEvents',
                'IdentityInfo', 'IdentityLogonEvents', 'EmailAttachmentInfo',
                'EmailEvents', 'EmailPostDeliveryEvents', 'EmailUrlInfo',
                'UrlClickEvents', 'ExposureGraphEdges', 'ExposureGraphNodes'
            ]
            
            successful_tables = []
            failed_tables = []
            
            for table in hunting_tables:
                print(f"{Fore.YELLOW}  📋 Processing {table} schema...{Style.RESET_ALL}")
                
                query = f"{table} | getschema | project Description=\"\", Type=split(DataType, \".\")[1], Entity=\"\", Name=ColumnName"
                
                request_body = {
                    "Query": query
                }
                
                try:
                    response = self.make_request(
                        f"{self.graph_beta_url}/security/runHuntingQuery",
                        method="POST",
                        body=request_body,
                        handle_paging=False
                    )
                    
                    if response and (response.get('results') or response.get('Results')):
                        results = response.get('results') or response.get('Results')
                        self.save_json_data(results, os.path.join(defender_path, f"{table}.json"))
                        successful_tables.append(table)
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ No results returned for {table}{Style.RESET_ALL}")
                        failed_tables.append(table)
                except Exception as e:
                    error_message = str(e)
                    if "Authorization_RequestDenied" in error_message or "Insufficient privileges" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {table} - skipping{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ Failed to get schema for {table}: {error_message}{Style.RESET_ALL}")
                    failed_tables.append(table)
            
            if successful_tables:
                print(f"{Fore.GREEN}✅ Defender Hunting Tables extraction completed for: {len(successful_tables)}/{len(hunting_tables)} tables{Style.RESET_ALL}")
            
            if failed_tables:
                print(f"{Fore.YELLOW}⚠️ Some Defender tables could not be accessed (permissions or availability): {len(failed_tables)} tables{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Defender Hunting Tables: {str(e)}{Style.RESET_ALL}")
            print(f"{Fore.YELLOW}⚠️ Defender hunting tables extraction will be skipped{Style.RESET_ALL}")

    def get_dcv2_configuration_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Settings."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Settings...{Style.RESET_ALL}")
            
            dcv2_path = os.path.join(export_path, "DCv2")
            if os.path.exists(dcv2_path):
                shutil.rmtree(dcv2_path)
            
            # Configuration Settings (Settings Catalog)
            print(f"{Fore.YELLOW}  📋 Processing configuration settings...{Style.RESET_ALL}")
            settings_path = os.path.join(dcv2_path, "Settings")
            os.makedirs(settings_path, exist_ok=True)
            
            config_settings_url = f"{self.graph_beta_url}/deviceManagement/configurationSettings"
            config_settings = self.make_request(config_settings_url)
            
            for setting in config_settings:
                if 'version' in setting:
                    del setting['version']
                setting_id = setting['id']
                file_path = os.path.join(settings_path, f"{setting_id}.json")
                self.save_json_data(setting, file_path)
            
            # Create backwards compatibility folder
            backwards_compat_path = os.path.join(export_path, "settings")
            if os.path.exists(backwards_compat_path):
                shutil.rmtree(backwards_compat_path)
            shutil.copytree(settings_path, backwards_compat_path)
            
            print(f"{Fore.GREEN}✅ DCv2 Configuration Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Configuration Settings: {str(e)}{Style.RESET_ALL}")
            raise

    def get_dcv2_compliance_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Compliance Settings."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Compliance Settings...{Style.RESET_ALL}")
            
            compliance_path = os.path.join(export_path, "DCv2", "Compliance")
            os.makedirs(compliance_path, exist_ok=True)
            
            compliance_settings_url = f"{self.graph_beta_url}/deviceManagement/complianceSettings"
            compliance_settings = self.make_request(compliance_settings_url)
            
            for setting in compliance_settings:
                if 'version' in setting:
                    del setting['version']
                setting_id = setting['id']
                file_path = os.path.join(compliance_path, f"{setting_id}.json")
                self.save_json_data(setting, file_path)
            
            print(f"{Fore.GREEN}✅ DCv2 Compliance Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Compliance Settings: {str(e)}{Style.RESET_ALL}")
            raise

    def get_dcv2_policy_templates(self, export_path: str) -> None:
        """Extract Device Configuration v2 Policy Templates."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Policy Templates...{Style.RESET_ALL}")
            
            templates_path = os.path.join(export_path, "DCv2", "Templates")
            os.makedirs(templates_path, exist_ok=True)
            
            templates_url = f"{self.graph_beta_url}/deviceManagement/configurationPolicyTemplates"
            templates = self.make_request(templates_url)
            
            for template in templates:
                template_id = template['id']
                file_path = os.path.join(templates_path, f"{template_id}.json")
                self.save_json_data(template, file_path)
            
            print(f"{Fore.GREEN}✅ DCv2 Policy Templates extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Policy Templates: {str(e)}{Style.RESET_ALL}")
            raise

    def get_dcv2_inventory_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Inventory Settings."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Inventory Settings...{Style.RESET_ALL}")
            
            inventory_path = os.path.join(export_path, "DCv2", "Inventory")
            os.makedirs(inventory_path, exist_ok=True)
            
            inventory_url = f"{self.graph_beta_url}/deviceManagement/inventorySettings"
            inventory_settings = self.make_request(inventory_url)
            
            for setting in inventory_settings:
                if 'version' in setting:
                    del setting['version']
                setting_id = setting['id']
                file_path = os.path.join(inventory_path, f"{setting_id}.json")
                self.save_json_data(setting, file_path)
            
            print(f"{Fore.GREEN}✅ DCv2 Inventory Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Inventory Settings: {str(e)}{Style.RESET_ALL}")
            raise

    def get_identity_product_changes(self, export_path: str) -> None:
        """Extract Identity Product Changes."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Identity Product Changes...{Style.RESET_ALL}")
            
            changes_url = f"{self.graph_beta_url}/identity/productChanges"
            all_changes = self.make_request(changes_url, handle_paging=True)
            
            if all_changes:
                self.save_json_data(all_changes, os.path.join(export_path, "IdentityProductChanges.json"))
                print(f"{Fore.GREEN}✅ Identity Product Changes extraction completed ({len(all_changes)} items){Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}⚠️ No Identity Product Changes data returned{Style.RESET_ALL}")
            
        except Exception as e:
            error_message = str(e)
            if "Authorization_RequestDenied" in error_message or "Insufficient privileges" in error_message or "403" in error_message:
                print(f"{Fore.YELLOW}⚠️ Insufficient privileges for Identity Product Changes - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.RED}❌ Failed to extract Identity Product Changes: {error_message}{Style.RESET_ALL}")
            # Don't throw - this is optional data
            print(f"{Fore.YELLOW}⚠️ Identity Product Changes extraction will be skipped{Style.RESET_ALL}")

    def get_office_cloud_policy_service(self, export_path: str) -> None:
        """Extract Office Cloud Policy Service (OCPS) Data."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Office Cloud Policy Service (OCPS) Data...{Style.RESET_ALL}")
            
            ocps_path = os.path.join(export_path, "OCPS")
            if os.path.exists(ocps_path):
                shutil.rmtree(ocps_path)
            os.makedirs(ocps_path)
            
            # Note: OCPS typically requires specific authentication, but we'll try with Graph token
            headers = {
                'Authorization': self.headers['Authorization'],
                'User-Agent': 'Microsoft Graph Python Data Exporter',
                'Content-Type': 'application/json'
            }
            
            ocps_endpoints = {
                'synchealth': 'https://clients.config.office.net/odbhealth/v1.0/synchealth/reports/versioncount',
                'languages': 'https://clients.config.office.net/releases/v1.0/FileList/languagesForProductIds?productId=O365ProPlusRetail',
                'userflights': 'https://config.office.com/appConfig/v1.0/userflights',
                'SettingsCatalog': 'https://clients.config.office.net/settings/v1.0/SettingsCatalog/Settings',
                'ServiceHealth': 'https://config.office.com/appConfig/v1.0/ServiceHealth',
                'OfficeReleases': 'https://clients.config.office.net/releases/v1.0/OfficeReleases'
            }
            
            successful_endpoints = []
            failed_endpoints = []
            
            for endpoint_name, endpoint_url in ocps_endpoints.items():
                print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
                
                try:
                    response = self.make_external_request(endpoint_url, headers=headers)
                    self.save_json_data(response, os.path.join(ocps_path, f"{endpoint_name}.json"))
                    successful_endpoints.append(endpoint_name)
                except Exception as e:
                    error_message = str(e)
                    if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
                    failed_endpoints.append(endpoint_name)
            
            # Try to get Feature data with different endpoint
            try:
                print(f"{Fore.YELLOW}  📋 Processing FeatureData...{Style.RESET_ALL}")
                feature_response = self.make_external_request(
                    'https://clients.config.office.net/onboarding/odata/v1.0/FeatureData', 
                    headers=headers
                )
                
                if isinstance(feature_response, dict) and 'value' in feature_response:
                    self.save_json_data(feature_response['value'], os.path.join(ocps_path, "FeatureData.json"))
                    successful_endpoints.append("FeatureData")
                else:
                    self.save_json_data(feature_response, os.path.join(ocps_path, "FeatureData.json"))
                    successful_endpoints.append("FeatureData")
            except Exception as e:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract FeatureData: {str(e)}{Style.RESET_ALL}")
                failed_endpoints.append("FeatureData")
            
            if successful_endpoints:
                print(f"{Fore.GREEN}✅ OCPS extraction completed for: {', '.join(successful_endpoints)}{Style.RESET_ALL}")
            
            if failed_endpoints:
                print(f"{Fore.YELLOW}⚠️ OCPS extraction failed for: {', '.join(failed_endpoints)} (authentication or availability){Style.RESET_ALL}")
            
            if not successful_endpoints:
                print(f"{Fore.YELLOW}⚠️ No OCPS data could be extracted - likely requires specific Office 365 authentication{Style.RESET_ALL}")
            
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract OCPS data: {str(e)}{Style.RESET_ALL}")
            # Don't throw - this is optional data that may not be accessible
            print(f"{Fore.YELLOW}⚠️ OCPS extraction will be skipped{Style.RESET_ALL}")

    def get_teams_admin_center(self, export_path: str) -> None:
        """Extract Teams Admin Center Data."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting Teams Admin Center Data...{Style.RESET_ALL}")
            
            teams_path = os.path.join(export_path, "Teams")
            if os.path.exists(teams_path):
                shutil.rmtree(teams_path)
            os.makedirs(teams_path)
            
            # Note: Teams Admin Center typically requires specific authentication
            headers = {
                'Authorization': self.headers['Authorization'],
                'User-Agent': 'Microsoft Graph Python Data Exporter',
                'Content-Type': 'application/json'
            }
            
            teams_endpoints = {
                'authz': {
                    'url': 'https://authsvc.teams.microsoft.com/v1.0/authz',
                    'method': 'POST'
                },
                'appsCatalog': {
                    'url': 'https://teams.microsoft.com/api/mt/part/au-01/beta/users/appsCatalog',
                    'method': 'GET'
                }
            }
            
            successful_endpoints = []
            failed_endpoints = []
            
            for endpoint_name, endpoint_config in teams_endpoints.items():
                print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
                
                try:
                    if endpoint_config['method'] == 'POST':
                        response = self.make_external_request(
                            endpoint_config['url'], 
                            method="POST", 
                            headers=headers, 
                            body="{}"
                        )
                    else:
                        response = self.make_external_request(endpoint_config['url'], headers=headers)
                    
                    # Remove sensitive token data if present
                    if isinstance(response, dict) and 'tokens' in response:
                        del response['tokens']
                    
                    self.save_json_data(response, os.path.join(teams_path, f"{endpoint_name}.json"))
                    successful_endpoints.append(endpoint_name)
                except Exception as e:
                    error_message = str(e)
                    if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
                    failed_endpoints.append(endpoint_name)
            
            if successful_endpoints:
                print(f"{Fore.GREEN}✅ Teams Admin Center extraction completed for: {', '.join(successful_endpoints)}{Style.RESET_ALL}")
            
            if failed_endpoints:
                print(f"{Fore.YELLOW}⚠️ Teams Admin Center extraction failed for: {', '.join(failed_endpoints)} (authentication or availability){Style.RESET_ALL}")
            
            if not successful_endpoints:
                print(f"{Fore.YELLOW}⚠️ No Teams Admin Center data could be extracted - likely requires specific Teams authentication{Style.RESET_ALL}")
            
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Teams Admin Center data: {str(e)}{Style.RESET_ALL}")
            # Don't throw - this is optional data that may not be accessible
            print(f"{Fore.YELLOW}⚠️ Teams Admin Center extraction will be skipped{Style.RESET_ALL}")

    def get_m365_admin_center(self, export_path: str) -> None:
        """Extract M365 Admin Center Data."""
        try:
            print(f"{Fore.CYAN}🔍 Extracting M365 Admin Center Data...{Style.RESET_ALL}")
            
            m365_admin_path = os.path.join(export_path, "M365Admin")
            if os.path.exists(m365_admin_path):
                shutil.rmtree(m365_admin_path)
            os.makedirs(m365_admin_path)
            
            # Note: M365 Admin Center typically requires specific authentication
            headers = {
                'Authorization': self.headers['Authorization'],
                'User-Agent': 'Microsoft Graph Python Data Exporter',
                'Content-Type': 'application/json'
            }
            
            # Simple endpoints (single response)
            simple_endpoints = {
                'features-config': 'https://admin.microsoft.com/admin/api/features/config',
                'features-all': 'https://admin.microsoft.com/admin/api/features/all',
                'partner-list': 'https://admin.microsoft.com/fd/bcws/api/v1/IntraTenantPartner/getPartnerList',
                'product-offers': 'https://admin.microsoft.com/fd/bsxcommerce/v1/ProductOffers/EligibleProductOffers?language=en-US'
            }
            
            # Complex endpoints (with nested data extraction)
            complex_endpoints = {
                'apps': {
                    'url': 'https://admin.microsoft.com/fd/addins/api/availableApps?workloads=MetaOS,Teams',
                    'property': 'apps'
                },
                'policy_definitions': {
                    'url': 'https://admin.microsoft.com/fd/edgeenterpriseextensionsmanagement/api/policies',
                    'property': 'policy_definitions'
                },
                'C2RReleaseInfo': {
                    'url': 'https://admin.microsoft.com/fd/dms/odata/C2RReleaseInfo',
                    'property': 'value'
                },
                'ProductOfferIndex': {
                    'url': 'https://admin.microsoft.com/fd/bsxcommerce/v1/ProductOfferIndex?language=en-US',
                    'property': 'results'
                },
                'licensedProducts': {
                    'url': 'https://admin.microsoft.com/fd/m365licensing/v3/licensedProducts',
                    'property': 'value'
                },
                'sidebarExtensions': {
                    'url': 'https://admin.microsoft.com/fd/edgeenterpriseextensionsmanagement/api/sidebarExtensions',
                    'property': 'hub_apps'
                }
            }
            
            successful_endpoints = []
            failed_endpoints = []
            
            # Process simple endpoints
            for endpoint_name, endpoint_url in simple_endpoints.items():
                print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
                
                try:
                    response = self.make_external_request(endpoint_url, headers=headers)
                    self.save_json_data(response, os.path.join(m365_admin_path, f"{endpoint_name}.json"))
                    successful_endpoints.append(endpoint_name)
                except Exception as e:
                    error_message = str(e)
                    if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
                    failed_endpoints.append(endpoint_name)
            
            # Process complex endpoints
            for endpoint_name, endpoint_config in complex_endpoints.items():
                print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
                
                try:
                    response = self.make_external_request(endpoint_config['url'], headers=headers)
                    
                    # Extract specific property if defined
                    if endpoint_config['property'] and isinstance(response, dict) and endpoint_config['property'] in response:
                        data_to_save = response[endpoint_config['property']]
                    else:
                        data_to_save = response
                    
                    # Extract filename from URL
                    url_parts = [part for part in endpoint_config['url'].split('/') if part and 
                                'https:' not in part and 'admin.microsoft.com' not in part and 
                                'fd' not in part and 'api' not in part]
                    file_name = url_parts[-1].split('?')[0] if url_parts else endpoint_name
                    if not file_name:
                        file_name = endpoint_name
                    
                    self.save_json_data(data_to_save, os.path.join(m365_admin_path, f"{file_name}.json"))
                    successful_endpoints.append(endpoint_name)
                except Exception as e:
                    error_message = str(e)
                    if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                        print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
                    failed_endpoints.append(endpoint_name)
            
            # Process Service Health (special handling)
            try:
                print(f"{Fore.YELLOW}  📋 Processing ServiceHealth...{Style.RESET_ALL}")
                service_health_response = self.make_external_request(
                    'https://admin.microsoft.com/admin/api/servicehealth/status/activeCM?showResolved=true', 
                    headers=headers
                )
                
                if isinstance(service_health_response, dict) and 'ServiceStatus' in service_health_response:
                    flattened_data = []
                    for service in service_health_response['ServiceStatus']:
                        if 'MessagesByClassification' in service:
                            if 'Incidents' in service['MessagesByClassification']:
                                flattened_data.extend(service['MessagesByClassification']['Incidents'])
                            if 'Advisories' in service['MessagesByClassification']:
                                flattened_data.extend(service['MessagesByClassification']['Advisories'])
                    self.save_json_data(flattened_data, os.path.join(m365_admin_path, "ServiceHealth.json"))
                    successful_endpoints.append("ServiceHealth")
            except Exception as e:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract ServiceHealth: {str(e)}{Style.RESET_ALL}")
                failed_endpoints.append("ServiceHealth")
            
            # Process Message Center (special handling)
            try:
                print(f"{Fore.YELLOW}  📋 Processing MessageCenter...{Style.RESET_ALL}")
                message_center_response = self.make_external_request(
                    'https://admin.microsoft.com/admin/api/messagecenter', 
                    headers=headers
                )
                
                if isinstance(message_center_response, dict) and 'Messages' in message_center_response:
                    # Remove sort-specific properties
                    for message in message_center_response['Messages']:
                        if 'ActionRequiredBySortValue' in message:
                            del message['ActionRequiredBySortValue']
                    self.save_json_data(message_center_response['Messages'], os.path.join(m365_admin_path, "messagecenter.json"))
                    successful_endpoints.append("MessageCenter")
            except Exception as e:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract MessageCenter: {str(e)}{Style.RESET_ALL}")
                failed_endpoints.append("MessageCenter")
            
            # Process Concierge Config (special handling)
            try:
                print(f"{Fore.YELLOW}  📋 Processing ConciergeConfig...{Style.RESET_ALL}")
                concierge_response = self.make_external_request(
                    'https://admin.microsoft.com/api/concierge/GetConciergeConfig', 
                    headers=headers
                )
                
                # Remove SessionID for security
                if isinstance(concierge_response, dict) and 'SessionID' in concierge_response:
                    del concierge_response['SessionID']
                
                self.save_json_data(concierge_response, os.path.join(m365_admin_path, "GetConciergeConfig.json"))
                successful_endpoints.append("ConciergeConfig")
            except Exception as e:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract ConciergeConfig: {str(e)}{Style.RESET_ALL}")
                failed_endpoints.append("ConciergeConfig")
            
            if successful_endpoints:
                print(f"{Fore.GREEN}✅ M365 Admin Center extraction completed for: {len(successful_endpoints)} endpoints{Style.RESET_ALL}")
            
            if failed_endpoints:
                print(f"{Fore.YELLOW}⚠️ M365 Admin Center extraction failed for: {len(failed_endpoints)} endpoints (authentication or availability){Style.RESET_ALL}")
            
            if not successful_endpoints:
                print(f"{Fore.YELLOW}⚠️ No M365 Admin Center data could be extracted - likely requires specific M365 Admin authentication{Style.RESET_ALL}")
            
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract M365 Admin Center data: {str(e)}{Style.RESET_ALL}")
            # Don't throw - this is optional data that may not be accessible
            print(f"{Fore.YELLOW}⚠️ M365 Admin Center extraction will be skipped{Style.RESET_ALL}")

    def test_service_principal_permissions(self) -> bool:
        """Validate service principal has required Microsoft Graph permissions.
        
        Performs a preflight check to ensure the authenticated service principal
        has all required application permissions and that admin consent has been granted.
        
        Returns:
            bool: True if all required permissions are granted, False otherwise
            
        Raises:
            requests.exceptions.RequestException: If permission validation API calls fail
        """
        try:
            print(f"\n{Fore.YELLOW}🔐 Validating service principal permissions...{Style.RESET_ALL}")
            
            # Get the service principal for this application
            sp_filter = f"appId eq '{self.client_id}'"
            encoded_filter = urllib.parse.quote(sp_filter)
            sp_url = f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_filter}"
            
            sp_data = self.make_request(sp_url)
            
            # Handle both list and dict responses from make_request
            if isinstance(sp_data, list):
                service_principals = sp_data
            elif isinstance(sp_data, dict) and "value" in sp_data:
                service_principals = sp_data["value"]
            else:
                service_principals = [sp_data] if sp_data else []
            
            if not service_principals:
                print(f"{Fore.RED}❌ Service principal not found for application: {self.client_id}{Style.RESET_ALL}")
                return False
            
            service_principal = service_principals[0]
            sp_id = service_principal['id']
            
            print(f"   {Fore.GREEN}✅ Found service principal: {service_principal['displayName']}{Style.RESET_ALL}")
            
            # Get Microsoft Graph service principal (resource)
            graph_sp_filter = "appId eq '00000003-0000-0000-c000-000000000000'"
            encoded_graph_filter = urllib.parse.quote(graph_sp_filter)
            graph_sp_url = f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_graph_filter}"
            
            graph_sp_data = self.make_request(graph_sp_url)
            
            # Handle both list and dict responses
            if isinstance(graph_sp_data, list):
                graph_service_principals = graph_sp_data
            elif isinstance(graph_sp_data, dict) and "value" in graph_sp_data:
                graph_service_principals = graph_sp_data["value"]
            else:
                graph_service_principals = [graph_sp_data] if graph_sp_data else []
            
            if not graph_service_principals:
                print(f"{Fore.RED}❌ Microsoft Graph service principal not found{Style.RESET_ALL}")
                return False
            
            graph_service_principal = graph_service_principals[0]
            
            # Get app role assignments for the service principal
            assignments_url = f"https://graph.microsoft.com/v1.0/servicePrincipals/{sp_id}/appRoleAssignments"
            assignments_data = self.make_request(assignments_url)
            
            # Handle both list and dict responses
            if isinstance(assignments_data, list):
                assignments = assignments_data
            elif isinstance(assignments_data, dict) and "value" in assignments_data:
                assignments = assignments_data["value"]
            else:
                assignments = [assignments_data] if assignments_data else []
            
            # Build a map of role names to role IDs for Microsoft Graph
            role_map = {}
            for role in graph_service_principal['appRoles']:
                role_map[role['value']] = role['id']
            
            print(f"\n   {Fore.CYAN}📋 Checking required permissions:{Style.RESET_ALL}")
            
            all_permissions_present = True
            missing_permissions = []
            
            for permission in self.required_permissions:
                role_id = role_map.get(permission)
                
                if not role_id:
                    print(f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Unknown permission{Style.RESET_ALL}")
                    all_permissions_present = False
                    missing_permissions.append(permission)
                    continue
                
                # Check if this permission is assigned
                assignment = next((a for a in assignments 
                                 if a['resourceId'] == graph_service_principal['id'] and a['appRoleId'] == role_id), None)
                
                if assignment:
                    print(f"   {Fore.GREEN}✅ {Fore.WHITE}{permission}{Fore.GREEN} - Granted{Style.RESET_ALL}")
                else:
                    print(f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Not granted{Style.RESET_ALL}")
                    all_permissions_present = False
                    missing_permissions.append(permission)
            
            if all_permissions_present:
                print(f"\n   {Fore.GREEN}✅ All required permissions are present{Style.RESET_ALL}")
                return True
            else:
                missing_perms_str = ', '.join(missing_permissions)
                print(f"\n   {Fore.RED}❌ Missing required permissions: {missing_perms_str}{Style.RESET_ALL}")
                return False
                
        except requests.exceptions.RequestException as e:
            print(f"{Fore.RED}❌ Failed to validate service principal permissions: {e}{Style.RESET_ALL}")
            return False

    def extract_all_metadata(self, export_path: str) -> None:
        """Extract all metadata from various sources."""
        try:
            print(f"\n{Fore.CYAN}📊 Microsoft Graph Data Extraction Tool with M365 Metadata{Style.RESET_ALL}")
            print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{Style.RESET_ALL}")
            
            # Create output directory
            if os.path.exists(export_path):
                print(f"{Fore.YELLOW}⚠️ Export directory exists. Some data may be overwritten.{Style.RESET_ALL}")
            else:
                os.makedirs(export_path)
                print(f"{Fore.GREEN}✅ Created export directory: {export_path}{Style.RESET_ALL}")
            
            # Execute data extraction modules
            print(f"\n{Fore.CYAN}🚀 Starting data extraction...{Style.RESET_ALL}")
            
            successful_modules = []
            failed_modules = []
            
            # Module 1: Setting Status Errors
            try:
                self.get_setting_status_errors(export_path)
                successful_modules.append("Setting Status Errors")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Setting Status Errors module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Setting Status Errors")
            
            # Module 2: DCv1 Policies
            try:
                self.get_dcv1_policies(export_path)
                successful_modules.append("DCv1 Policies")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ DCv1 Policies module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("DCv1 Policies")
            
            # Module 3: Service Principals and Endpoints
            try:
                self.get_service_principals_and_endpoints(export_path)
                successful_modules.append("Service Principals and Endpoints")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Service Principals and Endpoints module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Service Principals and Endpoints")
            
            # Module 4: Role Definitions
            try:
                self.get_role_definitions(export_path)
                successful_modules.append("Role Definitions")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Role Definitions module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Role Definitions")
            
            # Module 5: Resource Operations
            try:
                self.get_resource_operations(export_path)
                successful_modules.append("Resource Operations")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Resource Operations module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Resource Operations")
            
            # Module 6: Defender Hunting Tables
            try:
                self.get_defender_hunting_tables(export_path)
                successful_modules.append("Defender Hunting Tables")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Defender Hunting Tables module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Defender Hunting Tables")
            
            # Module 7: DCv2 Configuration Settings
            try:
                self.get_dcv2_configuration_settings(export_path)
                successful_modules.append("DCv2 Configuration Settings")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ DCv2 Configuration Settings module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("DCv2 Configuration Settings")
            
            # Module 8: DCv2 Compliance Settings
            try:
                self.get_dcv2_compliance_settings(export_path)
                successful_modules.append("DCv2 Compliance Settings")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ DCv2 Compliance Settings module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("DCv2 Compliance Settings")
            
            # Module 9: DCv2 Policy Templates
            try:
                self.get_dcv2_policy_templates(export_path)
                successful_modules.append("DCv2 Policy Templates")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ DCv2 Policy Templates module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("DCv2 Policy Templates")
            
            # Module 10: DCv2 Inventory Settings
            try:
                self.get_dcv2_inventory_settings(export_path)
                successful_modules.append("DCv2 Inventory Settings")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ DCv2 Inventory Settings module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("DCv2 Inventory Settings")
            
            # Module 11: Identity Product Changes
            try:
                self.get_identity_product_changes(export_path)
                successful_modules.append("Identity Product Changes")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Identity Product Changes module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Identity Product Changes")
            
            # Module 12: Office Cloud Policy Service (OCPS)
            try:
                self.get_office_cloud_policy_service(export_path)
                successful_modules.append("Office Cloud Policy Service")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Office Cloud Policy Service module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Office Cloud Policy Service")
            
            # Module 13: Teams Admin Center
            try:
                self.get_teams_admin_center(export_path)
                successful_modules.append("Teams Admin Center")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ Teams Admin Center module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("Teams Admin Center")
            
            # Module 14: M365 Admin Center
            try:
                self.get_m365_admin_center(export_path)
                successful_modules.append("M365 Admin Center")
            except Exception:
                print(f"{Fore.YELLOW}⚠️ M365 Admin Center module failed - continuing with next module{Style.RESET_ALL}")
                failed_modules.append("M365 Admin Center")
            
            print(f"\n{Fore.GREEN}✨ Data extraction process completed!{Style.RESET_ALL}")
            print(f"{Fore.CYAN}📁 All data saved to: {export_path}{Style.RESET_ALL}")
            
            # Display detailed summary
            if successful_modules:
                print(f"\n{Fore.GREEN}✅ Successful modules ({len(successful_modules)}):{Style.RESET_ALL}")
                for module in successful_modules:
                    print(f"{Fore.GREEN}   • {module}{Style.RESET_ALL}")
            
            if failed_modules:
                print(f"\n{Fore.YELLOW}⚠️ Failed modules ({len(failed_modules)}):{Style.RESET_ALL}")
                for module in failed_modules:
                    print(f"{Fore.YELLOW}   • {module}{Style.RESET_ALL}")
                print(f"\n{Fore.YELLOW}Note: Some failures may be due to insufficient permissions or feature availability.{Style.RESET_ALL}")
            
            # Display file summary
            try:
                subfolders = len([f for f in os.listdir(export_path) if os.path.isdir(os.path.join(export_path, f))])
                files = sum([len([f for f in os.listdir(os.path.join(export_path, d)) if os.path.isfile(os.path.join(export_path, d, f))]) for d in os.listdir(export_path) if os.path.isdir(os.path.join(export_path, d))])
                files += len([f for f in os.listdir(export_path) if os.path.isfile(os.path.join(export_path, f))])
                
                print(f"\n{Fore.CYAN}📊 Extraction Summary:{Style.RESET_ALL}")
                print(f"{Fore.WHITE}   • Total directories: {subfolders}{Style.RESET_ALL}")
                print(f"{Fore.WHITE}   • Total files: {files}{Style.RESET_ALL}")
                print(f"{Fore.WHITE}   • Export path: {export_path}{Style.RESET_ALL}")
            except Exception:
                print(f"\n{Fore.CYAN}📊 Data extraction completed (summary unavailable){Style.RESET_ALL}")
            
        except Exception as e:
            print(f"\n{Fore.RED}❌ Data extraction process failed: {str(e)}{Style.RESET_ALL}")
            raise

def main():
    """Main function to run the script."""
    parser = argparse.ArgumentParser(
        description="Microsoft Graph Metadata Extraction Tool",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument(
        "--tenant-id",
        required=True,
        help="Specify the Entra ID tenant ID (Directory ID) where the application is registered"
    )
    
    parser.add_argument(
        "--client-id", 
        required=True,
        help="Specify the application (client) ID of the Entra ID app registration"
    )
    
    parser.add_argument(
        "--client-secret",
        required=True,
        help="Specify the client secret of the Entra ID app registration"
    )
    
    parser.add_argument(
        "--export-path",
        default="./GraphMetadata",
        help="Export directory path (default: ./GraphMetadata)"
    )
    
    parser.add_argument(
        "--required-permissions",
        nargs="*",
        default=None,
        help="Required Microsoft Graph application permissions to validate (uses defaults if not specified)"
    )
    
    args = parser.parse_args()
    
    try:
        print(f"{Fore.CYAN}🚀 Starting Microsoft Graph Data Extraction with M365 Metadata...{Style.RESET_ALL}")
        print(f"{Fore.CYAN}📁 Export location: {args.export_path}{Style.RESET_ALL}")
        
        extractor = GraphMetadataExtractor(args.tenant_id, args.client_id, args.client_secret, args.required_permissions)
        
        # Authenticate
        if not extractor.authenticate():
            print(f"{Fore.RED}❌ Authentication failed. Please check your credentials.{Style.RESET_ALL}")
            sys.exit(2)
            
        # Validate permissions
        if not extractor.test_service_principal_permissions():
            print(f"{Fore.RED}❌ Cannot proceed due to missing permissions. Please grant the required permissions to your application.{Style.RESET_ALL}")
            sys.exit(2)
        
        extractor.extract_all_metadata(args.export_path)
        
        print(f"\n{Fore.GREEN}🎉 Script completed!{Style.RESET_ALL}")
        print(f"{Fore.GREEN}📋 Microsoft Graph and M365 metadata extraction process finished{Style.RESET_ALL}")
        print(f"{Fore.CYAN}💡 Check the summary above for any modules that may have failed due to permissions{Style.RESET_ALL}")
        print(f"{Fore.CYAN}🔑 Note: M365-specific modules may require additional authentication methods beyond Graph API{Style.RESET_ALL}")
        sys.exit(0)
        
    except Exception as e:
        print(f"\n{Fore.RED}💥 Script execution encountered errors!{Style.RESET_ALL}")
        print(f"{Fore.RED}Error: {str(e)}{Style.RESET_ALL}")
        print(f"\n{Fore.CYAN}💡 Some data may have been successfully extracted. Check your export directory: {args.export_path}{Style.RESET_ALL}")
        print(f"{Fore.CYAN}🔑 Note: M365-specific modules typically require specialized authentication tokens{Style.RESET_ALL}")
        sys.exit(2)

if __name__ == "__main__":
    main()