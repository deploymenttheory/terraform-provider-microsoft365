#!/usr/bin/env python3
"""
Microsoft Intune DMG Package Uploader

A comprehensive tool for uploading macOS DMG packages to Microsoft Intune
using Microsoft Graph APIs.

Usage:
    Basic upload:
        python upload_dmg_intune_package.py --tenant-id "your-tenant" --client-id "your-client"
                                           --client-secret "your-secret" --dmg-file
                                           "path/to/app.dmg"
                                           --app-name "App Name" --app-version "1.0.0"
                                           --bundle-id "com.example.app"

    Full options:
        python upload_dmg_intune_package.py --tenant-id "your-tenant" --client-id "your-client"
                                           --client-secret "your-secret" --dmg-file
                                           "path/to/app.dmg"
                                           --app-name "App Name" --app-version "1.0.0"
                                           --bundle-id "com.example.app" --description "App
                                           Description"
                                           --publisher "Publisher Name" --logo "path/to/logo.png"

Requirements:
    - Azure AD application with DeviceManagementManagedDevices.ReadWrite.All permission
    - Admin consent granted for the application permissions
    - Python packages: requests, colorama, azure-identity, cryptography
"""

import argparse
import base64
import os
import sys
import time
from dataclasses import dataclass
from typing import Any, Dict, NamedTuple, Optional, Tuple

import requests
from azure.core.exceptions import AzureError, ClientAuthenticationError
from azure.identity import ClientSecretCredential
from colorama import Fore, init
from cryptography.exceptions import UnsupportedAlgorithm
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, hmac
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from graph_authentication import ClientCredentials, GraphAuthentication, add_credential_arguments

PACKAGE_UPLOAD_ERRORS = (
    OSError,
    ValueError,
    TypeError,
    LookupError,
    RuntimeError,
    AzureError,
    UnsupportedAlgorithm,
)

# Initialize colorama for cross-platform color support
init(autoreset=True)


@dataclass
class AppPackage:
    """The DMG metadata submitted together as one application package."""

    dmg_file: str
    app_name: str
    app_version: str
    bundle_id: str
    description: Optional[str] = None
    publisher: Optional[str] = None
    logo_file: Optional[str] = None

    def validate(self) -> None:
        """Validate the source and fill the existing description/publisher defaults."""
        if not os.path.exists(self.dmg_file):
            raise FileNotFoundError(f"DMG file not found: {self.dmg_file}")
        if not self.dmg_file.lower().endswith(".dmg"):
            raise ValueError(f"File must be a DMG file: {self.dmg_file}")
        self.description = self.description or f"{self.app_name} {self.app_version}"
        self.publisher = self.publisher or self.app_name

    def graph_payload(self) -> Dict[str, Any]:
        """Serialize the existing macOS DMG application fields for Microsoft Graph."""
        return {
            "@odata.type": "#microsoft.graph.macOSDmgApp",
            "displayName": self.app_name,
            "description": self.description,
            "publisher": self.publisher,
            "fileName": os.path.basename(self.dmg_file),
            "packageIdentifier": self.bundle_id,
            "bundleId": self.bundle_id,
            "versionNumber": self.app_version,
            "primaryBundleId": self.bundle_id,
            "primaryBundleVersion": self.app_version,
            "minimumSupportedOperatingSystem": {
                "@odata.type": "#microsoft.graph.macOSMinimumOperatingSystem",
                "v11_0": True,
            },
            "includedApps": [
                {
                    "@odata.type": "#microsoft.graph.macOSIncludedApp",
                    "bundleId": self.bundle_id,
                    "bundleVersion": self.app_version,
                }
            ],
        }


class ContentFileReference(NamedTuple):
    """Identify the same content file across storage and commit requests."""

    app_id: str
    app_type: str
    content_version_id: str
    file_id: str


def _sha256(data: bytes) -> bytes:
    """Hash source data using the existing Intune file digest format."""
    digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest.update(data)
    return digest.finalize()


def _encrypt_content(data: bytes, key: bytes, iv: bytes) -> bytes:
    """Keep the existing AES-CBC encryption and PKCS7 byte padding."""
    encryptor = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend()).encryptor()
    padding_length = 16 - len(data) % 16
    padded = data + bytes([padding_length] * padding_length)
    return encryptor.update(padded) + encryptor.finalize()


