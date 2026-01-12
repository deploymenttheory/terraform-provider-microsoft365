"""Diff fetching and filtering."""

import json
import subprocess
import traceback
from typing import List, TYPE_CHECKING

if TYPE_CHECKING:
    from core.github_client import GitHubClient
    from core.progress_reporter import ProgressReporter


class DiffFetcher:
    """Handles fetching and filtering diffs from SDK repository."""

    def __init__(self, sdk_repo: str, github_client: 'GitHubClient', 
                 reporter: 'ProgressReporter', save_diff: bool = False):
        """Initialize diff fetcher.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
            github_client: GitHub client for API calls
            reporter: Progress reporter
            save_diff: If True, save diffs to files
        """
        self.sdk_repo = sdk_repo
        self.github_client = github_client
        self.reporter = reporter
        self.save_diff = save_diff

    def fetch_version_diff(self, old_version: str, new_version: str) -> str:
        """Fetch diff between two versions.
        
        Args:
            old_version: Old version tag
            new_version: New version tag
            
        Returns:
            Unified diff text for model files
        """
        try:
            self.reporter.info("Fetching diff from GitHub API...", indent=1)
            
            compare_data = self.github_client.get_compare_data(
                self.sdk_repo, old_version, new_version
            )
            
            files = compare_data.get('files', [])
            if not files:
                self.reporter.info("No files changed in this version", indent=1)
                return ""

            model_files = self._filter_model_files(files)
            if not model_files:
                self.reporter.info("No model files with changes found", indent=1)
                return ""

            diff_text = self._build_unified_diff(model_files)
            self.reporter.info(f"📦 Collected diffs from {len(model_files)} model file(s)", indent=1)

            if self.save_diff:
                self._save_diff_to_file(diff_text, old_version, new_version)

            return diff_text

        except (subprocess.CalledProcessError, json.JSONDecodeError, OSError, KeyError) as e:
            self.reporter.error(f"Error fetching diff: {e}")
            traceback.print_exc()
            return ""

    def _filter_model_files(self, files: List[dict]) -> List[dict]:
        """Filter for model files with patches.
        
        Args:
            files: List of file info from GitHub API
            
        Returns:
            Filtered list of model files
        """
        model_files = []
        
        for file_info in files:
            filename = file_info.get('filename', '')
            
            if not filename.startswith('models/'):
                continue

            patch = file_info.get('patch')
            status = file_info.get('status', 'unknown')

            if not patch:
                self.reporter.info(
                    f"⚠️  Skipping {filename} (no patch content, status: {status})",
                    indent=2
                )
                continue

            model_files.append(file_info)
            self.reporter.info(f"✓ Found changes in {filename} ({status})", indent=2)

        return model_files

    def _build_unified_diff(self, files: List[dict]) -> str:
        """Build unified diff format from file patches.
        
        Args:
            files: List of file info with patches
            
        Returns:
            Unified diff text
        """
        diff_parts = []
        
        for file_info in files:
            filename = file_info['filename']
            patch = file_info['patch']
            
            diff_parts.append(f"diff --git a/{filename} b/{filename}")
            diff_parts.append(f"--- a/{filename}")
            diff_parts.append(f"+++ b/{filename}")
            diff_parts.append(patch)
            diff_parts.append("")

        return '\n'.join(diff_parts)

    def _save_diff_to_file(self, diff_text: str, old_version: str, new_version: str):
        """Save diff to file for debugging.
        
        Args:
            diff_text: Diff content
            old_version: Old version
            new_version: New version
        """
        diff_file = f"sdk_diff_{old_version}_to_{new_version}.patch"
        with open(diff_file, 'w', encoding='utf-8') as f:
            f.write(diff_text)
        self.reporter.info(f"💾 Saved diff to {diff_file}", indent=1)

    def fetch_changelog_section(self, version: str) -> str:
        """Fetch changelog section for a specific version.
        
        Args:
            version: Version to fetch changelog for
            
        Returns:
            Changelog section text
        """
        try:
            changelog_url = f"https://raw.githubusercontent.com/{self.sdk_repo}/main/CHANGELOG.md"
            changelog_content = self.github_client.fetch_changelog(changelog_url)
            
            return self._extract_version_section(changelog_content, version)

        except subprocess.CalledProcessError as e:
            self.reporter.error(f"Error fetching changelog: {e}")
            return "Error fetching changelog"

    def _extract_version_section(self, changelog: str, version: str) -> str:
        """Extract version section from changelog.
        
        Args:
            changelog: Full changelog content
            version: Version to extract
            
        Returns:
            Version section text
        """
        lines = changelog.split('\n')
        section_lines = []
        in_section = False

        for line in lines:
            if line.startswith('##') and version.lstrip('v') in line:
                in_section = True
                section_lines.append(line)
            elif in_section:
                if line.startswith('##'):
                    break
                section_lines.append(line)

        return '\n'.join(section_lines) if section_lines else "Changelog section not found"
