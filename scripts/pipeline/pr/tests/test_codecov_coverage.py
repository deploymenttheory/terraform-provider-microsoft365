"""Exercise public/private Codecov access through the actual coverage gate."""

import io
import json
import os
import sys
import unittest
from unittest.mock import patch
from urllib.error import HTTPError

from support import PIPELINE  # Adds pipeline modules to sys.path.
from enforce_coverage import main


class CodecovCoverageTests(unittest.TestCase):
    def run_gate(self, response, token=None):
        argv = ['enforce_coverage.py', '--repo-slug', 'example/provider', '--pr-number', '123']
        environment = {} if token is None else {'CODECOV_API_TOKEN': token}
        requests = []

        def respond(request, **kwargs):
            requests.append(request)
            if isinstance(response, Exception):
                raise response
            return io.BytesIO(json.dumps(response).encode())

        with patch.object(sys, 'argv', argv), patch.dict(os.environ, environment, clear=True):
            with patch('enforce_coverage.load_pr_checks_config', return_value={'coverage_threshold': {'minimum_pct': 60}}):
                with patch('lib.get_codecov_coverage.urllib.request.urlopen', side_effect=respond):
                    with patch('lib.get_codecov_coverage.time.sleep'), patch('sys.stdout', new_callable=io.StringIO):
                        result = main()
        return result, requests

    def test_public_coverage_passes_without_authorization_header(self):
        for token in [None, '']:
            with self.subTest(token=token):
                result, requests = self.run_gate({'patch': {'coverage': 71.9, 'hits': 87, 'misses': 30, 'partials': 4}}, token)
                self.assertEqual(result, 0)
                self.assertEqual(len(requests), 1)
                self.assertNotIn('Authorization', dict(requests[0].header_items()))
                self.assertEqual(requests[0].get_header('Accept'), 'application/json')

    def test_private_coverage_keeps_supplied_api_token(self):
        result, requests = self.run_gate({'patch': {'coverage': 80, 'hits': 8, 'misses': 2}}, 'test-api-token')
        self.assertEqual(result, 0)
        self.assertEqual(requests[0].get_header('Authorization'), 'bearer test-api-token')

    def test_below_threshold_fails_without_token(self):
        for coverage in [0, 59.99]:
            with self.subTest(coverage=coverage):
                result, _ = self.run_gate({'patch': {'coverage': coverage, 'hits': 0, 'misses': 10}})
                self.assertEqual(result, 1)

    def test_exact_threshold_passes_without_token(self):
        result, _ = self.run_gate({'patch': {'coverage': 60, 'hits': 6, 'misses': 4}})
        self.assertEqual(result, 0)

    def test_missing_coverage_fails_instead_of_skipping_the_gate(self):
        for response in [{}, {'patch': None}, {'patch': {}}]:
            with self.subTest(response=response):
                result, requests = self.run_gate(response)
                self.assertEqual(result, 1)
                self.assertGreater(len(requests), 1)

    def test_access_denied_fails_without_retrying_anonymously(self):
        for code in [401, 403]:
            for token in [None, 'test-api-token']:
                with self.subTest(code=code, token=token):
                    error = HTTPError('https://api.codecov.io/example', code, 'denied', {}, None)
                    self.addCleanup(error.close)
                    result, requests = self.run_gate(error, token)
                    self.assertEqual(result, 1)
                    self.assertEqual(len(requests), 1)


if __name__ == '__main__':
    unittest.main()
