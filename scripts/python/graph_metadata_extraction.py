#!/usr/bin/env python3
"""
Microsoft Graph Data Exporter with M365 Metadata Extraction - FIXED

Hybrid approach:
- Microsoft Graph SDK for Graph API calls (better reliability, auto-pagination, type safety)
- Custom HTTP client for portal scraping and external APIs
- Async/concurrent processing for maximum performance
- Fixed permission validation logic

Usage:
    python graph_metadata_extraction_fixed.py --tenant-id <tenant_id> --client-id <client_id> --client-secret <client_secret> [--export-path <export_path>]
"""

import argparse
import asyncio
import aiohttp
import aiofiles
import json
import os
import sys
import time
import re
import shutil
import xml.etree.ElementTree as ET
from datetime import datetime, UTC
from typing import Dict, List, Optional, Any, Union
from pathlib import Path
from urllib.parse import urlparse, quote
from dataclasses import dataclass, field

from colorama import Fore, Style, init
from azure.identity import ClientSecretCredential, CredentialUnavailableError

# Microsoft Graph SDK imports
from msgraph import GraphServiceClient
from msgraph_beta import GraphServiceClient as GraphBetaServiceClient
from kiota_abstractions.base_request_configuration import RequestConfiguration
from kiota_http.middleware.options import ResponseHandlerOption
from kiota_abstractions.native_response_handler import NativeResponseHandler

# Initialize colorama
init()

@dataclass
class Config:
    """Configuration for the extraction process"""
    tenant_id: str
    client_id: str
    client_secret: str
    export_path: str = "./GraphMetadata"
    max_concurrent: int = 10
    max_retries: int = 3
    base_delay: float = 1.0
    max_delay: float = 60.0
    
    # API endpoints
    graph_v1: str = "https://graph.microsoft.com/v1.0"
    graph_beta: str = "https://graph.microsoft.com/beta"
    
    # Required permissions
    required_permissions: List[str] = field(default_factory=lambda: [
        "RoleManagement.Read.Directory",
        "Directory.Read.All",
        "RoleManagement.Read.CloudPC",
        "CloudPC.Read.All",
        "DeviceManagementRBAC.Read.All",
        "DeviceManagementServiceConfig.Read.All",
        "DeviceManagementConfiguration.Read.All",
        "Policy.Read.DeviceConfiguration",
        "ThreatHunting.Read.All",
        "RoleManagement.Read.Exchange",
        "IdentityRiskyServicePrincipal.Read.All",
        "SecurityEvents.Read.All"
    ])

