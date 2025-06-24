#!/usr/bin/env python3
"""
Terraform Provider Registry Manager

A tool to interact with the Terraform Provider Registry API to list, match, and remove provider versions.
Note: Version deletion requires provider maintainer privileges and may not be supported by all registries.

Usage:
    python terraform_provider_manager.py --namespace hashicorp --provider aws --version 5.0.0 [--token <api_token>]
"""

import argparse
import asyncio
import json
import sys
import re
from datetime import datetime
from typing import Dict, List, Optional, Any
from dataclasses import dataclass
import aiohttp
from packaging import version

from colorama import Fore, Style, init

# Initialize colorama
init()

@dataclass
class Config:
    """Configuration for the Terraform Provider Registry Manager"""
    namespace: str
    provider: str
    target_version: Optional[str] = None
    api_token: Optional[str] = None
    registry_url: str = "https://registry.terraform.io"
    max_retries: int = 3
    timeout: int = 30

class TerraformProviderRegistryClient:
    """Client for interacting with the Terraform Provider Registry API"""
    
    def __init__(self, config: Config):
        self.config = config
        self.session = None
        self.base_url = f"{config.registry_url}/v1"
    
    async def __aenter__(self):
        """Async context manager entry"""
        timeout = aiohttp.ClientTimeout(total=self.config.timeout)
        self.session = aiohttp.ClientSession(timeout=timeout)
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit"""
        if self.session:
            await self.session.close()
    
    def _get_headers(self) -> Dict[str, str]:
        """Get headers for API requests"""
        headers = {
            'User-Agent': 'Terraform-Provider-Registry-Manager/1.0',
            'Accept': 'application/json'
        }
        
        if self.config.api_token:
            headers['Authorization'] = f'Bearer {self.config.api_token}'
        
        return headers
    
    async def _make_request(
        self, 
        method: str, 
        endpoint: str, 
        **kwargs
    ) -> Dict[str, Any]:
        """Make HTTP request with retry logic"""
        url = f"{self.base_url}{endpoint}"
        headers = self._get_headers()
        
        for attempt in range(self.config.max_retries + 1):
            try:
                async with self.session.request(
                    method,
                    url,
                    headers=headers,
                    **kwargs
                ) as response:
                    if response.status == 429:  # Rate limited
                        retry_after = int(response.headers.get('Retry-After', 5))
                        print(f"{Fore.YELLOW}⚠️ Rate limited. Retrying in {retry_after}s... (Attempt {attempt + 1}/{self.config.max_retries + 1}){Style.RESET_ALL}")
                        await asyncio.sleep(retry_after)
                        continue
                    
                    response.raise_for_status()
                    
                    if response.headers.get('content-type', '').startswith('application/json'):
                        return await response.json()
                    else:
                        return {'text': await response.text()}
            
            except aiohttp.ClientResponseError as e:
                if e.status == 404:
                    raise TerraformRegistryError(f"Provider not found: {self.config.namespace}/{self.config.provider}") from e
                elif e.status == 401:
                    raise TerraformRegistryError("Authentication failed. Check your API token.") from e
                elif e.status == 403:
                    raise TerraformRegistryError("Access denied. Insufficient permissions.") from e
                elif attempt == self.config.max_retries:
                    raise TerraformRegistryError(f"HTTP {e.status}: {e.message}") from e
                else:
                    wait_time = 2 ** attempt
                    print(f"{Fore.YELLOW}⚠️ Request failed (HTTP {e.status}). Retrying in {wait_time}s... (Attempt {attempt + 1}/{self.config.max_retries + 1}){Style.RESET_ALL}")
                    await asyncio.sleep(wait_time)
            
            except (aiohttp.ClientError, asyncio.TimeoutError, ConnectionError) as e:
                if attempt == self.config.max_retries:
                    raise TerraformRegistryError(f"Request failed: {str(e)}") from e
                else:
                    wait_time = 2 ** attempt
                    print(f"{Fore.YELLOW}⚠️ Request failed: {str(e)}. Retrying in {wait_time}s... (Attempt {attempt + 1}/{self.config.max_retries + 1}){Style.RESET_ALL}")
                    await asyncio.sleep(wait_time)
        
        raise TerraformRegistryError("Max retries exceeded")
    
    async def get_provider_versions(self) -> List[Dict[str, Any]]:
        """Get all versions for a provider using the Provider Registry API"""
        print(f"{Fore.CYAN}🔍 Fetching versions for provider {self.config.namespace}/{self.config.provider}...{Style.RESET_ALL}")
        
        endpoint = f"/providers/{self.config.namespace}/{self.config.provider}/versions"
        
        try:
            response = await self._make_request('GET', endpoint)
            
            if 'versions' not in response:
                raise TerraformRegistryError("Invalid response format: missing 'versions' field")
            
            versions = response['versions']
            print(f"{Fore.GREEN}✅ Found {len(versions)} versions for {self.config.namespace}/{self.config.provider}{Style.RESET_ALL}")
            
            return versions
        
        except TerraformRegistryError:
            raise
        except Exception as e:
            raise TerraformRegistryError(f"Failed to fetch provider versions: {str(e)}") from e
    
    async def get_provider_details(self) -> Dict[str, Any]:
        """Get provider details from the extended registry API (not part of standard protocol)"""
        print(f"{Fore.YELLOW}📋 Fetching provider details...{Style.RESET_ALL}")
        
        # This is an extended endpoint that may not be available on all registries
        endpoint = f"/providers/{self.config.namespace}/{self.config.provider}"
        
        try:
            response = await self._make_request('GET', endpoint)
            print(f"{Fore.GREEN}✅ Retrieved provider details{Style.RESET_ALL}")
            return response
        
        except TerraformRegistryError as e:
            if "404" in str(e) or "not found" in str(e).lower():
                print(f"{Fore.YELLOW}⚠️ Extended provider details not available (this is normal for standard registries){Style.RESET_ALL}")
                return {}
            else:
                raise
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Could not fetch extended provider details: {str(e)}{Style.RESET_ALL}")
            return {}
    
    async def delete_version(self, version_str: str) -> bool:
        """Delete a specific version (requires maintainer privileges and registry support)"""
        print(f"{Fore.RED}🗑️ Attempting to delete version {version_str}...{Style.RESET_ALL}")
        
        if not self.config.api_token:
            raise TerraformRegistryError("API token is required for deletion operations")
        
        # Note: This endpoint may not be supported by all registries
        # The standard Provider Registry Protocol does not define deletion endpoints
        endpoint = f"/providers/{self.config.namespace}/{self.config.provider}/{version_str}"
        
        try:
            await self._make_request('DELETE', endpoint)
            print(f"{Fore.GREEN}✅ Successfully deleted version {version_str}{Style.RESET_ALL}")
            return True
        
        except TerraformRegistryError as e:
            if "403" in str(e) or "Access denied" in str(e):
                print(f"{Fore.RED}❌ Access denied: You don't have permission to delete versions for this provider{Style.RESET_ALL}")
                print(f"{Fore.YELLOW}💡 Version deletion requires provider maintainer privileges{Style.RESET_ALL}")
            elif "404" in str(e) or "not found" in str(e).lower():
                print(f"{Fore.RED}❌ Deletion not supported: This registry does not support version deletion{Style.RESET_ALL}")
                print(f"{Fore.YELLOW}💡 The standard Provider Registry Protocol does not define deletion endpoints{Style.RESET_ALL}")
            else:
                print(f"{Fore.RED}❌ Failed to delete version: {str(e)}{Style.RESET_ALL}")
            return False
        except Exception as e:
            print(f"{Fore.RED}❌ Error during deletion: {str(e)}{Style.RESET_ALL}")
            return False

class VersionMatcher:
    """Handles version matching logic"""
    
    @staticmethod
    def normalize_version(version_str: str) -> str:
        """Normalize version string by removing 'v' prefix if present"""
        return version_str.lstrip('v')
    
    @staticmethod
    def match_version(target: str, available_versions: List[Dict[str, Any]]) -> Optional[Dict[str, Any]]:
        """Match target version against available versions"""
        target_normalized = VersionMatcher.normalize_version(target)
        
        # Try exact match first
        for ver in available_versions:
            ver_normalized = VersionMatcher.normalize_version(ver['version'])
            if ver_normalized == target_normalized:
                return ver
        
        # Try semantic version matching
        try:
            target_version = version.parse(target_normalized)
            for ver in available_versions:
                ver_normalized = VersionMatcher.normalize_version(ver['version'])
                try:
                    if version.parse(ver_normalized) == target_version:
                        return ver
                except version.InvalidVersion:
                    continue
        except version.InvalidVersion:
            pass
        
        return None
    
    @staticmethod
    def find_similar_versions(target: str, available_versions: List[Dict[str, Any]], limit: int = 5) -> List[Dict[str, Any]]:
        """Find versions similar to the target version"""
        target_normalized = VersionMatcher.normalize_version(target)
        similar = []
        
        for ver in available_versions:
            ver_normalized = VersionMatcher.normalize_version(ver['version'])
            
            # Check if they share major.minor
            try:
                target_parts = target_normalized.split('.')
                ver_parts = ver_normalized.split('.')
                
                if len(target_parts) >= 2 and len(ver_parts) >= 2:
                    if target_parts[0] == ver_parts[0] and target_parts[1] == ver_parts[1]:
                        similar.append(ver)
            except (IndexError, ValueError):
                continue
        
        # Sort by version and return limited results
        try:
            similar.sort(key=lambda x: version.parse(VersionMatcher.normalize_version(x['version'])), reverse=True)
        except version.InvalidVersion:
            # Fallback to string sorting
            similar.sort(key=lambda x: x['version'], reverse=True)
        
        return similar[:limit]

class TerraformRegistryError(Exception):
    """Custom exception for Terraform Registry operations"""
    pass

class TerraformProviderRegistryManager:
    """Main manager class for Terraform Provider Registry operations"""
    
    def __init__(self, config: Config):
        self.config = config
    
    def _format_platforms(self, platforms: List[Dict[str, str]]) -> str:
        """Format platform information for display"""
        if not platforms:
            return "No platforms"
        
        platform_count = len(platforms)
        if platform_count <= 3:
            platform_strings = [f"{p.get('os', 'unknown')}/{p.get('arch', 'unknown')}" for p in platforms]
            return ", ".join(platform_strings)
        else:
            first_three = [f"{p.get('os', 'unknown')}/{p.get('arch', 'unknown')}" for p in platforms[:3]]
            return f"{', '.join(first_three)}, +{platform_count - 3} more"
    
    def _format_protocols(self, protocols: List[str]) -> str:
        """Format protocol information for display"""
        if not protocols:
            return "[Unknown]"
        return f"[{', '.join(protocols)}]"
    
    async def list_versions(self) -> None:
        """List all available versions for a provider"""
        print(f"\n{Fore.CYAN}📦 Terraform Provider Registry Manager{Style.RESET_ALL}")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{Style.RESET_ALL}")
        print(f"{Fore.CYAN}🔍 Listing versions for provider {self.config.namespace}/{self.config.provider}{Style.RESET_ALL}")
        
        async with TerraformProviderRegistryClient(self.config) as client:
            try:
                # Get provider details first (may not be available)
                provider_details = await client.get_provider_details()
                
                # Get versions
                versions = await client.get_provider_versions()
                
                if not versions:
                    print(f"{Fore.YELLOW}⚠️ No versions found for provider {self.config.namespace}/{self.config.provider}{Style.RESET_ALL}")
                    return
                
                # Display provider details if available
                if provider_details:
                    print(f"\n{Fore.GREEN}📋 Provider Information:{Style.RESET_ALL}")
                    if 'description' in provider_details:
                        print(f"{Fore.YELLOW}   Description: {Fore.WHITE}{provider_details['description']}{Style.RESET_ALL}")
                    if 'source' in provider_details:
                        print(f"{Fore.YELLOW}   Source: {Fore.WHITE}{provider_details['source']}{Style.RESET_ALL}")
                    if 'published_at' in provider_details:
                        try:
                            dt = datetime.fromisoformat(provider_details['published_at'].replace('Z', '+00:00'))
                            formatted_date = dt.strftime('%Y-%m-%d')
                            print(f"{Fore.YELLOW}   Published: {Fore.WHITE}{formatted_date}{Style.RESET_ALL}")
                        except (ValueError, AttributeError):
                            print(f"{Fore.YELLOW}   Published: {Fore.WHITE}{provider_details['published_at']}{Style.RESET_ALL}")
                
                print(f"\n{Fore.GREEN}📋 Available versions ({len(versions)} total):{Style.RESET_ALL}")
                
                # Sort versions by semantic version (newest first)
                try:
                    sorted_versions = sorted(
                        versions, 
                        key=lambda x: version.parse(VersionMatcher.normalize_version(x['version'])), 
                        reverse=True
                    )
                except version.InvalidVersion:
                    # Fallback to string sorting
                    sorted_versions = sorted(versions, key=lambda x: x['version'], reverse=True)
                
                for i, ver in enumerate(sorted_versions):
                    protocols = ver.get('protocols', [])
                    platforms = ver.get('platforms', [])
                    
                    protocol_str = self._format_protocols(protocols)
                    platform_str = self._format_platforms(platforms)
                    
                    print(f"{Fore.CYAN}   {i+1:3d}. {Fore.WHITE}{ver['version']:<15} {Fore.BLUE}🔌 {protocol_str:<15} {Fore.GREEN}💻 {platform_str}{Style.RESET_ALL}")
                
                print(f"\n{Fore.GREEN}✅ Listed {len(versions)} versions successfully{Style.RESET_ALL}")
                print(f"{Fore.YELLOW}💡 Note: Provider Registry API does not include publication dates{Style.RESET_ALL}")
                
            except TerraformRegistryError as e:
                print(f"{Fore.RED}❌ Failed to list versions: {str(e)}{Style.RESET_ALL}")
                sys.exit(1)
            except Exception as e:
                print(f"{Fore.RED}💥 Unexpected error: {str(e)}{Style.RESET_ALL}")
                sys.exit(1)
    
    async def match_and_remove_version(self) -> None:
        """Match a version and remove it if found"""
        if not self.config.target_version:
            print(f"{Fore.RED}❌ No target version specified{Style.RESET_ALL}")
            sys.exit(1)
        
        print(f"\n{Fore.CYAN}📦 Terraform Provider Registry Manager{Style.RESET_ALL}")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━{Style.RESET_ALL}")
        print(f"{Fore.CYAN}🎯 Searching for version {self.config.target_version} of {self.config.namespace}/{self.config.provider}{Style.RESET_ALL}")
        
        async with TerraformProviderRegistryClient(self.config) as client:
            try:
                # Get all versions
                versions = await client.get_provider_versions()
                
                if not versions:
                    print(f"{Fore.YELLOW}⚠️ No versions found for provider {self.config.namespace}/{self.config.provider}{Style.RESET_ALL}")
                    return
                
                # Try to match the version
                print(f"{Fore.YELLOW}🔍 Searching for version match...{Style.RESET_ALL}")
                matched_version = VersionMatcher.match_version(self.config.target_version, versions)
                
                if matched_version:
                    print(f"{Fore.GREEN}✅ Found matching version: {matched_version['version']}{Style.RESET_ALL}")
                    
                    # Display version details
                    print(f"\n{Fore.CYAN}📋 Version Details:{Style.RESET_ALL}")
                    print(f"{Fore.YELLOW}   Version: {Fore.WHITE}{matched_version['version']}{Style.RESET_ALL}")
                    
                    protocols = matched_version.get('protocols', [])
                    if protocols:
                        print(f"{Fore.YELLOW}   Protocols: {Fore.WHITE}{', '.join(protocols)}{Style.RESET_ALL}")
                    
                    platforms = matched_version.get('platforms', [])
                    if platforms:
                        print(f"{Fore.YELLOW}   Platforms: {Fore.WHITE}{len(platforms)} supported{Style.RESET_ALL}")
                        for platform in platforms[:5]:  # Show first 5 platforms
                            os_name = platform.get('os', 'unknown')
                            arch = platform.get('arch', 'unknown')
                            print(f"{Fore.YELLOW}     • {Fore.WHITE}{os_name}/{arch}{Style.RESET_ALL}")
                        if len(platforms) > 5:
                            print(f"{Fore.YELLOW}     • {Fore.WHITE}... and {len(platforms) - 5} more{Style.RESET_ALL}")
                    
                    # Attempt to delete the version
                    print(f"\n{Fore.RED}🗑️ Attempting to remove version {matched_version['version']}...{Style.RESET_ALL}")
                    
                    if not self.config.api_token:
                        print(f"{Fore.YELLOW}⚠️ No API token provided. Version deletion requires authentication.{Style.RESET_ALL}")
                        print(f"{Fore.YELLOW}💡 Use --token parameter to provide your API token{Style.RESET_ALL}")
                        return
                    
                    success = await client.delete_version(matched_version['version'])
                    
                    if success:
                        print(f"\n{Fore.GREEN}🎉 Version {matched_version['version']} has been successfully removed!{Style.RESET_ALL}")
                    else:
                        print(f"\n{Fore.RED}💔 Failed to remove version {matched_version['version']}{Style.RESET_ALL}")
                        print(f"{Fore.YELLOW}💡 Note: Most registries do not support version deletion{Style.RESET_ALL}")
                else:
                    print(f"{Fore.RED}❌ No exact match found for version {self.config.target_version}{Style.RESET_ALL}")
                    
                    # Find similar versions
                    similar = VersionMatcher.find_similar_versions(self.config.target_version, versions)
                    
                    if similar:
                        print(f"\n{Fore.YELLOW}💡 Similar versions found:{Style.RESET_ALL}")
                        for ver in similar:
                            print(f"{Fore.YELLOW}   • {Fore.WHITE}{ver['version']}{Style.RESET_ALL}")
                    else:
                        print(f"{Fore.YELLOW}💡 No similar versions found{Style.RESET_ALL}")
                        
                        # Show latest versions
                        try:
                            latest_versions = sorted(
                                versions, 
                                key=lambda x: version.parse(VersionMatcher.normalize_version(x['version'])), 
                                reverse=True
                            )[:5]
                            
                            print(f"\n{Fore.CYAN}📋 Latest 5 versions:{Style.RESET_ALL}")
                            for ver in latest_versions:
                                print(f"{Fore.CYAN}   • {Fore.WHITE}{ver['version']}{Style.RESET_ALL}")
                        except version.InvalidVersion:
                            pass
                
            except TerraformRegistryError as e:
                print(f"{Fore.RED}❌ Operation failed: {str(e)}{Style.RESET_ALL}")
                sys.exit(1)
            except Exception as e:
                print(f"{Fore.RED}💥 Unexpected error: {str(e)}{Style.RESET_ALL}")
                sys.exit(1)

async def main():
    """Main function to run the script"""
    parser = argparse.ArgumentParser(
        description="Terraform Provider Registry Manager",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # List all versions for a provider
  python terraform_provider_manager.py --namespace hashicorp --provider aws --list

  # Match and remove a specific version (if supported by registry)
  python terraform_provider_manager.py --namespace hashicorp --provider aws --version 5.0.0 --token YOUR_API_TOKEN

  # List versions for a custom provider
  python terraform_provider_manager.py --namespace myorg --provider myprovider --list
  
Note: 
  - Provider Registry API does not include publication dates
  - Version deletion is not part of the standard Provider Registry Protocol
  - Most public registries do not support version deletion
        """
    )
    
    parser.add_argument("--namespace", required=True, help="Provider namespace (e.g., 'hashicorp')")
    parser.add_argument("--provider", required=True, help="Provider name (e.g., 'aws')")
    parser.add_argument("--version", help="Target version to match and remove (e.g., '5.0.0')")
    parser.add_argument("--token", help="API token (required for deletions, if supported)")
    parser.add_argument("--list", action="store_true", help="List all available versions without removal")
    parser.add_argument("--timeout", type=int, default=30, help="Request timeout in seconds")
    parser.add_argument("--max-retries", type=int, default=3, help="Maximum number of retries for failed requests")
    
    args = parser.parse_args()
    
    # Validate arguments
    if not args.list and not args.version:
        print(f"{Fore.RED}❌ Either --list or --version must be specified{Style.RESET_ALL}")
        parser.print_help()
        sys.exit(1)
    
    if args.version and args.list:
        print(f"{Fore.RED}❌ Cannot specify both --list and --version{Style.RESET_ALL}")
        parser.print_help()
        sys.exit(1)
    
    # Create configuration
    config = Config(
        namespace=args.namespace,
        provider=args.provider,
        target_version=args.version,
        api_token=args.token,
        timeout=args.timeout,
        max_retries=args.max_retries
    )
    
    # Create manager and run operation
    manager = TerraformProviderRegistryManager(config)
    
    try:
        if args.list:
            await manager.list_versions()
        else:
            await manager.match_and_remove_version()
        
        print(f"\n{Fore.GREEN}🎉 Operation completed successfully!{Style.RESET_ALL}")
        
    except KeyboardInterrupt:
        print(f"\n{Fore.YELLOW}⚠️ Operation cancelled by user{Style.RESET_ALL}")
        sys.exit(1)
    except Exception as e:
        print(f"\n{Fore.RED}💥 Script execution failed: {str(e)}{Style.RESET_ALL}")
        sys.exit(1)

if __name__ == "__main__":
    asyncio.run(main())