def _content_mac(data: bytes, key: bytes) -> bytes:
    """Authenticate the IV and ciphertext using the existing HMAC format."""
    digest = hmac.HMAC(key, hashes.SHA256(), backend=default_backend())
    digest.update(data)
    return digest.finalize()


def _show_upload_progress(index: int, total: int, block_size: int, file_size_mb: float) -> None:
    """Display one completed block without retaining upload progress state."""
    percent = round((index + 1) / total * 100, 1)
    uploaded = min(round((index + 1) * block_size / (1024 * 1024), 1), file_size_mb)
    filled = int(percent / 2)
    bar = "[" + "▓" * filled + "░" * (50 - filled) + "]"
    print(f"\r{bar} {percent}% ({uploaded} MB / {file_size_mb} MB)", end="")


class IntuneDMGUploader(GraphAuthentication):
    """Microsoft Intune DMG Package Uploader"""

    def __init__(self, tenant_id: str, client_id: str, client_secret: str):
        super().__init__(ClientCredentials(tenant_id, client_id, client_secret))

    def _make_graph_request(self, method: str, url: str, data: dict = None) -> dict:
        """Submit Graph JSON while allowing successful empty response bodies."""
        return self._request_graph(method, url, allow_empty=True, json=data)

    def create_intune_app(self, app_data: Dict[str, Any]) -> Dict[str, Any]:
        """Create a new Intune app"""
        try:
            url = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
            return self._make_graph_request("POST", url, app_data)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error creating Intune app: {e}")
            raise

    def create_content_version(self, app_id: str, app_type: str) -> Dict[str, Any]:
        """Create a content version for an app"""
        try:
            url = (
                f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/"
                f"microsoft.graph.{app_type}/contentVersions"
            )
            return self._make_graph_request("POST", url, {})
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error creating content version: {e}")
            raise

    def create_content_file(
        self, app_id: str, app_type: str, content_version_id: str, file_data: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Create a content file for a content version"""
        try:
            url = (
                f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/"
                f"microsoft.graph.{app_type}/contentVersions/{content_version_id}/files"
            )
            return self._make_graph_request("POST", url, file_data)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error creating content file: {e}")
            raise

    def get_content_file_status(
        self, app_id: str, app_type: str, content_version_id: str, content_file_id: str
    ) -> Dict[str, Any]:
        """Get content file status"""
        try:
            url = (
                f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/"
                f"microsoft.graph.{app_type}/contentVersions/{content_version_id}/files/"
                f"{content_file_id}"
            )
            return self._make_graph_request("GET", url)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error getting content file status: {e}")
            raise

    def commit_content_file(
        self,
        app_id: str,
        app_type: str,
        content_version_id: str,
        content_file_id: str,
        commit_data: Dict[str, Any],
    ) -> Dict[str, Any]:
        """Commit a content file"""
        try:
            url = (
                f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/"
                f"microsoft.graph.{app_type}/contentVersions/{content_version_id}/files/"
                f"{content_file_id}/commit"
            )
            return self._make_graph_request("POST", url, commit_data)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error committing content file: {e}")
            raise

    def update_app_with_content_version(
        self, app_id: str, app_type: str, content_version_id: str
    ) -> Dict[str, Any]:
        """Update app with committed content version"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}"
            update_data = {
                "@odata.type": f"#microsoft.graph.{app_type}",
                "committedContentVersion": content_version_id,
            }
            return self._make_graph_request("PATCH", url, update_data)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error updating app with content version: {e}")
            raise

    def update_app_icon(self, app_id: str, app_type: str, base64_icon: str) -> Dict[str, Any]:
        """Update app icon"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}"
            update_data = {
                "@odata.type": f"#microsoft.graph.{app_type}",
                "largeIcon": {
                    "@odata.type": "#microsoft.graph.mimeContent",
                    "type": "image/png",
                    "value": base64_icon,
                },
            }
            return self._make_graph_request("PATCH", url, update_data)
        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error updating app icon: {e}")
            raise

    def encrypt_file_for_intune(self, source_file: str) -> Tuple[str, Dict[str, str]]:
        """Encrypt file using AES encryption for Intune upload"""
        try:
            target_file = f"{source_file}.bin"

            # Generate encryption key and IV
            encryption_key = os.urandom(32)  # AES-256 key
            iv = os.urandom(16)  # AES-CBC IV
            hmac_key = os.urandom(32)  # HMAC-SHA256 key

            # Read source file
            with open(source_file, "rb") as f:
                source_data = f.read()

            source_digest = _sha256(source_data)
            encrypted_data = _encrypt_content(source_data, encryption_key, iv)
            mac = _content_mac(iv + encrypted_data, hmac_key)

            # Write the encrypted file
            with open(target_file, "wb") as f:
                f.write(mac)  # Write HMAC
                f.write(iv)  # Write IV
                f.write(encrypted_data)  # Write encrypted data

            # Return encryption info
            encryption_info = {
                "encryptionKey": base64.b64encode(encryption_key).decode("utf-8"),
                "fileDigest": base64.b64encode(source_digest).decode("utf-8"),
                "fileDigestAlgorithm": "SHA256",
                "initializationVector": base64.b64encode(iv).decode("utf-8"),
                "mac": base64.b64encode(mac).decode("utf-8"),
                "macKey": base64.b64encode(hmac_key).decode("utf-8"),
                "profileIdentifier": "ProfileVersion1",
            }

            return target_file, encryption_info

        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error encrypting file: {e}")
            raise

    def analyze_encrypted_file(self, encrypted_file: str) -> Dict[str, str]:
        """Analyze encrypted file hex details"""
        try:
            with open(encrypted_file, "rb") as f:
                data = f.read()

            file_length = len(data)
            hmac_bytes = data[:32]
            iv_bytes = data[32:48]
            ciphertext_sample = data[48:64] if len(data) >= 64 else None

            return {
                "file_length": file_length,
                "hmac_hex": hmac_bytes.hex(),
                "iv_hex": iv_bytes.hex(),
                "ciphertext_sample": ciphertext_sample.hex() if ciphertext_sample else "N/A",
                "full_header_hex": " ".join(f"{b:02x}" for b in data[:64]),
            }

        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error analyzing encrypted file: {e}")
            raise

    def upload_to_azure_storage(self, sas_uri: str, file_path: str):
        """Upload file to Azure Storage"""
        try:
            block_size = 8 * 1024 * 1024  # 8 MB
            file_size = os.path.getsize(file_path)
            total_blocks = (file_size + block_size - 1) // block_size

            print(f"\n{Fore.CYAN}⬆️  Uploading to Azure Storage...")
            print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

            file_size_mb = round(file_size / (1024 * 1024), 2)
            print(f"{Fore.YELLOW}📦 File size: {file_size_mb} MB")

            with open(file_path, "rb") as f:
                block_ids = []

                for block_id in range(total_blocks):
                    block_ids.append(self._upload_block(sas_uri, block_id, f.read(block_size)))
                    _show_upload_progress(block_id, total_blocks, block_size, file_size_mb)

            print()  # New line after progress bar

            self._commit_block_list(sas_uri, block_ids)

            print(f"{Fore.GREEN}✅ Upload completed successfully")

        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.RED}❌ Error uploading to Azure Storage: {e}")
            raise

    def _upload_block(self, sas_uri: str, index: int, data: bytes) -> str:
        """Upload one numbered Azure block with a bounded HTTP timeout."""
        block_id = base64.b64encode(str(index).zfill(6).encode()).decode()
        response = requests.put(
            f"{sas_uri}&comp=block&blockid={block_id}",
            headers={"x-ms-blob-type": "BlockBlob"},
            data=data,
            timeout=300,
        )
        response.raise_for_status()
        return block_id

    def _commit_block_list(self, sas_uri: str, block_ids: list[str]) -> None:
        """Commit the same ordered Azure block list after upload completes."""
        block_list = f"""<?xml version="1.0" encoding="utf-8"?>
<BlockList>
{''.join(f'<Latest>{block_id}</Latest>' for block_id in block_ids)}
</BlockList>"""
        response = requests.put(f"{sas_uri}&comp=blocklist", data=block_list, timeout=300)
        response.raise_for_status()

    def get_app_logo(self, app_name: str, local_logo_path: Optional[str] = None) -> Optional[str]:
        """Get app logo for Intune app"""
        try:
            if not local_logo_path or not os.path.exists(local_logo_path):
                print(f"{Fore.YELLOW}⚠️ No valid logo file available")
                return None

            print(f"{Fore.WHITE}Using local logo file for {app_name}: {local_logo_path}")

            with open(local_logo_path, "rb") as f:
                logo_content = base64.b64encode(f.read()).decode("utf-8")

            return logo_content

        except PACKAGE_UPLOAD_ERRORS as e:
            print(f"{Fore.YELLOW}⚠️ Error processing logo: {e}")
            return None

    def _create_dmg_application(self, package: AppPackage) -> str:
        """Display package metadata and create the unchanged DMG application payload."""
        print(f"\n{Fore.CYAN}📋 Application Details:")
        for label, value in [
            ("Display Name", package.app_name),
            ("Version", package.app_version),
            ("Bundle ID", package.bundle_id),
            ("Publisher", package.publisher),
            ("Description", package.description),
            ("File", os.path.basename(package.dmg_file)),
        ]:
            print(f"   • {label}: {value}")
        print(f"\n{Fore.YELLOW}🔄 Creating macOS app (DMG) in Intune...")
        app_id = self.create_intune_app(package.graph_payload())["id"]
        print(f"{Fore.GREEN}✅ App created successfully (ID: {app_id})")
        return app_id

    def _prepare_content(self, package: AppPackage, app_id: str, version: str):
        """Encrypt and register the application content without changing its file metadata."""
        print(f"\n{Fore.YELLOW}🔐 Encrypting application file...")
        encrypted_file, encryption = self.encrypt_file_for_intune(package.dmg_file)
        print(f"{Fore.GREEN}✅ Encryption complete")
        analysis = self.analyze_encrypted_file(encrypted_file)
        print(f"\n{Fore.CYAN}🔍 Encrypted file analysis:")
        print(f"   • File Length: {analysis['file_length']} bytes")
        print(f"   • HMAC (Hex): {analysis['hmac_hex']}")
        print(f"   • IV (Hex): {analysis['iv_hex']}")
        print(f"   • Ciphertext Sample: {analysis['ciphertext_sample']}")
        print(f"   • Full Header Hex: {analysis['full_header_hex']}")
        print(f"\n{Fore.YELLOW}📦 Creating content file...")
        content = {
            "@odata.type": "#microsoft.graph.mobileAppContentFile",
            "name": os.path.basename(package.dmg_file),
            "size": os.path.getsize(package.dmg_file),
            "sizeEncrypted": os.path.getsize(encrypted_file),
            "isDependency": False,
        }
        file_id = self.create_content_file(app_id, "macOSDmgApp", version, content)["id"]
        return (
            ContentFileReference(app_id, "macOSDmgApp", version, file_id),
            encrypted_file,
            encryption,
        )

    def _wait_storage_uri(self, reference: ContentFileReference) -> str:
        """Bound storage URI polling to the existing thirty attempts."""
        print(f"\n{Fore.YELLOW}⏳ Waiting for Azure Storage URI...")
        for attempt in range(30):
            if attempt:
                print(f"Waiting for Azure Storage URI... (Attempt {attempt + 1}/30)")
                time.sleep(5)
            status = self.get_content_file_status(*reference)
            if status["uploadState"] == "azureStorageUriRequestSuccess":
                print(f"{Fore.GREEN}✅ Azure Storage URI received")
                return status["azureStorageUri"]
        raise RuntimeError("Failed to get Azure Storage URI after 30 attempts")

    def _commit_uploaded_content(self, reference: ContentFileReference, encryption: dict) -> None:
        """Bound commit polling even when Graph repeatedly reports a pending state."""
        commit_data = {"fileEncryptionInfo": encryption}
        print(f"\n{Fore.YELLOW}🔄 Committing file...")
        self.commit_content_file(*reference, commit_data)
        print(f"\n{Fore.YELLOW}⏳ Waiting for file commitment to complete...")
        for attempt in range(10):
            time.sleep(10)
            status = self.get_content_file_status(*reference)
            if status["uploadState"] == "commitFileFailed":
                print(f"Commit failed, retrying ({attempt + 1}/10)...")
                self.commit_content_file(*reference, commit_data)
            elif status["uploadState"] == "commitFileSuccess":
                print(f"{Fore.GREEN}✅ File committed successfully")
                return
            print(f"Current state: {status['uploadState']}. Waiting...")
        raise RuntimeError("Failed to commit file after 10 attempts")

    def _add_package_logo(self, package: AppPackage, app_id: str) -> None:
        """Apply an explicitly supplied logo after content is committed."""
        if package.logo_file and os.path.exists(package.logo_file):
            print(f"\n{Fore.YELLOW}🖼️  Adding app logo...")
            logo = self.get_app_logo(package.app_name, package.logo_file)
            if logo:
                self.update_app_icon(app_id, "macOSDmgApp", logo)
                print(f"{Fore.GREEN}✅ Logo added successfully")

    def publish_intune_package(
        self,
        dmg_file: str,
        app_name: str,
        app_version: str,
        bundle_id: str,
        *optional,
        **named_optional,
    ) -> str:
        """Upload the same DMG metadata through explicit bounded content phases."""
        try:
            package = AppPackage(
                dmg_file, app_name, app_version, bundle_id, *optional, **named_optional
            )
            package.validate()
            app_id = self._create_dmg_application(package)
            print(f"\n{Fore.YELLOW}🔒 Processing content version...")
            version = self.create_content_version(app_id, "macOSDmgApp")["id"]
            print(f"{Fore.GREEN}✅ Content version created (ID: {version})")
            reference, encrypted_file, encryption = self._prepare_content(package, app_id, version)
            self.upload_to_azure_storage(self._wait_storage_uri(reference), encrypted_file)
            self._commit_uploaded_content(reference, encryption)
            print(f"\n{Fore.YELLOW}🔄 Updating app with committed content...")
            self.update_app_with_content_version(app_id, "macOSDmgApp", version)
            print(f"{Fore.GREEN}✅ App updated successfully")
            self._add_package_logo(package, app_id)
            print(f"\n{Fore.YELLOW}🧹 Cleaning up temporary files...")
            if os.path.exists(encrypted_file):
                os.remove(encrypted_file)
            print(f"{Fore.GREEN}✅ Cleanup complete")
            print(f"\n{Fore.CYAN}✨ Successfully uploaded {app_name} to Intune")
            print(
                f"{Fore.CYAN}🔗 Intune Portal URL: "
                + "https://intune.microsoft.com/#view/Microsoft_Intune_Apps/"
                + f"SettingsMenu/~/0/appId/{app_id}"
            )
            return app_id
        except PACKAGE_UPLOAD_ERRORS as error:
            print(f"{Fore.RED}❌ Error publishing package to Intune: {error}")
            raise


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description="Microsoft Intune DMG Package Uploader",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )

    add_credential_arguments(parser)

    parser.add_argument("--dmg-file", required=True, help="Path to the DMG file to upload")

    parser.add_argument("--app-name", required=True, help="Display name for the application")

    parser.add_argument("--app-version", required=True, help="Version number of the application")

    parser.add_argument("--bundle-id", required=True, help="Bundle ID of the application")

    parser.add_argument("--description", help="Description of the application")

    parser.add_argument("--publisher", help="Publisher of the application")

    parser.add_argument("--logo", help="Path to a logo PNG file for the application")

    args = parser.parse_args()

    try:
        print(f"\n{Fore.CYAN}📦 Starting DMG upload process...")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

        uploader = IntuneDMGUploader(
            tenant_id=args.tenant_id, client_id=args.client_id, client_secret=args.client_secret
        )

        # Authenticate
        if not uploader.get_access_token():
            sys.exit(2)

        # Upload package
        uploader.publish_intune_package(
            dmg_file=args.dmg_file,
            app_name=args.app_name,
            app_version=args.app_version,
            bundle_id=args.bundle_id,
            description=args.description,
            publisher=args.publisher,
            logo_file=args.logo,
        )

        print(f"\n{Fore.GREEN}🎉 DMG upload process completed successfully!")

    except PACKAGE_UPLOAD_ERRORS as e:
        print(f"\n{Fore.RED}❌ DMG upload process failed: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