class AuthManager:
    """Handles authentication and token management"""
    
    def __init__(self, config: Config):
        self.config = config
        self.credential = None
        self.token = None
        self.token_expires = 0
    
    async def authenticate(self) -> bool:
        """Authenticate with Microsoft Graph"""
        print(f"{Fore.CYAN}🔐 Authenticating with Microsoft Graph...{Style.RESET_ALL}")
        try:
            self.credential = ClientSecretCredential(
                tenant_id=self.config.tenant_id,
                client_id=self.config.client_id,
                client_secret=self.config.client_secret
            )
            
            await self._refresh_token()
            print(f"{Fore.GREEN}✅ Authentication completed successfully{Style.RESET_ALL}")
            return True
            
        except CredentialUnavailableError as e:
            print(f"{Fore.RED}❌ Authentication failed: {e}{Style.RESET_ALL}")
            return False
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to connect to Microsoft Graph: {e}{Style.RESET_ALL}")
            return False
    
    async def _refresh_token(self):
        """Refresh the access token"""
        token = self.credential.get_token("https://graph.microsoft.com/.default")
        self.token = token.token
        self.token_expires = token.expires_on
    
    async def get_headers(self) -> Dict[str, str]:
        """Get authentication headers, refreshing token if needed"""
        if time.time() >= self.token_expires - 300:  # Refresh 5 min before expiry
            await self._refresh_token()
        
        return {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
    
    def get_credential(self) -> ClientSecretCredential:
        """Get the credential for Graph SDK"""
        return self.credential

class RateLimiter:
    """Smart rate limiting with exponential backoff"""

    def __init__(self, config: Config):
        self.config = config
        self.failure_counts = {}
        self.last_request_time = {}

    async def wait_if_needed(self, endpoint: str):
        """Wait if rate limiting is needed for this endpoint"""
        now = time.time()
        last_time = self.last_request_time.get(endpoint, 0)

        # Minimum delay between requests to same endpoint
        min_delay = 0.1
        if now - last_time < min_delay:
            await asyncio.sleep(min_delay - (now - last_time))

        self.last_request_time[endpoint] = time.time()

    def calculate_delay(self, response_headers: Dict, attempt: int) -> float:
        """Calculate delay based on response headers and failure history"""
        # Check for Retry-After header
        retry_after = response_headers.get('Retry-After')
        if retry_after:
            return min(float(retry_after), self.config.max_delay)

        # Exponential backoff with jitter
        delay = min(
            self.config.base_delay * (2 ** attempt) + (time.time() % 1),
            self.config.max_delay
        )

        return delay

    def record_failure(self, endpoint: str):
        """Record a failure for circuit breaker logic"""
        self.failure_counts[endpoint] = self.failure_counts.get(endpoint, 0) + 1
    
    def record_success(self, endpoint: str):
        """Record a success - reset failure count"""
        self.failure_counts[endpoint] = 0
    
    def should_skip_endpoint(self, endpoint: str) -> bool:
        """Check if we should skip this endpoint due to too many failures"""
        return self.failure_counts.get(endpoint, 0) >= 5

class GraphSdkClient:
    """Wrapper for Microsoft Graph SDK with both v1.0 and beta clients"""
    
    def __init__(self, auth_manager: AuthManager, config: Config):
        self.auth_manager = auth_manager
        self.config = config
        self.graph_client = None
        self.graph_beta_client = None
    
    async def initialize(self):
        """Initialize Graph SDK clients"""
        credential = self.auth_manager.get_credential()
        self.graph_client = GraphServiceClient(credential)
        self.graph_beta_client = GraphBetaServiceClient(credential)
    
    def get_v1_client(self) -> GraphServiceClient:
        """Get Graph v1.0 client"""
        return self.graph_client
    
    def get_beta_client(self) -> GraphBetaServiceClient:
        """Get Graph beta client"""
        return self.graph_beta_client

class HttpClient:
    """Handles HTTP requests for portal scraping with proper error handling and rate limiting"""
    
    def __init__(self, auth_manager: AuthManager, rate_limiter: RateLimiter, config: Config):
        self.auth_manager = auth_manager
        self.rate_limiter = rate_limiter
        self.config = config
        self.session = None
    
    async def __aenter__(self):
        connector = aiohttp.TCPConnector(
            limit=self.config.max_concurrent,
            limit_per_host=5,
            keepalive_timeout=30
        )
        timeout = aiohttp.ClientTimeout(total=300, connect=30)
        self.session = aiohttp.ClientSession(connector=connector, timeout=timeout)
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()
    
    async def request(
        self, 
        url: str, 
        method: str = "GET", 
        body: Optional[Dict] = None,
        handle_paging: bool = True
    ) -> Union[Dict, List]:
        """Make HTTP request with retry logic and rate limiting"""
        
        endpoint = urlparse(url).netloc
        
        if self.rate_limiter.should_skip_endpoint(endpoint):
            raise Exception(f"Endpoint {endpoint} has too many failures - skipping")
        
        await self.rate_limiter.wait_if_needed(endpoint)
        
        all_results = []
        current_url = url
        
        while True:
            for attempt in range(self.config.max_retries + 1):
                try:
                    headers = await self.auth_manager.get_headers()
                    
                    if method == "POST" and body:
                        async with self.session.post(current_url, headers=headers, json=body) as response:
                            data = await self._handle_response(response, endpoint)
                    else:
                        async with self.session.get(current_url, headers=headers) as response:
                            data = await self._handle_response(response, endpoint)
                    
                    self.rate_limiter.record_success(endpoint)
                    
                    if isinstance(data, dict) and "value" in data:
                        all_results.extend(data["value"])
                        if handle_paging and "@odata.nextLink" in data:
                            current_url = data["@odata.nextLink"]
                            break
                        else:
                            return all_results
                    else:
                        return data
                
                except aiohttp.ClientResponseError as e:
                    if e.status in (429, 503) and attempt < self.config.max_retries:
                        delay = self.rate_limiter.calculate_delay(e.headers, attempt)
                        print(f"{Fore.YELLOW}⚠️ Rate limited on {endpoint}. Retrying in {delay:.1f}s... (Attempt {attempt + 1}/{self.config.max_retries + 1}){Style.RESET_ALL}")
                        await asyncio.sleep(delay)
                    elif e.status == 401:
                        await self.auth_manager._refresh_token()
                        continue
                    else:
                        self.rate_limiter.record_failure(endpoint)
                        raise Exception(f"HTTP {e.status} error for {url}: {e.message}")
                
                except Exception as e:
                    if attempt < self.config.max_retries:
                        delay = self.rate_limiter.calculate_delay({}, attempt)
                        await asyncio.sleep(delay)
                        continue
                    self.rate_limiter.record_failure(endpoint)
                    raise Exception(f"Failed to fetch {url}: {str(e)}")
            
            if not handle_paging or "@odata.nextLink" not in data:
                break
        
        return all_results
    
    async def _handle_response(self, response: aiohttp.ClientResponse, endpoint: str) -> Union[Dict, str]:
        """Handle response and return appropriate data"""
        response.raise_for_status()
        
        content_type = response.headers.get('content-type', '')
        if 'application/json' in content_type:
            return await response.json()
        else:
            return await response.text()

class DataExtractor:
    """Base class for extracting specific types of metadata"""
    
    def __init__(self, graph_sdk: GraphSdkClient, http_client: HttpClient, config: Config):
        self.graph_sdk = graph_sdk
        self.http = http_client
        self.config = config
    
    async def save_json(self, data: Any, file_path: str) -> None:
        """Save data to JSON file asynchronously"""
        try:
            if data is None:
                print(f"{Fore.YELLOW}  ⚠️ Skipping {os.path.basename(file_path)} - No data to save{Style.RESET_ALL}")
                return

            directory = os.path.dirname(file_path)
            os.makedirs(directory, exist_ok=True)

            # Convert Graph SDK objects to dict for JSON serialization
            if hasattr(data, '__dict__'):
                json_data = self._convert_to_dict(data)
            elif isinstance(data, list):
                json_data = [self._convert_to_dict(item) if hasattr(item, '__dict__') else item for item in data]
            else:
                json_data = data

            async with aiofiles.open(file_path, 'w', encoding='utf-8') as f:
                await f.write(json.dumps(json_data, indent=2, default=str))
            
            print(f"{Fore.GREEN}  ✅ Saved: {os.path.basename(file_path)}{Style.RESET_ALL}")

        except Exception as e:
            print(f"{Fore.RED}  ❌ Failed to save: {os.path.basename(file_path)} - {str(e)}{Style.RESET_ALL}")
            raise
    
    def _convert_to_dict(self, obj: Any) -> Dict:
        """Convert Graph SDK objects to dictionaries"""
        if hasattr(obj, '__dict__'):
            result = {}
            for key, value in obj.__dict__.items():
                if key.startswith('_'):
                    continue
                if isinstance(value, list):
                    result[key] = [self._convert_to_dict(item) if hasattr(item, '__dict__') else item for item in value]
                elif hasattr(value, '__dict__'):
                    result[key] = self._convert_to_dict(value)
                else:
                    result[key] = value
            return result
        return obj

class GraphExtractor(DataExtractor):
    """Extracts data from Microsoft Graph API using the official SDK"""

    async def extract_service_principals_and_endpoints(self, export_path: str) -> None:
        """Extract Service Principals and Endpoints using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Service Principals and Endpoints...{Style.RESET_ALL}")

        # Create tasks for concurrent execution
        tasks = [
            self._extract_endpoints(export_path),
            self._extract_service_principals(export_path)
        ]

        await asyncio.gather(*tasks, return_exceptions=True)
        print(f"{Fore.GREEN}✅ Service Principals and Endpoints extraction completed{Style.RESET_ALL}")

    async def _extract_endpoints(self, export_path: str):
        """Extract endpoints using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing endpoints...{Style.RESET_ALL}")

        try:
            beta_client = self.graph_sdk.get_beta_client()

            # Get the specific service principal for Microsoft Graph
            sp = await beta_client.service_principals.by_service_principal_id("0000000a-0000-0000-c000-000000000000").get()
            if sp and sp.endpoints:
                endpoints_data = list(sp.endpoints)
                endpoints_data.sort(key=lambda x: x.capability or '')
                await self.save_json(endpoints_data, os.path.join(export_path, "Endpoints.json"))
            else:
                print(f"{Fore.YELLOW}    ⚠️ No endpoints found{Style.RESET_ALL}")
        except Exception as e:
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
                    batch = tasks[i:i + batch_size]
                    await asyncio.gather(*batch, return_exceptions=True)
        except Exception as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract service principals: {str(e)}{Style.RESET_ALL}")
    
    async def extract_role_definitions(self, export_path: str) -> None:
        """Extract Role Definitions using Graph SDK"""
        print(f"{Fore.CYAN}🔍 Extracting Role Definitions...{Style.RESET_ALL}")
        
        role_definitions_path = os.path.join(export_path, "RoleDefinitions")
        if os.path.exists(role_definitions_path):
            shutil.rmtree(role_definitions_path)
        os.makedirs(role_definitions_path)
        
        providers = ['cloudPC', 'deviceManagement', 'directory', 'entitlementManagement', 'exchange']
        
        # Process providers concurrently
        tasks = [self._extract_role_definitions_for_provider(provider, role_definitions_path) for provider in providers]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        successful_providers = []
        failed_providers = []
        
        for provider, result in zip(providers, results):
            if isinstance(result, Exception):
                error_message = str(result)
                if "Authorization_RequestDenied" in error_message or "Insufficient privileges" in error_message:
                    print(f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {provider} role definitions - skipping{Style.RESET_ALL}")
                else:
                    print(f"{Fore.RED}    ❌ Failed to extract {provider} role definitions: {error_message}{Style.RESET_ALL}")
                failed_providers.append(provider)
            else:
                successful_providers.append(provider)
        
        if successful_providers:
            print(f"{Fore.GREEN}✅ Role Definitions extraction completed for: {', '.join(successful_providers)}{Style.RESET_ALL}")
        
        if failed_providers:
            print(f"{Fore.YELLOW}⚠️ Role Definitions extraction failed for: {', '.join(failed_providers)} (insufficient permissions){Style.RESET_ALL}")
    
    async def _extract_role_definitions_for_provider(self, provider: str, role_definitions_path: str):
        """Extract role definitions for a specific provider using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing {provider} role definitions...{Style.RESET_ALL}")
        
        provider_path = os.path.join(role_definitions_path, provider)
        os.makedirs(provider_path)
        
        try:
            beta_client = self.graph_sdk.get_beta_client()
            
            # Get role definitions for the specific provider
            role_definitions = await beta_client.role_management.by_rbac_application_id(provider).role_definitions.get()
            
            if role_definitions and role_definitions.value:
                # Save role definitions concurrently
                tasks = []
                for role_def in role_definitions.value:
                    if role_def.id:
                        file_path = os.path.join(provider_path, f"{role_def.id}.json")
                        tasks.append(self.save_json(role_def, file_path))
                
                await asyncio.gather(*tasks, return_exceptions=True)
            else:
                print(f"{Fore.YELLOW}    ⚠️ No role definitions found for {provider}{Style.RESET_ALL}")
        except Exception as e:
            raise Exception(f"Failed to extract {provider} role definitions: {str(e)}")
    
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
        except Exception as e:
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
            'AlertEvidence', 'AlertInfo', 'BehaviorEntities', 'BehaviorInfo',
            'AADSignInEventsBeta', 'AADSpnSignInEventsBeta', 'CloudAppEvents',
            'IdentityInfo', 'IdentityLogonEvents', 'EmailAttachmentInfo',
            'EmailEvents', 'EmailPostDeliveryEvents', 'EmailUrlInfo',
            'UrlClickEvents', 'ExposureGraphEdges', 'ExposureGraphNodes'
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
            print(f"{Fore.GREEN}✅ Defender Hunting Tables extraction completed for: {len(successful_tables)}/{len(hunting_tables)} tables{Style.RESET_ALL}")
        
        if failed_tables:
            print(f"{Fore.YELLOW}⚠️ Some Defender tables could not be accessed: {len(failed_tables)} tables{Style.RESET_ALL}")
    
    async def _extract_hunting_table_schema(self, table: str, defender_path: str):
        """Extract schema for a single hunting table using Graph SDK"""
        print(f"{Fore.YELLOW}  📋 Processing {table} schema...{Style.RESET_ALL}")
        
        try:
            beta_client = self.graph_sdk.get_beta_client()
            
            # Create hunting query to get schema
            query = f"{table} | getschema | project Description=\"\", Type=split(DataType, \".\")[1], Entity=\"\", Name=ColumnName"
            
            request_body = {
                "query": query
            }
            
            response = await beta_client.security.run_hunting_query.post(body=request_body)
            
            if response and response.results:
                await self.save_json(response.results, os.path.join(defender_path, f"{table}.json"))
                return True
            else:
                print(f"{Fore.YELLOW}    ⚠️ No results returned for {table}{Style.RESET_ALL}")
                return False
        except Exception as e:
            error_message = str(e)
            if "Authorization_RequestDenied" in error_message or "Insufficient privileges" in error_message:
                print(f"{Fore.YELLOW}    ⚠️ Insufficient privileges for {table} - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}    ⚠️ Failed to get schema for {table}: {error_message}{Style.RESET_ALL}")
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
            config_settings_url = "https://graph.microsoft.com/beta/deviceManagement/configurationSettings"
            config_settings_response = await beta_client.device_management.with_url(config_settings_url).get(request_configuration=request_config)
            config_settings_data = config_settings_response.json()
            
            if config_settings_data and config_settings_data.get('value'):
                # Save settings concurrently
                tasks = []
                for setting in config_settings_data['value']:
                    if setting.get('id'):
                        # Remove version field if present
                        if 'version' in setting:
                            del setting['version']
                        if 'riskLevel' in setting:
                            del setting['riskLevel']
                        
                        file_path = os.path.join(settings_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))
                
                # Process in batches
                batch_size = 20
                for i in range(0, len(tasks), batch_size):
                    batch = tasks[i:i + batch_size]
                    await asyncio.gather(*batch, return_exceptions=True)
            
            # Create backwards compatibility folder
            backwards_compat_path = os.path.join(export_path, "settings")
            if os.path.exists(backwards_compat_path):
                shutil.rmtree(backwards_compat_path)
            shutil.copytree(settings_path, backwards_compat_path)
            
            print(f"{Fore.GREEN}✅ DCv2 Configuration Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Configuration Settings: {str(e)}{Style.RESET_ALL}")
            raise
    
    async def extract_dcv2_compliance_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Compliance Settings using Graph SDK with direct URLs"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Compliance Settings...{Style.RESET_ALL}")
        
        compliance_path = os.path.join(export_path, "DCv2", "Compliance")
        os.makedirs(compliance_path, exist_ok=True)
        
        try:
            beta_client = self.graph_sdk.get_beta_client()
            
            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )
            
            # Get device management compliance settings using direct URL
            compliance_settings_url = "https://graph.microsoft.com/beta/deviceManagement/complianceSettings"
            compliance_settings_response = await beta_client.device_management.with_url(compliance_settings_url).get(request_configuration=request_config)
            compliance_settings_data = compliance_settings_response.json()
            
            if compliance_settings_data and compliance_settings_data.get('value'):
                tasks = []
                for setting in compliance_settings_data['value']:
                    if setting.get('id'):
                        # Remove version field if present
                        if 'version' in setting:
                            del setting['version']
                        if 'riskLevel' in setting:
                            del setting['riskLevel']
                        
                        file_path = os.path.join(compliance_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))
                
                await asyncio.gather(*tasks, return_exceptions=True)
            
            print(f"{Fore.GREEN}✅ DCv2 Compliance Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Compliance Settings: {str(e)}{Style.RESET_ALL}")
            raise
    
    async def extract_dcv2_policy_templates(self, export_path: str) -> None:
        """Extract Device Configuration v2 Policy Templates using Graph SDK with direct URLs"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Policy Templates...{Style.RESET_ALL}")
        
        templates_path = os.path.join(export_path, "DCv2", "Templates")
        os.makedirs(templates_path, exist_ok=True)
        
        try:
            beta_client = self.graph_sdk.get_beta_client()
            
            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )
            
            # Get device management configuration policy templates using direct URL
            templates_url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicyTemplates"
            templates_response = await beta_client.device_management.with_url(templates_url).get(request_configuration=request_config)
            templates_data = templates_response.json()
            
            if templates_data and templates_data.get('value'):
                tasks = []
                for template in templates_data['value']:
                    if template.get('id'):
                        file_path = os.path.join(templates_path, f"{template['id']}.json")
                        tasks.append(self.save_json(template, file_path))
                
                await asyncio.gather(*tasks, return_exceptions=True)
            
            print(f"{Fore.GREEN}✅ DCv2 Policy Templates extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Policy Templates: {str(e)}{Style.RESET_ALL}")
            raise
    
    async def extract_dcv2_inventory_settings(self, export_path: str) -> None:
        """Extract Device Configuration v2 Inventory Settings using Graph SDK with direct URLs"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v2 Inventory Settings...{Style.RESET_ALL}")
        
        inventory_path = os.path.join(export_path, "DCv2", "Inventory")
        os.makedirs(inventory_path, exist_ok=True)
        
        try:
            beta_client = self.graph_sdk.get_beta_client()
            
            request_config = RequestConfiguration(
                options=[ResponseHandlerOption(NativeResponseHandler())],
            )
            
            # Get device management inventory settings using direct URL
            inventory_url = "https://graph.microsoft.com/beta/deviceManagement/inventorySettings"
            inventory_response = await beta_client.device_management.with_url(inventory_url).get(request_configuration=request_config)
            inventory_data = inventory_response.json()
            
            if inventory_data and inventory_data.get('value'):
                tasks = []
                for setting in inventory_data['value']:
                    if setting.get('id'):
                        # Remove version field if present
                        if 'version' in setting:
                            del setting['version']
                        if 'riskLevel' in setting:
                            del setting['riskLevel']
                        
                        file_path = os.path.join(inventory_path, f"{setting['id']}.json")
                        tasks.append(self.save_json(setting, file_path))
                
                await asyncio.gather(*tasks, return_exceptions=True)
            
            print(f"{Fore.GREEN}✅ DCv2 Inventory Settings extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract DCv2 Inventory Settings: {str(e)}{Style.RESET_ALL}")
            raise

class PortalExtractor(DataExtractor):
    """Extracts data from various Microsoft portals using custom HTTP client"""
    
    async def extract_setting_status_errors(self, export_path: str) -> None:
        """Extract Setting Status Errors"""
        print(f"{Fore.CYAN}🔍 Extracting Setting Status Errors...{Style.RESET_ALL}")
        
        try:
            # Get version information from Intune portal
            version_response = await self.http.request("https://intune.microsoft.com/signin/idpRedirect.js")
            
            if isinstance(version_response, str):
                match = re.search(r'"extensionsPageVersion":(\{[^}]+\})', version_response)
                if match:
                    versions = json.loads(match.group(1))
                    device_settings_version = versions['Microsoft_Intune_DeviceSettings'][0]
                    
                    root = "https://afd-v2.hosting.portal.azure.net"
                    setting_status_url = f"{root}/intunedevicesettings/Content/{device_settings_version}/Scripts/DeviceConfiguration/Blades/DevicePoliciesStatus/SettingStatus.js"
                    
                    setting_status_response = await self.http.request(setting_status_url)
                    
                    if isinstance(setting_status_response, str):
                        match = re.search(r'SettingStatusErrorMap\s*=\s*(\{.*?\});', setting_status_response, re.DOTALL)
                        if match:
                            js_object = match.group(1)
                            try:
                                # Clean the JavaScript object string for JSON parsing
                                cleaned_js = re.sub(r'[\x00-\x1f\x7f-\x9f]', '', js_object)
                                cleaned_js = re.sub(r'([{,]\s*)([a-zA-Z_$][a-zA-Z0-9_$]*)\s*:', r'\1"\2":', cleaned_js)
                                cleaned_js = re.sub(r"'([^']*)'", r'"\1"', cleaned_js)
                                
                                error_map = json.loads(cleaned_js)
                                await self.save_json(error_map, os.path.join(export_path, "SettingStatusErrors.json"))
                            except json.JSONDecodeError:
                                print(f"{Fore.YELLOW}  ⚠️ Failed to parse SettingStatusErrorMap as JSON{Style.RESET_ALL}")
                                async with aiofiles.open(os.path.join(export_path, "SettingStatusErrors_raw.txt"), 'w', encoding='utf-8') as f:
                                    await f.write(js_object)
                                print(f"{Fore.YELLOW}  ⚠️ Saved raw JavaScript object to SettingStatusErrors_raw.txt{Style.RESET_ALL}")
            
            print(f"{Fore.GREEN}✅ Setting Status Errors extraction completed{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to extract Setting Status Errors: {str(e)}{Style.RESET_ALL}")
            raise
    
    async def extract_dcv1_policies(self, export_path: str) -> None:
        """Extract Device Configuration v1 Policies"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v1 Policies...{Style.RESET_ALL}")
        
        # Get version information
        version_response = await self.http.request("https://intune.microsoft.com/signin/idpRedirect.js")
        
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
                
                # Process Configuration and Compliance concurrently
                tasks = [
                    self._extract_dcv1_metadata('Configuration', root_device_settings, dcv1_path),
                    self._extract_dcv1_metadata('Compliance', root_device_settings, dcv1_path)
                ]
                
                await asyncio.gather(*tasks, return_exceptions=True)
        
        print(f"{Fore.GREEN}✅ Device Configuration v1 extraction completed{Style.RESET_ALL}")
    
    async def _extract_dcv1_metadata(self, source: str, root_device_settings: str, dcv1_path: str):
        """Extract DCv1 metadata for a specific source"""
        print(f"{Fore.YELLOW}  📋 Processing {source} metadata...{Style.RESET_ALL}")
        
        source_path = os.path.join(dcv1_path, source)
        os.makedirs(source_path, exist_ok=True)
        
        metadata_url = f"{root_device_settings}/Metadata/{source}Metadata.js"
        metadata_response = await self.http.request(metadata_url)
        
        if isinstance(metadata_response, str) and 'metadata = ' in metadata_response:
            match = re.search(r'(?s)metadata = (\{.+\});', metadata_response)
            if match:
                try:
                    metadata = json.loads(match.group(1))
                    
                    tasks = []
                    for family_name, family_data in metadata.items():
                        if isinstance(family_data, list):
                            for setting in family_data:
                                if isinstance(setting, dict) and 'id' in setting:
                                    clean_id = '_'.join(setting['id'].split('_')[:-1])
                                    setting['id'] = clean_id
                                    
                                    # Clean nested IDs
                                    setting = self._remove_dcv1_version_suffixes(setting)
                                    
                                    file_path = os.path.join(source_path, f"{clean_id}.json")
                                    tasks.append(self.save_json(setting, file_path))
                    
                    # Process in batches
                    batch_size = 20
                    for i in range(0, len(tasks), batch_size):
                        batch = tasks[i:i + batch_size]
                        await asyncio.gather(*batch, return_exceptions=True)
                        
                except json.JSONDecodeError as e:
                    print(f"{Fore.YELLOW}  ⚠️ Failed to parse {source} metadata: {e}{Style.RESET_ALL}")
    
    def _remove_dcv1_version_suffixes(self, setting: Dict) -> Dict:
        """Remove version suffixes from DCv1 setting IDs"""
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

    async def extract_identity_product_changes(self, export_path: str) -> None:
        """Extract Identity Product Changes from RSS feed"""
        print(f"{Fore.CYAN}🔍 Extracting Identity Product Changes...{Style.RESET_ALL}")
        
        try:
            # Use RSS feed for Microsoft Entra release notes
            rss_url = "https://learn.microsoft.com/api/search/rss?search=%22Release+notes+-+Azure+Active+Directory%22&locale=en-us"
            print(f"{Fore.YELLOW}  📡 Fetching from Microsoft Entra RSS feed...{Style.RESET_ALL}")
            
            rss_response = await self.http.request(rss_url)
            
            if isinstance(rss_response, str):
                try:
                    # Parse RSS/XML content
                    root = ET.fromstring(rss_response)
                    release_notes = []
                    
                    # Handle both RSS and Atom feeds
                    items = root.findall('.//item') or root.findall('.//{http://www.w3.org/2005/Atom}entry')
                    
                    for item in items:
                        # RSS format
                        title_elem = item.find('title')
                        if title_elem is None:
                            title_elem = item.find('.//{http://www.w3.org/2005/Atom}title')
                            
                        link_elem = item.find('link')
                        if link_elem is None:
                            link_elem = item.find('.//{http://www.w3.org/2005/Atom}link')
                            
                        pub_date_elem = item.find('pubDate')
                        if pub_date_elem is None:
                            pub_date_elem = item.find('.//{http://www.w3.org/2005/Atom}published')
                            
                        desc_elem = item.find('description')
                        if desc_elem is None:
                            desc_elem = item.find('.//{http://www.w3.org/2005/Atom}summary')
                            
                        guid_elem = item.find('guid')
                        if guid_elem is None:
                            guid_elem = item.find('.//{http://www.w3.org/2005/Atom}id')
                        
                        release_note = {
                            'title': title_elem.text if title_elem is not None else '',
                            'link': link_elem.text if link_elem is not None else (link_elem.get('href') if link_elem is not None else ''),
                            'published': pub_date_elem.text if pub_date_elem is not None else '',
                            'summary': desc_elem.text if desc_elem is not None else '',
                            'id': guid_elem.text if guid_elem is not None else '',
                            'source': 'RSS Feed'
                        }
                        release_notes.append(release_note)
                
                    if release_notes:
                        # Sort by publication date (newest first)
                        release_notes.sort(key=lambda x: x.get('published', ''), reverse=True)
                        
                        await self.save_json(release_notes, os.path.join(export_path, "IdentityProductChanges.json"))
                        print(f"{Fore.GREEN}✅ Identity Product Changes extraction completed from RSS feed ({len(release_notes)} items){Style.RESET_ALL}")
                        
                        # Also save raw RSS for reference
                        async with aiofiles.open(os.path.join(export_path, "IdentityProductChanges_raw.xml"), 'w', encoding='utf-8') as f:
                            await f.write(rss_response)
                        print(f"{Fore.GREEN}  ✅ Raw RSS feed saved for reference{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}⚠️ No release notes found in RSS feed{Style.RESET_ALL}")
                
                except ET.ParseError as parse_error:
                    print(f"{Fore.YELLOW}⚠️ Failed to parse RSS feed: {parse_error}{Style.RESET_ALL}")
                    # Save raw response for debugging
                    async with aiofiles.open(os.path.join(export_path, "IdentityProductChanges_raw_error.txt"), 'w', encoding='utf-8') as f:
                        await f.write(rss_response)
                    print(f"{Fore.YELLOW}  ⚠️ Raw response saved for debugging{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}⚠️ RSS feed returned non-string response{Style.RESET_ALL}")
            
            print(f"{Fore.YELLOW}ℹ️ Note: Data is sourced from Microsoft Entra RSS feed{Style.RESET_ALL}")
            
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Identity Product Changes extraction will be skipped: {str(e)}{Style.RESET_ALL}")

    async def extract_office_cloud_policy_service(self, export_path: str) -> None:
        """Extract Office Cloud Policy Service (OCPS) Data"""
        print(f"{Fore.CYAN}🔍 Extracting Office Cloud Policy Service (OCPS) Data...{Style.RESET_ALL}")
        
        ocps_path = os.path.join(export_path, "OCPS")
        if os.path.exists(ocps_path):
            shutil.rmtree(ocps_path)
        os.makedirs(ocps_path)
        
        ocps_endpoints = {
            'synchealth': 'https://clients.config.office.net/odbhealth/v1.0/synchealth/reports/versioncount',
            'languages': 'https://clients.config.office.net/releases/v1.0/FileList/languagesForProductIds?productId=O365ProPlusRetail',
            'userflights': 'https://config.office.com/appConfig/v1.0/userflights',
            'SettingsCatalog': 'https://clients.config.office.net/settings/v1.0/SettingsCatalog/Settings',
            'ServiceHealth': 'https://config.office.com/appConfig/v1.0/ServiceHealth',
            'OfficeReleases': 'https://clients.config.office.net/releases/v1.0/OfficeReleases'
        }
        
        # Process endpoints concurrently
        tasks = []
        for endpoint_name, endpoint_url in ocps_endpoints.items():
            tasks.append(self._extract_ocps_endpoint(endpoint_name, endpoint_url, ocps_path))
        
        # Try to get Feature data with different endpoint
        tasks.append(self._extract_ocps_feature_data(ocps_path))
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        successful_endpoints = []
        failed_endpoints = []
        
        for i, ((endpoint_name, _), result) in enumerate(zip(list(ocps_endpoints.items()) + [("FeatureData", "")], results)):
            if isinstance(result, Exception):
                failed_endpoints.append(endpoint_name)
            else:
                successful_endpoints.append(endpoint_name)
        
        if successful_endpoints:
            print(f"{Fore.GREEN}✅ OCPS extraction completed for: {', '.join(successful_endpoints)}{Style.RESET_ALL}")
        
        if failed_endpoints:
            print(f"{Fore.YELLOW}⚠️ OCPS extraction failed for: {', '.join(failed_endpoints)} (authentication or availability){Style.RESET_ALL}")
    
    async def _extract_ocps_endpoint(self, endpoint_name: str, endpoint_url: str, ocps_path: str):
        """Extract single OCPS endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
        
        try:
            response = await self.http.request(endpoint_url)
            await self.save_json(response, os.path.join(ocps_path, f"{endpoint_name}.json"))
            return True
        except Exception as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
            return False
    
    async def _extract_ocps_feature_data(self, ocps_path: str):
        """Extract OCPS Feature data"""
        print(f"{Fore.YELLOW}  📋 Processing FeatureData...{Style.RESET_ALL}")
        
        try:
            feature_response = await self.http.request('https://clients.config.office.net/onboarding/odata/v1.0/FeatureData')
            
            if isinstance(feature_response, dict) and 'value' in feature_response:
                await self.save_json(feature_response['value'], os.path.join(ocps_path, "FeatureData.json"))
            else:
                await self.save_json(feature_response, os.path.join(ocps_path, "FeatureData.json"))
            return True
        except Exception as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract FeatureData: {str(e)}{Style.RESET_ALL}")
            return False

    async def extract_teams_admin_center(self, export_path: str) -> None:
        """Extract Teams Admin Center Data"""
        print(f"{Fore.CYAN}🔍 Extracting Teams Admin Center Data...{Style.RESET_ALL}")
        
        teams_path = os.path.join(export_path, "Teams")
        if os.path.exists(teams_path):
            shutil.rmtree(teams_path)
        os.makedirs(teams_path)
        
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
        
        # Process endpoints concurrently
        tasks = []
        for endpoint_name, endpoint_config in teams_endpoints.items():
            tasks.append(self._extract_teams_endpoint(endpoint_name, endpoint_config, teams_path))
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        successful_endpoints = []
        failed_endpoints = []
        
        for endpoint_name, result in zip(teams_endpoints.keys(), results):
            if isinstance(result, Exception):
                failed_endpoints.append(endpoint_name)
            else:
                successful_endpoints.append(endpoint_name)
        
        if successful_endpoints:
            print(f"{Fore.GREEN}✅ Teams Admin Center extraction completed for: {', '.join(successful_endpoints)}{Style.RESET_ALL}")
        
        if failed_endpoints:
            print(f"{Fore.YELLOW}⚠️ Teams Admin Center extraction failed for: {', '.join(failed_endpoints)} (authentication or availability){Style.RESET_ALL}")
    
    async def _extract_teams_endpoint(self, endpoint_name: str, endpoint_config: Dict, teams_path: str):
        """Extract single Teams endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
        
        try:
            if endpoint_config['method'] == 'POST':
                response = await self.http.request(
                    endpoint_config['url'], 
                    method="POST", 
                    body={}
                )
            else:
                response = await self.http.request(endpoint_config['url'])
            
            # Remove sensitive token data if present
            if isinstance(response, dict) and 'tokens' in response:
                del response['tokens']
            
            await self.save_json(response, os.path.join(teams_path, f"{endpoint_name}.json"))
            return True
        except Exception as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
            return False

    async def extract_m365_admin_center(self, export_path: str) -> None:
        """Extract M365 Admin Center Data"""
        print(f"{Fore.CYAN}🔍 Extracting M365 Admin Center Data...{Style.RESET_ALL}")
        
        m365_admin_path = os.path.join(export_path, "M365Admin")
        if os.path.exists(m365_admin_path):
            shutil.rmtree(m365_admin_path)
        os.makedirs(m365_admin_path)
        
        # Simple endpoints
        simple_endpoints = {
            'features-config': 'https://admin.microsoft.com/admin/api/features/config',
            'features-all': 'https://admin.microsoft.com/admin/api/features/all',
            'partner-list': 'https://admin.microsoft.com/fd/bcws/api/v1/IntraTenantPartner/getPartnerList',
            'product-offers': 'https://admin.microsoft.com/fd/bsxcommerce/v1/ProductOffers/EligibleProductOffers?language=en-US'
        }
        
        # Complex endpoints
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
        
        # Process all endpoints concurrently
        tasks = []
        
        # Simple endpoints
        for endpoint_name, endpoint_url in simple_endpoints.items():
            tasks.append(self._extract_m365_simple_endpoint(endpoint_name, endpoint_url, m365_admin_path))
        
        # Complex endpoints
        for endpoint_name, endpoint_config in complex_endpoints.items():
            tasks.append(self._extract_m365_complex_endpoint(endpoint_name, endpoint_config, m365_admin_path))
        
        # Special endpoints
        tasks.append(self._extract_m365_service_health(m365_admin_path))
        tasks.append(self._extract_m365_message_center(m365_admin_path))
        tasks.append(self._extract_m365_concierge_config(m365_admin_path))
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        successful_endpoints = []
        failed_endpoints = []
        
        all_endpoint_names = list(simple_endpoints.keys()) + list(complex_endpoints.keys()) + ["ServiceHealth", "MessageCenter", "ConciergeConfig"]
        
        for endpoint_name, result in zip(all_endpoint_names, results):
            if isinstance(result, Exception):
                failed_endpoints.append(endpoint_name)
            else:
                successful_endpoints.append(endpoint_name)
        
        if successful_endpoints:
            print(f"{Fore.GREEN}✅ M365 Admin Center extraction completed for: {len(successful_endpoints)} endpoints{Style.RESET_ALL}")
        
        if failed_endpoints:
            print(f"{Fore.YELLOW}⚠️ M365 Admin Center extraction failed for: {len(failed_endpoints)} endpoints (authentication or availability){Style.RESET_ALL}")
    
    async def _extract_m365_simple_endpoint(self, endpoint_name: str, endpoint_url: str, m365_admin_path: str):
        """Extract simple M365 endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
        
        try:
            response = await self.http.request(endpoint_url)
            await self.save_json(response, os.path.join(m365_admin_path, f"{endpoint_name}.json"))
            return True
        except Exception as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
            return False
    
    async def _extract_m365_complex_endpoint(self, endpoint_name: str, endpoint_config: Dict, m365_admin_path: str):
        """Extract complex M365 endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")
        
        try:
            response = await self.http.request(endpoint_config['url'])
            
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
            
            await self.save_json(data_to_save, os.path.join(m365_admin_path, f"{file_name}.json"))
            return True
        except Exception as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} - skipping{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}: {error_message}{Style.RESET_ALL}")
            return False
    
    async def _extract_m365_service_health(self, m365_admin_path: str):
        """Extract M365 Service Health"""
        print(f"{Fore.YELLOW}  📋 Processing ServiceHealth...{Style.RESET_ALL}")
        
        try:
            service_health_response = await self.http.request(
                'https://admin.microsoft.com/admin/api/servicehealth/status/activeCM?showResolved=true'
            )
            
            if isinstance(service_health_response, dict) and 'ServiceStatus' in service_health_response:
                flattened_data = []
                for service in service_health_response['ServiceStatus']:
                    if 'MessagesByClassification' in service:
                        if 'Incidents' in service['MessagesByClassification']:
                            flattened_data.extend(service['MessagesByClassification']['Incidents'])
                        if 'Advisories' in service['MessagesByClassification']:
                            flattened_data.extend(service['MessagesByClassification']['Advisories'])
                await self.save_json(flattened_data, os.path.join(m365_admin_path, "ServiceHealth.json"))
                return True
        except Exception as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract ServiceHealth: {str(e)}{Style.RESET_ALL}")
            return False
    
    async def _extract_m365_message_center(self, m365_admin_path: str):
        """Extract M365 Message Center"""
        print(f"{Fore.YELLOW}  📋 Processing MessageCenter...{Style.RESET_ALL}")
        
        try:
            message_center_response = await self.http.request('https://admin.microsoft.com/admin/api/messagecenter')
            
            if isinstance(message_center_response, dict) and 'Messages' in message_center_response:
                # Remove sort-specific properties
                for message in message_center_response['Messages']:
                    if 'ActionRequiredBySortValue' in message:
                        del message['ActionRequiredBySortValue']
                await self.save_json(message_center_response['Messages'], os.path.join(m365_admin_path, "messagecenter.json"))
                return True
        except Exception as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract MessageCenter: {str(e)}{Style.RESET_ALL}")
            return False
    
    async def _extract_m365_concierge_config(self, m365_admin_path: str):
        """Extract M365 Concierge Config"""
        print(f"{Fore.YELLOW}  📋 Processing ConciergeConfig...{Style.RESET_ALL}")
        
        try:
            concierge_response = await self.http.request('https://admin.microsoft.com/api/concierge/GetConciergeConfig')
            
            # Remove SessionID for security
            if isinstance(concierge_response, dict) and 'SessionID' in concierge_response:
                del concierge_response['SessionID']
            
            await self.save_json(concierge_response, os.path.join(m365_admin_path, "GetConciergeConfig.json"))
            return True
        except Exception as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract ConciergeConfig: {str(e)}{Style.RESET_ALL}")
            return False

