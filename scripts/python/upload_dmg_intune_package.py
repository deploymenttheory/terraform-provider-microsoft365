#!/usr/bin/env python3
"""
Microsoft Intune DMG Package Uploader

A comprehensive tool for uploading macOS DMG packages to Microsoft Intune
using Microsoft Graph APIs.

Usage:
    Basic upload:
        python upload_dmg_intune_package.py --tenant-id "your-tenant" --client-id "your-client" 
                                           --client-secret "your-secret" --dmg-file "path/to/app.dmg"
                                           --app-name "App Name" --app-version "1.0.0" 
                                           --bundle-id "com.example.app"
    
    Full options:
        python upload_dmg_intune_package.py --tenant-id "your-tenant" --client-id "your-client" 
                                           --client-secret "your-secret" --dmg-file "path/to/app.dmg"
                                           --app-name "App Name" --app-version "1.0.0" 
                                           --bundle-id "com.example.app" --description "App Description"
                                           --publisher "Publisher Name" --logo "path/to/logo.png"

Requirements:
    - Azure AD application with DeviceManagementManagedDevices.ReadWrite.All permission
    - Admin consent granted for the application permissions
    - Python packages: requests, colorama, azure-identity, cryptography
"""

import argparse
import base64
import json
import os
import sys
import time
import urllib.parse
from datetime import datetime, UTC
from typing import Dict, List, Optional, Any, Tuple

import requests
from colorama import init, Fore, Style
from azure.identity import ClientSecretCredential
from azure.core.exceptions import ClientAuthenticationError
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import hashes, hmac
from cryptography.hazmat.backends import default_backend

# Initialize colorama for cross-platform color support
init(autoreset=True)

