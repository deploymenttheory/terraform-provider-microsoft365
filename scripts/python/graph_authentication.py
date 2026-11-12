"""Shared client credentials and bounded Graph requests for standalone tools."""

from dataclasses import dataclass, field

import requests
from azure.core.exceptions import AzureError, ClientAuthenticationError
from azure.identity import ClientSecretCredential
from colorama import Fore


@dataclass
class ClientCredentials:
    """Identify an application without exposing its secret in object representations."""

    tenant_id: str
    client_id: str
    client_secret: str = field(repr=False)


def add_credential_arguments(parser) -> None:
    """Add the existing required Entra ID flags in their original order."""
    for flag, help_text in [
        (
            "--tenant-id",
            "Specify the Entra ID tenant ID (Directory ID) where the application is registered",
        ),
        ("--client-id", "Specify the application (client) ID of the Entra ID app registration"),
        ("--client-secret", "Specify the client secret of the Entra ID app registration"),
    ]:
        parser.add_argument(flag, required=True, help=help_text)


class GraphAuthentication:
    """Authenticate once and retain the tools' existing refresh-and-retry behavior."""

    def __init__(self, credentials: ClientCredentials):
        self.credentials = credentials
        self.credential = None
        self.headers = None

    @property
    def tenant_id(self):
        """Return the configured tenant identifier."""
        return self.credentials.tenant_id

    @property
    def client_id(self):
        """Return the configured application identifier."""
        return self.credentials.client_id

    @property
    def client_secret(self):
        """Return the application secret only when authentication needs it."""
        return self.credentials.client_secret

    def get_access_token(self) -> bool:
        """Acquire the existing application token and JSON request headers."""
        print(f"{Fore.CYAN}Connecting to Microsoft Graph...")
        try:
            self.credential = ClientSecretCredential(
                tenant_id=self.tenant_id, client_id=self.client_id, client_secret=self.client_secret
            )
            token = self.credential.get_token("https://graph.microsoft.com/.default")
            self.headers = {
                "Authorization": f"Bearer {token.token}",
                "Content-Type": "application/json",
            }
            print(f"{Fore.GREEN}✅ Connected to Microsoft Graph successfully")
            return True
        except ClientAuthenticationError as error:
            print(f"{Fore.RED}❌ Authentication failed: {error}")
            return False
        except (AzureError, OSError, ValueError) as error:
            print(f"{Fore.RED}❌ Unexpected error: {error}")
            return False

    def _refresh_token_if_needed(self) -> None:
        try:
            if self.credential:
                token = self.credential.get_token("https://graph.microsoft.com/.default")
                self.headers["Authorization"] = f"Bearer {token.token}"
        except (AzureError, OSError, ValueError) as error:
            print(f"{Fore.YELLOW}⚠️ Token refresh warning: {error}")

    @staticmethod
    def _response_json(response, allow_empty):
        return {} if allow_empty and not response.content else response.json()

    def _request_graph(self, method: str, url: str, allow_empty=False, **options):
        """Bound HTTP calls and retain both 401 and request/JSON error refresh retries."""
        options.setdefault("timeout", 30)
        try:
            response = requests.request(method, url, headers=self.headers, **options)
            if response.status_code == 401:
                self._refresh_token_if_needed()
                response = requests.request(method, url, headers=self.headers, **options)
            response.raise_for_status()
            return self._response_json(response, allow_empty)
        except requests.RequestException:
            self._refresh_token_if_needed()
            response = requests.request(method, url, headers=self.headers, **options)
            response.raise_for_status()
            return self._response_json(response, allow_empty)
