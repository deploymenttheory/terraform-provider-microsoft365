"""Protect script refactors with fake HTTP responses and deterministic file fixtures."""

import base64
import contextlib
import io
import json
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest.mock import AsyncMock, MagicMock, Mock, patch

import aiohttp
import requests

from scripts.pipeline.test_support import TemporaryDirectoryTestCase

SCRIPTS = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SCRIPTS))
# Direct script modules resolve their sibling helpers from this directory.
# pylint: disable=wrong-import-position
import graph_authentication as authentication
import upload_dmg_intune_package as upload
from get_m365_service_health import M365ServiceHealthMonitor
from graph_metadata_core import Config, HttpClient

# pylint: enable=wrong-import-position


class GraphRequestTests(unittest.TestCase):
    """Keep bounded retries, response parsing, and constructor options compatible."""

    def test_network_error_refreshes_once_and_keeps_timeouts(self):
        """Verify both network attempts use a timeout and the existing refresh retry."""
        client = authentication.GraphAuthentication(
            authentication.ClientCredentials("tenant", "client", "secret")
        )
        response = Mock(status_code=200, content=b"{}")
        response.json.return_value = {"ok": True}
        with (
            patch.object(
                authentication.requests,
                "request",
                side_effect=[requests.ConnectionError("temporary"), response],
            ) as request,
            patch.object(client, "_refresh_token_if_needed") as refresh,
        ):
            self.assertEqual(client._request_graph("GET", "https://example.invalid"), {"ok": True})
        refresh.assert_called_once()
        self.assertEqual(request.call_count, 2)
        self.assertTrue(all(call.kwargs["timeout"] == 30 for call in request.call_args_list))

    def test_json_error_retains_refresh_retry(self):
        """Verify response JSON errors retain the previous request retry behavior."""
        client = authentication.GraphAuthentication(
            authentication.ClientCredentials("tenant", "client", "secret")
        )
        response = Mock(status_code=200, content=b"{}")
        response.json.side_effect = [requests.JSONDecodeError("invalid", "", 0), {"ok": True}]
        with (
            patch.object(authentication.requests, "request", return_value=response) as request,
            patch.object(client, "_refresh_token_if_needed") as refresh,
        ):
            self.assertEqual(client._request_graph("GET", "https://example.invalid"), {"ok": True})
        self.assertEqual(request.call_count, 2)
        refresh.assert_called_once()

    def test_empty_success_does_not_parse_json_when_allowed(self):
        """Verify the uploader still accepts successful empty Graph responses."""
        client = authentication.GraphAuthentication(
            authentication.ClientCredentials("tenant", "client", "secret")
        )
        response = Mock(status_code=204, content=b"")
        with patch.object(authentication.requests, "request", return_value=response):
            self.assertEqual(
                client._request_graph("PATCH", "https://example.invalid", allow_empty=True), {}
            )
        response.json.assert_not_called()

    def test_health_constructor_retains_positional_options(self):
        """Verify existing positional reporting and permission overrides remain accepted."""
        monitor = M365ServiceHealthMonitor(
            "tenant", "client", "secret", True, "report.json", True, ["issue"], ["Custom.Read"]
        )
        self.assertTrue(monitor.include_resolved)
        self.assertTrue(monitor.problems_only)
        self.assertEqual(monitor.options.export_path, "report.json")
        self.assertEqual(monitor.options.watch_issues, ["issue"])
        self.assertEqual(monitor.required_permissions, ["Custom.Read"])
        self.assertEqual(monitor.client_id, "client")
        self.assertNotIn("secret", repr(monitor.credentials))


class MetadataRequestTests(unittest.IsolatedAsyncioTestCase):
    """Preserve metadata pagination and throttling without external API calls."""

    def make_client(self):
        """Build an HTTP client whose authentication and session are entirely fake."""
        auth = SimpleNamespace(get_headers=AsyncMock(return_value={}), refresh_token=AsyncMock())
        limiter = Mock()
        limiter.should_skip_endpoint.return_value = False
        limiter.wait_if_needed = AsyncMock()
        limiter.calculate_delay.return_value = 0
        client = HttpClient(auth, limiter, Config("tenant", "client", "secret", max_retries=1))
        client.session = MagicMock()
        response = client.session.get.return_value.__aenter__.return_value
        response.headers = {"content-type": "application/json"}
        response.raise_for_status = Mock()
        return client, response

    async def test_pagination_preserves_all_values(self):
        """Verify each page contributes its values before the next-link loop completes."""
        client, response = self.make_client()
        response.json = AsyncMock(
            side_effect=[
                {"value": [1], "@odata.nextLink": "https://example.invalid/second"},
                {"value": [2]},
            ]
        )
        self.assertEqual(await client.request("https://example.invalid/first"), [1, 2])
        self.assertEqual(client.session.get.call_count, 2)

    async def test_throttle_retries_without_recording_a_failure(self):
        """Verify the throttle branch retries instead of falling into the error branch."""
        client, response = self.make_client()
        error = aiohttp.ClientResponseError(
            SimpleNamespace(real_url="https://example.invalid"), (), status=429, message="retry"
        )
        response.raise_for_status.side_effect = [error, None]
        response.json = AsyncMock(return_value={"value": [1]})
        with patch("graph_metadata_core.asyncio.sleep", new_callable=AsyncMock):
            self.assertEqual(await client.request("https://example.invalid"), [1])
        self.assertEqual(client.session.get.call_count, 2)
        client.rate_limiter.record_failure.assert_not_called()


class UploadContractTests(TemporaryDirectoryTestCase):
    """Keep encryption compatible and ensure commit polling always terminates."""

    def setUp(self):
        super().setUp()
        self.resources.enter_context(contextlib.redirect_stdout(io.StringIO()))

    def test_encryption_matches_original_fixture(self):
        """Verify encrypted bytes and metadata against the original implementation's fixture."""
        fixture = json.loads(
            (SCRIPTS / "tests/fixtures/dmg_encryption.json").read_text(encoding="utf-8")
        )
        source = Path(self.directory) / "fixture.dmg"
        source.write_bytes(fixture["source"].encode())
        uploader = upload.IntuneDMGUploader("tenant", "client", "secret")
        with patch.object(
            upload.os,
            "urandom",
            side_effect=[bytes(range(32)), bytes(range(16)), bytes(range(32, 64))],
        ):
            encrypted, metadata = uploader.encrypt_file_for_intune(str(source))
        self.assertEqual(
            base64.b64encode(Path(encrypted).read_bytes()).decode(), fixture["ciphertext_base64"]
        )
        self.assertEqual(metadata, fixture["encryption_info"])

    def test_pending_commit_is_bounded(self):
        """Verify an indefinitely pending response stops after ten polls."""
        uploader = upload.IntuneDMGUploader("tenant", "client", "secret")
        reference = upload.ContentFileReference("app", "macOSDmgApp", "version", "file")
        with (
            patch.object(
                uploader,
                "get_content_file_status",
                return_value={"uploadState": "commitFilePending"},
            ) as status,
            patch.object(uploader, "commit_content_file") as commit,
            patch.object(upload.time, "sleep"),
        ):
            with self.assertRaisesRegex(RuntimeError, "10 attempts"):
                uploader._commit_uploaded_content(reference, {})
        self.assertEqual(status.call_count, 10)
        commit.assert_called_once()