class IntuneDMGUploader:
    """Microsoft Intune DMG Package Uploader"""
    
    def __init__(self, tenant_id: str, client_id: str, client_secret: str):
        self.tenant_id = tenant_id
        self.client_id = client_id
        self.client_secret = client_secret
        self.credential = None
        self.headers = None
        
    def get_access_token(self) -> bool:
        """Authenticate with Microsoft Graph using Azure Identity SDK"""
        print(f"{Fore.CYAN}Connecting to Microsoft Graph...")
        
        try:
            self.credential = ClientSecretCredential(
                tenant_id=self.tenant_id,
                client_id=self.client_id,
                client_secret=self.client_secret
            )
            
            token = self.credential.get_token("https://graph.microsoft.com/.default")
            
            self.headers = {
                'Authorization': f'Bearer {token.token}',
                'Content-Type': 'application/json'
            }
            
            print(f"{Fore.GREEN}✅ Connected to Microsoft Graph successfully")
            return True
            
        except ClientAuthenticationError as e:
            print(f"{Fore.RED}❌ Authentication failed: {e}")
            return False
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to connect to Microsoft Graph: {e}")
            return False
    
    def _refresh_token_if_needed(self):
        """Refresh access token if needed"""
        try:
            if self.credential:
                token = self.credential.get_token("https://graph.microsoft.com/.default")
                self.headers['Authorization'] = f'Bearer {token.token}'
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Token refresh warning: {e}")
    
    def _make_graph_request(self, method: str, url: str, data: dict = None) -> dict:
        """Make a Microsoft Graph API request with automatic token refresh"""
        try:
            response = requests.request(method, url, headers=self.headers, json=data)
            
            if response.status_code == 401:
                self._refresh_token_if_needed()
                response = requests.request(method, url, headers=self.headers, json=data)
            
            response.raise_for_status()
            return response.json() if response.content else {}
            
        except requests.exceptions.RequestException:
            self._refresh_token_if_needed()
            response = requests.request(method, url, headers=self.headers, json=data)
            response.raise_for_status()
            return response.json() if response.content else {}
    
    def create_intune_app(self, app_data: Dict[str, Any]) -> Dict[str, Any]:
        """Create a new Intune app"""
        try:
            url = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
            return self._make_graph_request("POST", url, app_data)
        except Exception as e:
            print(f"{Fore.RED}❌ Error creating Intune app: {e}")
            raise
    
    def create_content_version(self, app_id: str, app_type: str) -> Dict[str, Any]:
        """Create a content version for an app"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/microsoft.graph.{app_type}/contentVersions"
            return self._make_graph_request("POST", url, {})
        except Exception as e:
            print(f"{Fore.RED}❌ Error creating content version: {e}")
            raise
    
    def create_content_file(self, app_id: str, app_type: str, content_version_id: str, file_data: Dict[str, Any]) -> Dict[str, Any]:
        """Create a content file for a content version"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/microsoft.graph.{app_type}/contentVersions/{content_version_id}/files"
            return self._make_graph_request("POST", url, file_data)
        except Exception as e:
            print(f"{Fore.RED}❌ Error creating content file: {e}")
            raise
    
    def get_content_file_status(self, app_id: str, app_type: str, content_version_id: str, content_file_id: str) -> Dict[str, Any]:
        """Get content file status"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/microsoft.graph.{app_type}/contentVersions/{content_version_id}/files/{content_file_id}"
            return self._make_graph_request("GET", url)
        except Exception as e:
            print(f"{Fore.RED}❌ Error getting content file status: {e}")
            raise
    
    def commit_content_file(self, app_id: str, app_type: str, content_version_id: str, content_file_id: str, commit_data: Dict[str, Any]) -> Dict[str, Any]:
        """Commit a content file"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}/microsoft.graph.{app_type}/contentVersions/{content_version_id}/files/{content_file_id}/commit"
            return self._make_graph_request("POST", url, commit_data)
        except Exception as e:
            print(f"{Fore.RED}❌ Error committing content file: {e}")
            raise
    
    def update_app_with_content_version(self, app_id: str, app_type: str, content_version_id: str) -> Dict[str, Any]:
        """Update app with committed content version"""
        try:
            url = f"https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/{app_id}"
            update_data = {
                "@odata.type": f"#microsoft.graph.{app_type}",
                "committedContentVersion": content_version_id
            }
            return self._make_graph_request("PATCH", url, update_data)
        except Exception as e:
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
                    "value": base64_icon
                }
            }
            return self._make_graph_request("PATCH", url, update_data)
        except Exception as e:
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
            with open(source_file, 'rb') as f:
                source_data = f.read()
            
            # Calculate source file hash
            source_hash = hashes.Hash(hashes.SHA256(), backend=default_backend())
            source_hash.update(source_data)
            source_digest = source_hash.finalize()
            
            # Encrypt the data
            cipher = Cipher(algorithms.AES(encryption_key), modes.CBC(iv), backend=default_backend())
            encryptor = cipher.encryptor()
            
            # Pad the data to AES block size
            block_size = 16
            padding_length = block_size - (len(source_data) % block_size)
            padded_data = source_data + bytes([padding_length] * padding_length)
            
            encrypted_data = encryptor.update(padded_data) + encryptor.finalize()
            
            # Calculate HMAC
            h = hmac.HMAC(hmac_key, hashes.SHA256(), backend=default_backend())
            h.update(iv + encrypted_data)
            mac = h.finalize()
            
            # Write the encrypted file
            with open(target_file, 'wb') as f:
                f.write(mac)  # Write HMAC
                f.write(iv)   # Write IV
                f.write(encrypted_data)  # Write encrypted data
            
            # Return encryption info
            encryption_info = {
                "encryptionKey": base64.b64encode(encryption_key).decode('utf-8'),
                "fileDigest": base64.b64encode(source_digest).decode('utf-8'),
                "fileDigestAlgorithm": "SHA256",
                "initializationVector": base64.b64encode(iv).decode('utf-8'),
                "mac": base64.b64encode(mac).decode('utf-8'),
                "macKey": base64.b64encode(hmac_key).decode('utf-8'),
                "profileIdentifier": "ProfileVersion1"
            }
            
            return target_file, encryption_info
            
        except Exception as e:
            print(f"{Fore.RED}❌ Error encrypting file: {e}")
            raise
    
    def analyze_encrypted_file(self, encrypted_file: str) -> Dict[str, str]:
        """Analyze encrypted file hex details"""
        try:
            with open(encrypted_file, 'rb') as f:
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
                "full_header_hex": " ".join(f"{b:02x}" for b in data[:64])
            }
            
        except Exception as e:
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
            
            with open(file_path, 'rb') as f:
                block_ids = []
                
                for block_id in range(total_blocks):
                    block_data = f.read(block_size)
                    block_id_base64 = base64.b64encode(str(block_id).zfill(6).encode()).decode()
                    
                    # Upload block
                    block_url = f"{sas_uri}&comp=block&blockid={block_id_base64}"
                    response = requests.put(
                        block_url,
                        headers={"x-ms-blob-type": "BlockBlob"},
                        data=block_data
                    )
                    response.raise_for_status()
                    
                    block_ids.append(block_id_base64)
                    
                    # Show progress
                    percent_complete = round((block_id + 1) / total_blocks * 100, 1)
                    uploaded_mb = min(round((block_id + 1) * block_size / (1024 * 1024), 1), file_size_mb)
                    
                    progress_width = 50
                    filled_blocks = int(percent_complete / 2)
                    empty_blocks = progress_width - filled_blocks
                    progress_bar = "[" + "▓" * filled_blocks + "░" * empty_blocks + "]"
                    
                    print(f"\r{progress_bar} {percent_complete}% ({uploaded_mb} MB / {file_size_mb} MB)", end="")
            
            print()  # New line after progress bar
            
            # Commit block list
            block_list_xml = f"""<?xml version="1.0" encoding="utf-8"?>
<BlockList>
{''.join(f'<Latest>{block_id}</Latest>' for block_id in block_ids)}
</BlockList>"""
            
            commit_url = f"{sas_uri}&comp=blocklist"
            response = requests.put(commit_url, data=block_list_xml)
            response.raise_for_status()
            
            print(f"{Fore.GREEN}✅ Upload completed successfully")
            
        except Exception as e:
            print(f"{Fore.RED}❌ Error uploading to Azure Storage: {e}")
            raise
    
    def get_app_logo(self, app_name: str, local_logo_path: Optional[str] = None) -> Optional[str]:
        """Get app logo for Intune app"""
        try:
            if not local_logo_path or not os.path.exists(local_logo_path):
                print(f"{Fore.YELLOW}⚠️ No valid logo file available")
                return None
            
            print(f"{Fore.WHITE}Using local logo file: {local_logo_path}")
            
            with open(local_logo_path, 'rb') as f:
                logo_content = base64.b64encode(f.read()).decode('utf-8')
            
            return logo_content
            
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Error processing logo: {e}")
            return None
    
    def publish_intune_package(self, dmg_file: str, app_name: str, app_version: str, 
                            bundle_id: str, description: Optional[str] = None,
                            publisher: Optional[str] = None, logo_file: Optional[str] = None) -> str:
        """Main function to upload a DMG file to Intune"""
        try:
            # Validate file exists and is DMG
            if not os.path.exists(dmg_file):
                raise FileNotFoundError(f"DMG file not found: {dmg_file}")
            
            if not dmg_file.lower().endswith('.dmg'):
                raise ValueError(f"File must be a DMG file: {dmg_file}")
            
            # Set defaults if not provided
            if not description:
                description = f"{app_name} {app_version}"
            
            if not publisher:
                publisher = app_name
            
            file_name = os.path.basename(dmg_file)
            app_type = "macOSDmgApp"
            
            print(f"\n{Fore.CYAN}📋 Application Details:")
            print(f"   • Display Name: {app_name}")
            print(f"   • Version: {app_version}")
            print(f"   • Bundle ID: {bundle_id}")
            print(f"   • Publisher: {publisher}")
            print(f"   • Description: {description}")
            print(f"   • File: {file_name}")
            
            # Step 1: Create the app in Intune
            print(f"\n{Fore.YELLOW}🔄 Creating macOS app (DMG) in Intune...")
            
            app = {
                "@odata.type": f"#microsoft.graph.{app_type}",
                "displayName": app_name,
                "description": description,
                "publisher": publisher,
                "fileName": file_name,
                "packageIdentifier": bundle_id,
                "bundleId": bundle_id,
                "versionNumber": app_version,
                "primaryBundleId": bundle_id,
                "primaryBundleVersion": app_version,
                "minimumSupportedOperatingSystem": {
                    "@odata.type": "#microsoft.graph.macOSMinimumOperatingSystem",
                    "v11_0": True
                },
                "includedApps": [
                    {
                        "@odata.type": "#microsoft.graph.macOSIncludedApp",
                        "bundleId": bundle_id,
                        "bundleVersion": app_version
                    }
                ]
            }
            
            new_app = self.create_intune_app(app)
            app_id = new_app['id']
            print(f"{Fore.GREEN}✅ App created successfully (ID: {app_id})")
            
            # Step 2: Create content version
            print(f"\n{Fore.YELLOW}🔒 Processing content version...")
            content_version = self.create_content_version(app_id, app_type)
            content_version_id = content_version['id']
            print(f"{Fore.GREEN}✅ Content version created (ID: {content_version_id})")
            
            # Step 3: Encrypt the file
            print(f"\n{Fore.YELLOW}🔐 Encrypting application file...")
            encrypted_file, file_encryption_info = self.encrypt_file_for_intune(dmg_file)
            print(f"{Fore.GREEN}✅ Encryption complete")
            
            # Analyze the encrypted file
            analysis = self.analyze_encrypted_file(encrypted_file)
            print(f"\n{Fore.CYAN}🔍 Encrypted file analysis:")
            print(f"   • File Length: {analysis['file_length']} bytes")
            print(f"   • HMAC (Hex): {analysis['hmac_hex']}")
            print(f"   • IV (Hex): {analysis['iv_hex']}")
            print(f"   • Ciphertext Sample: {analysis['ciphertext_sample']}")
            print(f"   • Full Header Hex: {analysis['full_header_hex']}")
            
            # Step 4: Create content file
            print(f"\n{Fore.YELLOW}📦 Creating content file...")
            file_content = {
                "@odata.type": "#microsoft.graph.mobileAppContentFile",
                "name": file_name,
                "size": os.path.getsize(dmg_file),
                "sizeEncrypted": os.path.getsize(encrypted_file),
                "isDependency": False
            }
            
            content_file = self.create_content_file(app_id, app_type, content_version_id, file_content)
            content_file_id = content_file['id']
            
            # Step 5: Wait for Azure Storage Uri
            print(f"\n{Fore.YELLOW}⏳ Waiting for Azure Storage URI...")
            
            attempts = 0
            max_attempts = 30
            file_status = None
            
            while attempts < max_attempts:
                if attempts > 0:
                    print(f"Waiting for Azure Storage URI... (Attempt {attempts + 1}/{max_attempts})")
                    time.sleep(5)
                
                file_status = self.get_content_file_status(app_id, app_type, content_version_id, content_file_id)
                if file_status['uploadState'] == 'azureStorageUriRequestSuccess':
                    break
                
                attempts += 1
            
            if file_status['uploadState'] != 'azureStorageUriRequestSuccess':
                raise Exception(f"Failed to get Azure Storage URI after {max_attempts} attempts")
            
            print(f"{Fore.GREEN}✅ Azure Storage URI received")
            
            # Step 6: Upload file to Azure Storage
            self.upload_to_azure_storage(file_status['azureStorageUri'], encrypted_file)
            
            # Step 7: Commit the file
            print(f"\n{Fore.YELLOW}🔄 Committing file...")
            commit_data = {
                "fileEncryptionInfo": file_encryption_info
            }
            
            self.commit_content_file(app_id, app_type, content_version_id, content_file_id, commit_data)
            
            # Step 8: Wait for commit to complete
            print(f"\n{Fore.YELLOW}⏳ Waiting for file commitment to complete...")
            retry_count = 0
            max_retries = 10
            
            while retry_count < max_retries:
                time.sleep(10)
                file_status = self.get_content_file_status(app_id, app_type, content_version_id, content_file_id)
                
                if file_status['uploadState'] == 'commitFileFailed':
                    retry_count += 1
                    print(f"Commit failed, retrying ({retry_count}/{max_retries})...")
                    self.commit_content_file(app_id, app_type, content_version_id, content_file_id, commit_data)
                elif file_status['uploadState'] == 'commitFileSuccess':
                    print(f"{Fore.GREEN}✅ File committed successfully")
                    break
                else:
                    print(f"Current state: {file_status['uploadState']}. Waiting...")
            
            if file_status['uploadState'] != 'commitFileSuccess':
                raise Exception(f"Failed to commit file after {max_retries} attempts")
            
            # Step 9: Update app with committed content version
            print(f"\n{Fore.YELLOW}🔄 Updating app with committed content...")
            self.update_app_with_content_version(app_id, app_type, content_version_id)
            print(f"{Fore.GREEN}✅ App updated successfully")
            
            # Step 10: Add logo if one was provided
            if logo_file and os.path.exists(logo_file):
                print(f"\n{Fore.YELLOW}🖼️  Adding app logo...")
                logo_content = self.get_app_logo(app_name, logo_file)
                if logo_content:
                    self.update_app_icon(app_id, app_type, logo_content)
                    print(f"{Fore.GREEN}✅ Logo added successfully")
            
            # Step 11: Clean up temporary files
            print(f"\n{Fore.YELLOW}🧹 Cleaning up temporary files...")
            if os.path.exists(encrypted_file):
                os.remove(encrypted_file)
            print(f"{Fore.GREEN}✅ Cleanup complete")
            
            # Step 12: Final success message
            print(f"\n{Fore.CYAN}✨ Successfully uploaded {app_name} to Intune")
            print(f"{Fore.CYAN}🔗 Intune Portal URL: https://intune.microsoft.com/#view/Microsoft_Intune_Apps/SettingsMenu/~/0/appId/{app_id}")
            
            return app_id
            
        except Exception as e:
            print(f"{Fore.RED}❌ Error publishing package to Intune: {e}")
            raise


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description="Microsoft Intune DMG Package Uploader",
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
        "--dmg-file",
        required=True,
        help="Path to the DMG file to upload"
    )
    
    parser.add_argument(
        "--app-name",
        required=True,
        help="Display name for the application"
    )
    
    parser.add_argument(
        "--app-version",
        required=True,
        help="Version number of the application"
    )
    
    parser.add_argument(
        "--bundle-id",
        required=True,
        help="Bundle ID of the application"
    )
    
    parser.add_argument(
        "--description",
        help="Description of the application"
    )
    
    parser.add_argument(
        "--publisher",
        help="Publisher of the application"
    )
    
    parser.add_argument(
        "--logo",
        help="Path to a logo PNG file for the application"
    )
    
    args = parser.parse_args()
    
    try:
        print(f"\n{Fore.CYAN}📦 Starting DMG upload process...")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        
        uploader = IntuneDMGUploader(
            tenant_id=args.tenant_id,
            client_id=args.client_id,
            client_secret=args.client_secret
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
            logo_file=args.logo
        )
        
        print(f"\n{Fore.GREEN}🎉 DMG upload process completed successfully!")
        
    except Exception as e:
        print(f"\n{Fore.RED}❌ DMG upload process failed: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