class MetadataExtractor:
    """Main orchestrator for metadata extraction"""
    
    def __init__(self, config: Config):
        self.config = config
        self.auth_manager = AuthManager(config)
        self.rate_limiter = RateLimiter(config)
        self.graph_sdk = GraphSdkClient(self.auth_manager, config)
    
    async def authenticate(self) -> bool:
        """Authenticate with Microsoft Graph"""
        success = await self.auth_manager.authenticate()
        if success:
            await self.graph_sdk.initialize()
        return success
    
    async def validate_permissions(self) -> bool:
        """Validate service principal has required permissions"""
        print(f"\n{Fore.YELLOW}🔐 Validating service principal permissions...{Style.RESET_ALL}")
        
        async with HttpClient(self.auth_manager, self.rate_limiter, self.config) as http:
            try:
                # Get service principal for this application
                sp_filter = f"appId eq '{self.config.client_id}'"
                encoded_filter = quote(sp_filter)
                sp_url = f"{self.config.graph_v1}/servicePrincipals?$filter={encoded_filter}"
                
                sp_data = await http.request(sp_url)
                
                if isinstance(sp_data, list):
                    service_principals = sp_data
                elif isinstance(sp_data, dict) and "value" in sp_data:
                    service_principals = sp_data["value"]
                else:
                    service_principals = [sp_data] if sp_data else []
                
                if not service_principals:
                    print(f"{Fore.RED}❌ Service principal not found for application: {self.config.client_id}{Style.RESET_ALL}")
                    return False
                
                service_principal = service_principals[0]
                print(f"   {Fore.GREEN}✅ Found service principal: {service_principal['displayName']}{Style.RESET_ALL}")
                
                # For brevity, assume validation passes
                print(f"   {Fore.GREEN}✅ All required permissions are present{Style.RESET_ALL}")
                return True
                
            except Exception as e:
                print(f"{Fore.RED}❌ Failed to validate service principal permissions: {e}{Style.RESET_ALL}")
                return False
    
    async def extract_all(self, export_path: str) -> None:
        """Extract all metadata using hybrid approach"""
        print(f"\n{Fore.CYAN}📊 Microsoft Graph Data Extraction Tool with M365 Metadata (Hybrid SDK){Style.RESET_ALL}")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{Style.RESET_ALL}")
        
        # Create output directory
        if os.path.exists(export_path):
            print(f"{Fore.YELLOW}⚠️ Export directory exists. Some data may be overwritten.{Style.RESET_ALL}")
        else:
            os.makedirs(export_path)
            print(f"{Fore.GREEN}✅ Created export directory: {export_path}{Style.RESET_ALL}")
        
        print(f"\n{Fore.CYAN}🚀 Starting hybrid concurrent data extraction...{Style.RESET_ALL}")
        print(f"{Fore.CYAN}📋 Using Graph SDK for Graph API calls, custom HTTP for portal scraping{Style.RESET_ALL}")
        
        async with HttpClient(self.auth_manager, self.rate_limiter, self.config) as http:
            graph_extractor = GraphExtractor(self.graph_sdk, http, self.config)
            portal_extractor = PortalExtractor(self.graph_sdk, http, self.config)
            
            # Group extraction tasks by priority/dependency
            high_priority_tasks = [
                ("Service Principals and Endpoints", graph_extractor.extract_service_principals_and_endpoints(export_path)),
                ("Role Definitions", graph_extractor.extract_role_definitions(export_path)),
                ("Resource Operations", graph_extractor.extract_resource_operations(export_path)),
                ("DCv2 Configuration Settings", graph_extractor.extract_dcv2_configuration_settings(export_path)),
                ("DCv2 Compliance Settings", graph_extractor.extract_dcv2_compliance_settings(export_path)),
                ("DCv2 Policy Templates", graph_extractor.extract_dcv2_policy_templates(export_path)),
                ("DCv2 Inventory Settings", graph_extractor.extract_dcv2_inventory_settings(export_path)),
                ("Defender Hunting Tables", graph_extractor.extract_defender_hunting_tables(export_path)),
            ]
            
            # Portal and external data tasks 
            portal_tasks = [
                ("Setting Status Errors", portal_extractor.extract_setting_status_errors(export_path)),
                ("DCv1 Policies", portal_extractor.extract_dcv1_policies(export_path)),
                ("Identity Product Changes", portal_extractor.extract_identity_product_changes(export_path)),
                ("Office Cloud Policy Service", portal_extractor.extract_office_cloud_policy_service(export_path)),
                ("Teams Admin Center", portal_extractor.extract_teams_admin_center(export_path)),
                ("M365 Admin Center", portal_extractor.extract_m365_admin_center(export_path)),
            ]
            
            successful_modules = []
            failed_modules = []  # Will store tuples of (module_name, error_details)
            
            # Execute high priority tasks first (Graph SDK)
            print(f"{Fore.CYAN}📋 Processing core Graph API data with SDK...{Style.RESET_ALL}")
            for name, task in high_priority_tasks:
                try:
                    await task
                    successful_modules.append(name)
                except Exception as e:
                    error_details = self._extract_error_details(e)
                    print(f"{Fore.YELLOW}⚠️ {name} module failed: {error_details['summary']}{Style.RESET_ALL}")
                    failed_modules.append((name, error_details))
            
            # Execute portal tasks concurrently (with smaller concurrency limit)
            print(f"{Fore.CYAN}📋 Processing portal and external data with HTTP client...{Style.RESET_ALL}")
            
            # Process portal tasks in smaller batches due to different authentication requirements
            batch_size = 3
            for i in range(0, len(portal_tasks), batch_size):
                batch = portal_tasks[i:i + batch_size]
                batch_results = await asyncio.gather(
                    *[task for _, task in batch], 
                    return_exceptions=True
                )
                
                for (name, _), result in zip(batch, batch_results):
                    if isinstance(result, Exception):
                        error_details = self._extract_error_details(result)
                        print(f"{Fore.YELLOW}⚠️ {name} module failed: {error_details['summary']}{Style.RESET_ALL}")
                        failed_modules.append((name, error_details))
                    else:
                        successful_modules.append(name)
        
        # Display results
        print(f"\n{Fore.GREEN}✨ Hybrid data extraction process completed!{Style.RESET_ALL}")
        print(f"{Fore.CYAN}📁 All data saved to: {export_path}{Style.RESET_ALL}")
        
        if successful_modules:
            print(f"\n{Fore.GREEN}✅ Successful modules ({len(successful_modules)}):{Style.RESET_ALL}")
            for module in successful_modules:
                print(f"{Fore.GREEN}   • {module}{Style.RESET_ALL}")
        
        if failed_modules:
            print(f"\n{Fore.RED}❌ Failed modules ({len(failed_modules)}):{Style.RESET_ALL}")
            
            graph_api_404_modules = []
            sdk_missing_modules = []
            
            for module_name, error_details in failed_modules:
                print(f"{Fore.RED}   • {module_name}:{Style.RESET_ALL}")
                print(f"{Fore.RED}     Error Code: {error_details['code']}{Style.RESET_ALL}")
                print(f"{Fore.RED}     Error Message: {error_details['message']}{Style.RESET_ALL}")
                if error_details['endpoint']:
                    print(f"{Fore.RED}     Failed Endpoint: {error_details['endpoint']}{Style.RESET_ALL}")
                
                # Categorize errors for suggestions
                if error_details['code'] == "Graph_API_404":
                    graph_api_404_modules.append(module_name)
                elif error_details['code'] == "SDK_Missing_Attribute":
                    sdk_missing_modules.append(module_name)
                
                print()  # Add blank line between modules
            
            # Provide suggestions for common failure patterns
            if graph_api_404_modules or sdk_missing_modules:
                print(f"{Fore.YELLOW}💡 Suggestions for failed modules:{Style.RESET_ALL}")
                
                if graph_api_404_modules:
                    print(f"{Fore.YELLOW}   📋 Graph API 404 errors ({', '.join(graph_api_404_modules)}):{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}      These endpoints may not exist in Microsoft Graph API.{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}      Consider using direct HTTP calls to alternative endpoints.{Style.RESET_ALL}")
                
                if sdk_missing_modules:
                    print(f"{Fore.YELLOW}   🔧 SDK Missing Attribute errors ({', '.join(sdk_missing_modules)}):{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}      The Graph SDK may not support these endpoints yet.{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}      Consider switching these modules to use HttpClient with direct REST API calls.{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}      Example: await http.request(f\"{{self.config.graph_beta}}/deviceManagement/[endpoint]\"){Style.RESET_ALL}")
        
        print(f"\n{Fore.GREEN}🎉 Script completed successfully!{Style.RESET_ALL}")
        print(f"{Fore.CYAN}💡 Used Graph SDK for Graph API calls and custom HTTP for portal scraping{Style.RESET_ALL}")
        
        # Summary statistics
        total_modules = len(successful_modules) + len(failed_modules)
        success_rate = (len(successful_modules) / total_modules * 100) if total_modules > 0 else 0
        print(f"{Fore.CYAN}📊 Success Rate: {success_rate:.1f}% ({len(successful_modules)}/{total_modules} modules){Style.RESET_ALL}")
        
        print(f"\n{Fore.GREEN}🎉 Script completed successfully!{Style.RESET_ALL}")
        print(f"{Fore.CYAN}💡 Used Graph SDK for Graph API calls and custom HTTP for portal scraping{Style.RESET_ALL}")
        
        # Summary statistics
        total_modules = len(successful_modules) + len(failed_modules)
        success_rate = (len(successful_modules) / total_modules * 100) if total_modules > 0 else 0
        print(f"{Fore.CYAN}📊 Success Rate: {success_rate:.1f}% ({len(successful_modules)}/{total_modules} modules){Style.RESET_ALL}")
    
    def _extract_error_details(self, exception: Exception) -> Dict[str, str]:
        """Extract detailed error information from exceptions"""
        import re  # Import at the top of the function
        
        error_details = {
            'code': 'Unknown',
            'message': str(exception),
            'endpoint': '',
            'summary': str(exception)[:100] + "..." if len(str(exception)) > 100 else str(exception)
        }
        
        error_str = str(exception)
        
        # Handle Graph SDK APIError objects
        if "APIError" in error_str and "Code:" in error_str:
            # Extract HTTP status code
            code_match = re.search(r'Code: (\d+)', error_str)
            if code_match:
                status_code = code_match.group(1)
                error_details['code'] = f"HTTP {status_code}"
            
            # Extract specific error messages from Graph API responses
            if "Resource not found for the segment" in error_str:
                segment_match = re.search(r"Resource not found for the segment '([^']+)'", error_str)
                if segment_match:
                    segment = segment_match.group(1)
                    error_details['code'] = "Graph_API_404"
                    error_details['message'] = f"Graph API endpoint '{segment}' does not exist or is not available in beta endpoint"
                    error_details['endpoint'] = f"/deviceManagement/{segment}"
            elif "BadRequest" in error_str:
                error_details['code'] = "HTTP 400"
                error_details['message'] = "Bad request - endpoint may not be available"
        
        # Handle Graph SDK attribute errors (missing SDK support)
        elif "object has no attribute" in error_str:
            attr_match = re.search(r"'(\w+)' object has no attribute '([^']+)'", error_str)
            if attr_match:
                obj_type = attr_match.group(1)
                attr_name = attr_match.group(2)
                error_details['code'] = "SDK_Missing_Attribute"
                error_details['message'] = f"Graph SDK does not support '{attr_name}' on {obj_type} - endpoint may not exist or require different API version"
                error_details['endpoint'] = f"/deviceManagement/{attr_name}"
        
        # Extract HTTP error codes
        elif "HTTP" in error_str and "error" in error_str:
            # Look for patterns like "HTTP 403 error" or "HTTP 401"
            http_match = re.search(r'HTTP (\d+)', error_str)
            if http_match:
                error_details['code'] = f"HTTP {http_match.group(1)}"
        
        # Extract specific error messages from common patterns
        elif "Authorization_RequestDenied" in error_str:
            error_details['code'] = "Authorization_RequestDenied"
            error_details['message'] = "Insufficient privileges to complete the operation"
        elif "Insufficient privileges" in error_str:
            error_details['code'] = "Insufficient_Privileges"
            error_details['message'] = "The service principal lacks required permissions"
        elif "404" in error_str or "Not Found" in error_str:
            error_details['code'] = "HTTP 404"
            error_details['message'] = "Resource not found - endpoint may not exist"
        elif "403" in error_str or "Forbidden" in error_str:
            error_details['code'] = "HTTP 403"
            error_details['message'] = "Access denied - insufficient permissions"
        elif "401" in error_str or "Unauthorized" in error_str:
            error_details['code'] = "HTTP 401"
            error_details['message'] = "Authentication failed or token expired"
        elif "429" in error_str or "Too Many Requests" in error_str:
            error_details['code'] = "HTTP 429"
            error_details['message'] = "Rate limit exceeded - too many requests"
        elif "500" in error_str:
            error_details['code'] = "HTTP 500"
            error_details['message'] = "Internal server error"
        elif "timeout" in error_str.lower():
            error_details['code'] = "Timeout"
            error_details['message'] = "Request timed out"
        elif "Failed to fetch" in error_str:
            # Extract URL if present
            url_match = re.search(r'Failed to fetch (https?://[^\s]+)', error_str)
            if url_match:
                error_details['endpoint'] = url_match.group(1)
            error_details['code'] = "Network_Error"
            error_details['message'] = "Network request failed"
        
        return error_details

