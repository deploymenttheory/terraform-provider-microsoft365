"""Portal and public endpoint metadata exports."""

import asyncio
import json
import os
import re
import shutil
import xml.etree.ElementTree as ET
from typing import Dict

import aiofiles
from colorama import Fore, Style
from graph_metadata_core import GRAPH_EXPORT_ERRORS, DataExtractor


class PortalExtractor(DataExtractor):
    """Extracts data from various Microsoft portals using custom HTTP client"""

    async def extract_setting_status_errors(self, export_path: str) -> None:
        """Extract Setting Status Errors"""
        print(f"{Fore.CYAN}🔍 Extracting Setting Status Errors...{Style.RESET_ALL}")

        try:
            # Get version information from Intune portal
            version_response = await self.http.request(
                "https://intune.microsoft.com/signin/idpRedirect.js"
            )

            if isinstance(version_response, str):
                match = re.search(r'"extensionsPageVersion":(\{[^}]+\})', version_response)
                if match:
                    versions = json.loads(match.group(1))
                    device_settings_version = versions["Microsoft_Intune_DeviceSettings"][0]

                    root = "https://afd-v2.hosting.portal.azure.net"
                    setting_status_url = f"{root}/intunedevicesettings/Content/{device_settings_version}/Scripts/DeviceConfiguration/Blades/DevicePoliciesStatus/SettingStatus.js"

                    setting_status_response = await self.http.request(setting_status_url)

                    if isinstance(setting_status_response, str):
                        match = re.search(
                            r"SettingStatusErrorMap\s*=\s*(\{.*?\});",
                            setting_status_response,
                            re.DOTALL,
                        )
                        if match:
                            js_object = match.group(1)
                            try:
                                # Clean the JavaScript object string for JSON parsing
                                cleaned_js = re.sub(r"[\x00-\x1f\x7f-\x9f]", "", js_object)
                                cleaned_js = re.sub(
                                    r"([{,]\s*)([a-zA-Z_$][a-zA-Z0-9_$]*)\s*:",
                                    r'\1"\2":',
                                    cleaned_js,
                                )
                                cleaned_js = re.sub(r"'([^']*)'", r'"\1"', cleaned_js)

                                error_map = json.loads(cleaned_js)
                                await self.save_json(
                                    error_map, os.path.join(export_path, "SettingStatusErrors.json")
                                )
                            except json.JSONDecodeError:
                                print(
                                    f"{Fore.YELLOW}  ⚠️ Failed to parse SettingStatusErrorMap as"
                                    f" JSON{Style.RESET_ALL}"
                                )
                                async with aiofiles.open(
                                    os.path.join(export_path, "SettingStatusErrors_raw.txt"),
                                    "w",
                                    encoding="utf-8",
                                ) as f:
                                    await f.write(js_object)
                                print(
                                    f"{Fore.YELLOW}  ⚠️ Saved raw JavaScript object to"
                                    f" SettingStatusErrors_raw.txt{Style.RESET_ALL}"
                                )

            print(f"{Fore.GREEN}✅ Setting Status Errors extraction completed{Style.RESET_ALL}")
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.RED}❌ Failed to extract Setting Status Errors: {str(e)}{Style.RESET_ALL}"
            )
            raise

    async def extract_dcv1_policies(self, export_path: str) -> None:
        """Extract Device Configuration v1 Policies"""
        print(f"{Fore.CYAN}🔍 Extracting Device Configuration v1 Policies...{Style.RESET_ALL}")

        # Get version information
        version_response = await self.http.request(
            "https://intune.microsoft.com/signin/idpRedirect.js"
        )

        if isinstance(version_response, str):
            match = re.search(r'"extensionsPageVersion":(\{[^}]+\})', version_response)
            if match:
                versions = json.loads(match.group(1))
                device_settings_version = versions["Microsoft_Intune_DeviceSettings"][0]

                root = "https://afd-v2.hosting.portal.azure.net"
                root_device_settings = f"{root}/intunedevicesettings/Content/{device_settings_version}/Scripts/DeviceConfiguration"

                # Clean existing DCv1 directory
                dcv1_path = os.path.join(export_path, "DCv1")
                if os.path.exists(dcv1_path):
                    shutil.rmtree(dcv1_path)

                # Process Configuration and Compliance concurrently
                tasks = [
                    self._extract_dcv1_metadata("Configuration", root_device_settings, dcv1_path),
                    self._extract_dcv1_metadata("Compliance", root_device_settings, dcv1_path),
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

        if isinstance(metadata_response, str) and "metadata = " in metadata_response:
            match = re.search(r"(?s)metadata = (\{.+\});", metadata_response)
            if match:
                try:
                    metadata = json.loads(match.group(1))

                    tasks = []
                    for family_data in metadata.values():
                        if isinstance(family_data, list):
                            for setting in family_data:
                                if isinstance(setting, dict) and "id" in setting:
                                    clean_id = "_".join(setting["id"].split("_")[:-1])
                                    setting["id"] = clean_id

                                    # Clean nested IDs
                                    setting = self._remove_dcv1_version_suffixes(setting)

                                    file_path = os.path.join(source_path, f"{clean_id}.json")
                                    tasks.append(self.save_json(setting, file_path))

                    # Process in batches
                    batch_size = 20
                    for i in range(0, len(tasks), batch_size):
                        batch = tasks[i : i + batch_size]
                        await asyncio.gather(*batch, return_exceptions=True)

                except json.JSONDecodeError as e:
                    print(
                        f"{Fore.YELLOW}  ⚠️ Failed to parse {source} metadata: {e}{Style.RESET_ALL}"
                    )

    def _remove_dcv1_version_suffixes(self, setting: Dict) -> Dict:
        """Remove version suffixes from DCv1 setting IDs"""
        # Clean child settings
        if "childSettings" in setting and setting["childSettings"]:
            for child in setting["childSettings"]:
                if "id" in child:
                    child["id"] = "_".join(child["id"].split("_")[:-1])
                child = self._remove_dcv1_version_suffixes(child)

        # Clean options
        if "options" in setting and setting["options"]:
            for option in setting["options"]:
                if "children" in option and option["children"]:
                    for child in option["children"]:
                        if "id" in child:
                            child["id"] = "_".join(child["id"].split("_")[:-1])
                        child = self._remove_dcv1_version_suffixes(child)

        # Clean complex options
        if "complexOptions" in setting and setting["complexOptions"]:
            for complex_option in setting["complexOptions"]:
                if "id" in complex_option:
                    complex_option["id"] = "_".join(complex_option["id"].split("_")[:-1])
                complex_option = self._remove_dcv1_version_suffixes(complex_option)

        # Clean columns
        if "columns" in setting and setting["columns"]:
            for column in setting["columns"]:
                if "metadata" in column and column["metadata"]:
                    if "id" in column["metadata"]:
                        column["metadata"]["id"] = "_".join(
                            column["metadata"]["id"].split("_")[:-1]
                        )
                    column["metadata"] = self._remove_dcv1_version_suffixes(column["metadata"])

        return setting

    async def extract_identity_product_changes(self, export_path: str) -> None:
        """Extract Identity Product Changes from RSS feed"""
        print(f"{Fore.CYAN}🔍 Extracting Identity Product Changes...{Style.RESET_ALL}")

        try:
            # Use RSS feed for Microsoft Entra release notes
            rss_url = (
                "https://learn.microsoft.com/api/search/rss?"
                + "search=%22Release+notes+-+Azure+Active+Directory%22&"
                + "locale=en-us"
            )
            print(f"{Fore.YELLOW}  📡 Fetching from Microsoft Entra RSS feed...{Style.RESET_ALL}")

            rss_response = await self.http.request(rss_url)

            if isinstance(rss_response, str):
                try:
                    # Parse RSS/XML content
                    root = ET.fromstring(rss_response)
                    release_notes = []

                    # Handle both RSS and Atom feeds
                    items = root.findall(".//item") or root.findall(
                        ".//{http://www.w3.org/2005/Atom}entry"
                    )

                    for item in items:
                        # RSS format
                        title_elem = item.find("title")
                        if title_elem is None:
                            title_elem = item.find(".//{http://www.w3.org/2005/Atom}title")

                        link_elem = item.find("link")
                        if link_elem is None:
                            link_elem = item.find(".//{http://www.w3.org/2005/Atom}link")

                        pub_date_elem = item.find("pubDate")
                        if pub_date_elem is None:
                            pub_date_elem = item.find(".//{http://www.w3.org/2005/Atom}published")

                        desc_elem = item.find("description")
                        if desc_elem is None:
                            desc_elem = item.find(".//{http://www.w3.org/2005/Atom}summary")

                        guid_elem = item.find("guid")
                        if guid_elem is None:
                            guid_elem = item.find(".//{http://www.w3.org/2005/Atom}id")

                        release_note = {
                            "title": title_elem.text if title_elem is not None else "",
                            "link": (
                                link_elem.text
                                if link_elem is not None
                                else (link_elem.get("href") if link_elem is not None else "")
                            ),
                            "published": pub_date_elem.text if pub_date_elem is not None else "",
                            "summary": desc_elem.text if desc_elem is not None else "",
                            "id": guid_elem.text if guid_elem is not None else "",
                            "source": "RSS Feed",
                        }
                        release_notes.append(release_note)

                    if release_notes:
                        # Sort by publication date (newest first)
                        release_notes.sort(key=lambda x: x.get("published", ""), reverse=True)

                        await self.save_json(
                            release_notes, os.path.join(export_path, "IdentityProductChanges.json")
                        )
                        print(
                            f"{Fore.GREEN}✅ Identity Product Changes extraction completed from RSS"
                            f" feed ({len(release_notes)} items){Style.RESET_ALL}"
                        )

                        # Also save raw RSS for reference
                        async with aiofiles.open(
                            os.path.join(export_path, "IdentityProductChanges_raw.xml"),
                            "w",
                            encoding="utf-8",
                        ) as f:
                            await f.write(rss_response)
                        print(f"{Fore.GREEN}  ✅ Raw RSS feed saved for reference{Style.RESET_ALL}")
                    else:
                        print(
                            f"{Fore.YELLOW}⚠️ No release notes found in RSS feed{Style.RESET_ALL}"
                        )

                except ET.ParseError as parse_error:
                    print(
                        f"{Fore.YELLOW}⚠️ Failed to parse RSS feed: {parse_error}{Style.RESET_ALL}"
                    )
                    # Save raw response for debugging
                    async with aiofiles.open(
                        os.path.join(export_path, "IdentityProductChanges_raw_error.txt"),
                        "w",
                        encoding="utf-8",
                    ) as f:
                        await f.write(rss_response)
                    print(f"{Fore.YELLOW}  ⚠️ Raw response saved for debugging{Style.RESET_ALL}")
            else:
                print(f"{Fore.YELLOW}⚠️ RSS feed returned non-string response{Style.RESET_ALL}")

            print(
                f"{Fore.YELLOW}ℹ️ Note: Data is sourced from Microsoft Entra RSS"
                f" feed{Style.RESET_ALL}"
            )

        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.YELLOW}⚠️ Identity Product Changes extraction will be skipped:"
                f" {str(e)}{Style.RESET_ALL}"
            )

    async def extract_office_cloud_policy_service(self, export_path: str) -> None:
        """Extract Office Cloud Policy Service (OCPS) Data"""
        print(
            f"{Fore.CYAN}🔍 Extracting Office Cloud Policy Service (OCPS) Data...{Style.RESET_ALL}"
        )

        ocps_path = os.path.join(export_path, "OCPS")
        if os.path.exists(ocps_path):
            shutil.rmtree(ocps_path)
        os.makedirs(ocps_path)

        ocps_endpoints = {
            "synchealth": (
                "https://clients.config.office.net/odbhealth/v1.0/synchealth/reports/versioncount"
            ),
            "languages": (
                "https://clients.config.office.net/releases/v1.0/FileList/"
                + "languagesForProductIds?productId=O365ProPlusRetail"
            ),
            "userflights": "https://config.office.com/appConfig/v1.0/userflights",
            "SettingsCatalog": (
                "https://clients.config.office.net/settings/v1.0/SettingsCatalog/Settings"
            ),
            "ServiceHealth": "https://config.office.com/appConfig/v1.0/ServiceHealth",
            "OfficeReleases": "https://clients.config.office.net/releases/v1.0/OfficeReleases",
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

        for _index, ((endpoint_name, _), result) in enumerate(
            zip(list(ocps_endpoints.items()) + [("FeatureData", "")], results)
        ):
            if isinstance(result, Exception):
                failed_endpoints.append(endpoint_name)
            else:
                successful_endpoints.append(endpoint_name)

        if successful_endpoints:
            print(
                f"{Fore.GREEN}✅ OCPS extraction completed for:"
                f" {', '.join(successful_endpoints)}{Style.RESET_ALL}"
            )

        if failed_endpoints:
            print(
                f"{Fore.YELLOW}⚠️ OCPS extraction failed for: {', '.join(failed_endpoints)}"
                f" (authentication or availability){Style.RESET_ALL}"
            )

    async def _extract_ocps_endpoint(self, endpoint_name: str, endpoint_url: str, ocps_path: str):
        """Extract single OCPS endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")

        try:
            response = await self.http.request(endpoint_url)
            await self.save_json(response, os.path.join(ocps_path, f"{endpoint_name}.json"))
            return True
        except GRAPH_EXPORT_ERRORS as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(
                    f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} -"
                    f" skipping{Style.RESET_ALL}"
                )
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}:"
                    f" {error_message}{Style.RESET_ALL}"
                )
            return False

    async def _extract_ocps_feature_data(self, ocps_path: str):
        """Extract OCPS Feature data"""
        print(f"{Fore.YELLOW}  📋 Processing FeatureData...{Style.RESET_ALL}")

        try:
            feature_response = await self.http.request(
                "https://clients.config.office.net/onboarding/odata/v1.0/FeatureData"
            )

            if isinstance(feature_response, dict) and "value" in feature_response:
                await self.save_json(
                    feature_response["value"], os.path.join(ocps_path, "FeatureData.json")
                )
            else:
                await self.save_json(feature_response, os.path.join(ocps_path, "FeatureData.json"))
            return True
        except GRAPH_EXPORT_ERRORS as e:
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
            "authz": {"url": "https://authsvc.teams.microsoft.com/v1.0/authz", "method": "POST"},
            "appsCatalog": {
                "url": "https://teams.microsoft.com/api/mt/part/au-01/beta/users/appsCatalog",
                "method": "GET",
            },
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
            print(
                f"{Fore.GREEN}✅ Teams Admin Center extraction completed for:"
                f" {', '.join(successful_endpoints)}{Style.RESET_ALL}"
            )

        if failed_endpoints:
            print(
                f"{Fore.YELLOW}⚠️ Teams Admin Center extraction failed for:"
                f" {', '.join(failed_endpoints)} (authentication or availability){Style.RESET_ALL}"
            )

    async def _extract_teams_endpoint(
        self, endpoint_name: str, endpoint_config: Dict, teams_path: str
    ):
        """Extract single Teams endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")

        try:
            if endpoint_config["method"] == "POST":
                response = await self.http.request(endpoint_config["url"], method="POST", body={})
            else:
                response = await self.http.request(endpoint_config["url"])

            # Remove sensitive token data if present
            if isinstance(response, dict) and "tokens" in response:
                del response["tokens"]

            await self.save_json(response, os.path.join(teams_path, f"{endpoint_name}.json"))
            return True
        except GRAPH_EXPORT_ERRORS as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(
                    f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} -"
                    f" skipping{Style.RESET_ALL}"
                )
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}:"
                    f" {error_message}{Style.RESET_ALL}"
                )
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
            "features-config": "https://admin.microsoft.com/admin/api/features/config",
            "features-all": "https://admin.microsoft.com/admin/api/features/all",
            "partner-list": (
                "https://admin.microsoft.com/fd/bcws/api/v1/IntraTenantPartner/getPartnerList"
            ),
            "product-offers": (
                "https://admin.microsoft.com/fd/bsxcommerce/v1/ProductOffers/"
                + "EligibleProductOffers?language=en-US"
            ),
        }

        # Complex endpoints
        complex_endpoints = {
            "apps": {
                "url": (
                    "https://admin.microsoft.com/fd/addins/api/availableApps?workloads=MetaOS,Teams"
                ),
                "property": "apps",
            },
            "policy_definitions": {
                "url": (
                    "https://admin.microsoft.com/fd/edgeenterpriseextensionsmanagement/api/policies"
                ),
                "property": "policy_definitions",
            },
            "C2RReleaseInfo": {
                "url": "https://admin.microsoft.com/fd/dms/odata/C2RReleaseInfo",
                "property": "value",
            },
            "ProductOfferIndex": {
                "url": (
                    "https://admin.microsoft.com/fd/bsxcommerce/v1/ProductOfferIndex?language=en-US"
                ),
                "property": "results",
            },
            "licensedProducts": {
                "url": "https://admin.microsoft.com/fd/m365licensing/v3/licensedProducts",
                "property": "value",
            },
            "sidebarExtensions": {
                "url": (
                    "https://admin.microsoft.com/fd/edgeenterpriseextensionsmanagement/"
                    + "api/sidebarExtensions"
                ),
                "property": "hub_apps",
            },
        }

        # Process all endpoints concurrently
        tasks = []

        # Simple endpoints
        for endpoint_name, endpoint_url in simple_endpoints.items():
            tasks.append(
                self._extract_m365_simple_endpoint(endpoint_name, endpoint_url, m365_admin_path)
            )

        # Complex endpoints
        for endpoint_name, endpoint_config in complex_endpoints.items():
            tasks.append(
                self._extract_m365_complex_endpoint(endpoint_name, endpoint_config, m365_admin_path)
            )

        # Special endpoints
        tasks.append(self._extract_m365_service_health(m365_admin_path))
        tasks.append(self._extract_m365_message_center(m365_admin_path))
        tasks.append(self._extract_m365_concierge_config(m365_admin_path))

        results = await asyncio.gather(*tasks, return_exceptions=True)

        successful_endpoints = []
        failed_endpoints = []

        all_endpoint_names = (
            list(simple_endpoints.keys())
            + list(complex_endpoints.keys())
            + ["ServiceHealth", "MessageCenter", "ConciergeConfig"]
        )

        for endpoint_name, result in zip(all_endpoint_names, results):
            if isinstance(result, Exception):
                failed_endpoints.append(endpoint_name)
            else:
                successful_endpoints.append(endpoint_name)

        if successful_endpoints:
            print(
                f"{Fore.GREEN}✅ M365 Admin Center extraction completed for:"
                f" {len(successful_endpoints)} endpoints{Style.RESET_ALL}"
            )

        if failed_endpoints:
            print(
                f"{Fore.YELLOW}⚠️ M365 Admin Center extraction failed for: {len(failed_endpoints)}"
                f" endpoints (authentication or availability){Style.RESET_ALL}"
            )

    async def _extract_m365_simple_endpoint(
        self, endpoint_name: str, endpoint_url: str, m365_admin_path: str
    ):
        """Extract simple M365 endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")

        try:
            response = await self.http.request(endpoint_url)
            await self.save_json(response, os.path.join(m365_admin_path, f"{endpoint_name}.json"))
            return True
        except GRAPH_EXPORT_ERRORS as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(
                    f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} -"
                    f" skipping{Style.RESET_ALL}"
                )
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}:"
                    f" {error_message}{Style.RESET_ALL}"
                )
            return False

    async def _extract_m365_complex_endpoint(
        self, endpoint_name: str, endpoint_config: Dict, m365_admin_path: str
    ):
        """Extract complex M365 endpoint"""
        print(f"{Fore.YELLOW}  📋 Processing {endpoint_name}...{Style.RESET_ALL}")

        try:
            response = await self.http.request(endpoint_config["url"])

            # Extract specific property if defined
            if (
                endpoint_config["property"]
                and isinstance(response, dict)
                and endpoint_config["property"] in response
            ):
                data_to_save = response[endpoint_config["property"]]
            else:
                data_to_save = response

            # Extract filename from URL
            url_parts = [
                part
                for part in endpoint_config["url"].split("/")
                if part
                and "https:" not in part
                and "admin.microsoft.com" not in part
                and "fd" not in part
                and "api" not in part
            ]
            file_name = url_parts[-1].split("?")[0] if url_parts else endpoint_name
            if not file_name:
                file_name = endpoint_name

            await self.save_json(data_to_save, os.path.join(m365_admin_path, f"{file_name}.json"))
            return True
        except GRAPH_EXPORT_ERRORS as e:
            error_message = str(e)
            if "401" in error_message or "403" in error_message or "Authorization" in error_message:
                print(
                    f"{Fore.YELLOW}    ⚠️ Authentication/authorization failed for {endpoint_name} -"
                    f" skipping{Style.RESET_ALL}"
                )
            else:
                print(
                    f"{Fore.YELLOW}    ⚠️ Failed to extract {endpoint_name}:"
                    f" {error_message}{Style.RESET_ALL}"
                )
            return False

    async def _extract_m365_service_health(self, m365_admin_path: str):
        """Extract M365 Service Health"""
        print(f"{Fore.YELLOW}  📋 Processing ServiceHealth...{Style.RESET_ALL}")

        try:
            service_health_response = await self.http.request(
                (
                    "https://admin.microsoft.com/admin/api/servicehealth/status/activeCM?"
                    + "showResolved=true"
                )
            )

            if (
                isinstance(service_health_response, dict)
                and "ServiceStatus" in service_health_response
            ):
                flattened_data = []
                for service in service_health_response["ServiceStatus"]:
                    if "MessagesByClassification" in service:
                        if "Incidents" in service["MessagesByClassification"]:
                            flattened_data.extend(service["MessagesByClassification"]["Incidents"])
                        if "Advisories" in service["MessagesByClassification"]:
                            flattened_data.extend(service["MessagesByClassification"]["Advisories"])
                await self.save_json(
                    flattened_data, os.path.join(m365_admin_path, "ServiceHealth.json")
                )
                return True
        except GRAPH_EXPORT_ERRORS as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract ServiceHealth: {str(e)}{Style.RESET_ALL}")
            return False

    async def _extract_m365_message_center(self, m365_admin_path: str):
        """Extract M365 Message Center"""
        print(f"{Fore.YELLOW}  📋 Processing MessageCenter...{Style.RESET_ALL}")

        try:
            message_center_response = await self.http.request(
                "https://admin.microsoft.com/admin/api/messagecenter"
            )

            if isinstance(message_center_response, dict) and "Messages" in message_center_response:
                # Remove sort-specific properties
                for message in message_center_response["Messages"]:
                    if "ActionRequiredBySortValue" in message:
                        del message["ActionRequiredBySortValue"]
                await self.save_json(
                    message_center_response["Messages"],
                    os.path.join(m365_admin_path, "messagecenter.json"),
                )
                return True
        except GRAPH_EXPORT_ERRORS as e:
            print(f"{Fore.YELLOW}    ⚠️ Failed to extract MessageCenter: {str(e)}{Style.RESET_ALL}")
            return False

    async def _extract_m365_concierge_config(self, m365_admin_path: str):
        """Extract M365 Concierge Config"""
        print(f"{Fore.YELLOW}  📋 Processing ConciergeConfig...{Style.RESET_ALL}")

        try:
            concierge_response = await self.http.request(
                "https://admin.microsoft.com/api/concierge/GetConciergeConfig"
            )

            # Remove SessionID for security
            if isinstance(concierge_response, dict) and "SessionID" in concierge_response:
                del concierge_response["SessionID"]

            await self.save_json(
                concierge_response, os.path.join(m365_admin_path, "GetConciergeConfig.json")
            )
            return True
        except GRAPH_EXPORT_ERRORS as e:
            print(
                f"{Fore.YELLOW}    ⚠️ Failed to extract ConciergeConfig: {str(e)}{Style.RESET_ALL}"
            )
            return False
