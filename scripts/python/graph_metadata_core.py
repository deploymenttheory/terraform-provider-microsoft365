"""Shared authentication, HTTP, and persistence for Graph metadata exports."""

import asyncio
import json
import os
import time
from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional, Union
from urllib.parse import urlparse

import aiofiles
import aiohttp
import httpx
from azure.core.exceptions import AzureError
from azure.identity import ClientSecretCredential, CredentialUnavailableError
from colorama import Fore, Style
from kiota_abstractions.api_error import APIError
from msgraph import GraphServiceClient
from msgraph_beta import GraphServiceClient as GraphBetaServiceClient

# Failures expected from remote APIs, authentication, SDK availability, or response data.
GRAPH_EXPORT_ERRORS = (
    OSError,
    ValueError,
    TypeError,
    LookupError,
    AttributeError,
    RuntimeError,
    AzureError,
    APIError,
    httpx.HTTPError,
    aiohttp.ClientError,
    SyntaxError,  # ElementTree.ParseError inherits SyntaxError.
)


@dataclass
# Explicit data fields mirror the configuration or external metadata contract.
class Config:  # pylint: disable=too-many-instance-attributes
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
    required_permissions: List[str] = field(
        default_factory=lambda: [
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
            "SecurityEvents.Read.All",
        ]
    )


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
                client_secret=self.config.client_secret,
            )

            await self.refresh_token()
            print(f"{Fore.GREEN}✅ Authentication completed successfully{Style.RESET_ALL}")
            return True

        except CredentialUnavailableError as e:
            print(f"{Fore.RED}❌ Authentication failed: {e}{Style.RESET_ALL}")
            return False
        except GRAPH_EXPORT_ERRORS as e:
            print(f"{Fore.RED}❌ Failed to connect to Microsoft Graph: {e}{Style.RESET_ALL}")
            return False

    async def refresh_token(self):
        """Refresh the access token"""
        token = self.credential.get_token("https://graph.microsoft.com/.default")
        self.token = token.token
        self.token_expires = token.expires_on

    async def get_headers(self) -> Dict[str, str]:
        """Get authentication headers, refreshing token if needed"""
        if time.time() >= self.token_expires - 300:  # Refresh 5 min before expiry
            await self.refresh_token()

        return {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}

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
        retry_after = response_headers.get("Retry-After")
        if retry_after:
            return min(float(retry_after), self.config.max_delay)

        # Exponential backoff with jitter
        delay = min(
            self.config.base_delay * (2**attempt) + (time.time() % 1), self.config.max_delay
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
            limit=self.config.max_concurrent, limit_per_host=5, keepalive_timeout=30
        )
        timeout = aiohttp.ClientTimeout(total=300, connect=30)
        self.session = aiohttp.ClientSession(connector=connector, timeout=timeout)
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.session:
            await self.session.close()

    async def request(
        self, url: str, method: str = "GET", body: Optional[Dict] = None, handle_paging: bool = True
    ) -> Union[Dict, List]:
        """Make HTTP request with retry logic and rate limiting"""

        endpoint = urlparse(url).netloc

        if self.rate_limiter.should_skip_endpoint(endpoint):
            raise RuntimeError(f"Endpoint {endpoint} has too many failures - skipping")

        await self.rate_limiter.wait_if_needed(endpoint)

        all_results = []
        current_url = url

        while True:
            for attempt in range(self.config.max_retries + 1):
                try:
                    headers = await self.auth_manager.get_headers()

                    if method == "POST" and body:
                        async with self.session.post(
                            current_url, headers=headers, json=body
                        ) as response:
                            data = await self._handle_response(response)
                    else:
                        async with self.session.get(current_url, headers=headers) as response:
                            data = await self._handle_response(response)

                    self.rate_limiter.record_success(endpoint)

                    if isinstance(data, dict) and "value" in data:
                        all_results.extend(data["value"])
                        if handle_paging and "@odata.nextLink" in data:
                            current_url = data["@odata.nextLink"]
                            break
                        return all_results
                    return data

                except aiohttp.ClientResponseError as e:
                    if e.status in (429, 503) and attempt < self.config.max_retries:
                        delay = self.rate_limiter.calculate_delay(e.headers, attempt)
                        print(
                            f"{Fore.YELLOW}⚠️ Rate limited on {endpoint}. Retrying in"
                            f" {delay:.1f}s... (Attempt"
                            f" {attempt + 1}/{self.config.max_retries + 1}){Style.RESET_ALL}"
                        )
                        await asyncio.sleep(delay)
                        continue
                    elif e.status == 401:
                        await self.auth_manager.refresh_token()
                        continue
                    self.rate_limiter.record_failure(endpoint)
                    raise RuntimeError(f"HTTP {e.status} error for {url}: {e.message}") from e

                except GRAPH_EXPORT_ERRORS as e:
                    if attempt < self.config.max_retries:
                        delay = self.rate_limiter.calculate_delay({}, attempt)
                        await asyncio.sleep(delay)
                        continue
                    self.rate_limiter.record_failure(endpoint)
                    raise RuntimeError(f"Failed to fetch {url}: {str(e)}") from e

            if not handle_paging or "@odata.nextLink" not in data:
                break

        return all_results

    async def _handle_response(self, response: aiohttp.ClientResponse) -> Union[Dict, str]:
        """Handle response and return appropriate data"""
        response.raise_for_status()

        content_type = response.headers.get("content-type", "")
        if "application/json" in content_type:
            return await response.json()
        return await response.text()


class DataExtractor:  # pylint: disable=too-few-public-methods
    """Base class for extracting specific types of metadata"""

    def __init__(self, graph_sdk: GraphSdkClient, http_client: HttpClient, config: Config):
        self.graph_sdk = graph_sdk
        self.http = http_client
        self.config = config

    async def save_json(self, data: Any, file_path: str) -> None:
        """Save data to JSON file asynchronously"""
        try:
            if data is None:
                print(
                    f"{Fore.YELLOW}  ⚠️ Skipping {os.path.basename(file_path)} - No data to"
                    f" save{Style.RESET_ALL}"
                )
                return

            directory = os.path.dirname(file_path)
            os.makedirs(directory, exist_ok=True)

            # Convert Graph SDK objects to dict for JSON serialization
            if hasattr(data, "__dict__"):
                json_data = self._convert_to_dict(data)
            elif isinstance(data, list):
                json_data = [
                    self._convert_to_dict(item) if hasattr(item, "__dict__") else item
                    for item in data
                ]
            else:
                json_data = data

            async with aiofiles.open(file_path, "w", encoding="utf-8") as f:
                await f.write(json.dumps(json_data, indent=2, default=str))

            print(f"{Fore.GREEN}  ✅ Saved: {os.path.basename(file_path)}{Style.RESET_ALL}")

        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}  ❌ Failed to save: {os.path.basename(file_path)} -"
                f" {str(e)}{Style.RESET_ALL}"
            )
            raise

    def _convert_to_dict(self, obj: Any) -> Dict:
        """Convert Graph SDK objects to dictionaries"""
        if hasattr(obj, "__dict__"):
            result = {}
            for key, value in obj.__dict__.items():
                if key.startswith("_"):
                    continue
                if isinstance(value, list):
                    result[key] = [
                        self._convert_to_dict(item) if hasattr(item, "__dict__") else item
                        for item in value
                    ]
                elif hasattr(value, "__dict__"):
                    result[key] = self._convert_to_dict(value)
                else:
                    result[key] = value
            return result
        return obj