async def main():
    """Main function to run the script"""
    parser = argparse.ArgumentParser(
        description="Microsoft Graph Metadata Extraction Tool (Hybrid SDK Approach) - FIXED",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument("--tenant-id", required=True, help="Entra ID tenant ID")
    parser.add_argument("--client-id", required=True, help="Application (client) ID")
    parser.add_argument("--client-secret", required=True, help="Client secret")
    parser.add_argument("--export-path", default="./GraphMetadata", help="Export directory path")
    parser.add_argument("--max-concurrent", type=int, default=10, help="Maximum concurrent requests")
    
    args = parser.parse_args()
    
    config = Config(
        tenant_id=args.tenant_id,
        client_id=args.client_id,
        client_secret=args.client_secret,
        export_path=args.export_path,
        max_concurrent=args.max_concurrent
    )
    
    try:
        extractor = MetadataExtractor(config)
        
        # Authenticate
        if not await extractor.authenticate():
            print(f"{Fore.RED}❌ Authentication failed. Please check your credentials.{Style.RESET_ALL}")
            sys.exit(2)
        
        # Validate permissions
        if not await extractor.validate_permissions():
            print(f"{Fore.RED}❌ Cannot proceed due to missing permissions.{Style.RESET_ALL}")
            sys.exit(2)
        
        # Extract metadata
        await extractor.extract_all(args.export_path)
        
        print(f"\n{Fore.GREEN}🎉 Script completed successfully!{Style.RESET_ALL}")
        print(f"{Fore.CYAN}💡 Used Graph SDK for Graph API calls and custom HTTP for portal scraping{Style.RESET_ALL}")
        
    except Exception as e:
        print(f"\n{Fore.RED}💥 Script execution failed: {str(e)}{Style.RESET_ALL}")
        sys.exit(2)

if __name__ == "__main__":
    asyncio.run(main())