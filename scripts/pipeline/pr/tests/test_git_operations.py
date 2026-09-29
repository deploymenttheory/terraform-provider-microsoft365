"""Regression tests for exact changed-package selection using real Git diffs."""

from pathlib import Path

from support import RepositoryTest
from lib.git_operations import get_changed_files, get_changed_packages, get_merge_base


class ChangedPackagesTests(RepositoryTest):
    def test_no_changes(self):
        self.assertEqual(get_changed_packages(self.base), [])

    def test_non_go_changes(self):
        for name in ['README.md', 'example.tf', 'go.sum', '.github/workflows/lint.yml']:
            self.write(name, 'changed\n')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), [])

    def test_root_file(self):
        self.write('main.go', 'package main\nfunc main() {}\n')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), ['.'])

    def test_test_only_change(self):
        self.write('selected/value_test.go', 'package selected\n')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), ['selected'])

    def test_exact_packages_without_children(self):
        self.write('selected/value.go', 'package selected\nfunc Value() int { return 2 }\n')
        self.write('selected/child/value.go', 'package child\n')
        self.write('another/value.go', 'package another\n')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), ['another', 'selected', 'selected/child'])

    def test_deleted_file_with_surviving_package(self):
        self.write('selected/other.go', 'package selected\n')
        self.commit()
        base = self.git('rev-parse', 'HEAD')
        Path('selected/value.go').unlink()
        self.commit()
        self.assertEqual(get_changed_packages(base), ['selected'])

    def test_deleted_package(self):
        Path('selected/value.go').unlink()
        self.commit()
        self.assertEqual(get_changed_packages(self.base), [])

    def test_rename_includes_both_surviving_packages(self):
        self.write('selected/other.go', 'package selected\n')
        self.commit()
        base = self.git('rev-parse', 'HEAD')
        self.git('mv', 'selected/value.go', 'unrelated/moved.go')
        self.commit()
        self.assertEqual(get_changed_packages(base), ['selected', 'unrelated'])

    def test_renamed_entire_package(self):
        self.git('mv', 'selected', 'renamed')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), ['renamed'])

    def test_go_file_renamed_to_non_go(self):
        self.git('mv', 'selected/value.go', 'selected/value.txt')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), [])

    def test_more_than_300_files(self):
        for index in range(305):
            self.write(f'docs/{index}.md', 'example\n')
        self.write('selected/value_test.go', 'package selected\n')
        self.commit()
        self.assertEqual(get_changed_packages(self.base), ['selected'])

    def test_null_delimited_filenames(self):
        self.write('selected/with space.go', 'package selected\n')
        self.commit()
        self.assertEqual(get_changed_files(self.base), ['selected/with space.go'])

    def test_merge_base_excludes_base_branch_only_changes(self):
        self.git('checkout', '-qb', 'base-branch')
        self.write('unrelated/base.go', 'package unrelated\n')
        self.commit()
        self.git('checkout', '-qb', 'feature', self.base)
        self.write('selected/new.go', 'package selected\n')
        self.commit()
        self.assertEqual(get_merge_base('base-branch'), self.base)
        self.assertEqual(get_changed_packages('base-branch'), ['selected'])
