#!/usr/bin/env python3
#
# Copyright (C) 2015 Canonical Ltd.
# Author: Xavi Garcia <xavi.garcia.mena@canonical.com>
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 3.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#

"""
This is a test of a go-bindings scope.
"""

from scope_harness import ( ScopeHarness, CategoryMatcher, CategoryMatcherMode, CategoryListMatcher,
                            CategoryListMatcherMode, ResultMatcher, PreviewMatcher, PreviewWidgetMatcher,
                            PreviewColumnMatcher, PreviewView,
                            Parameters, DepartmentMatcher, ChildDepartmentMatcher )
from scope_harness.testing import ScopeHarnessTestCase
import unittest
import sys
import os
from shutil import copyfile
import inspect

TEST_DATA_DIR = os.path.dirname(os.path.abspath(inspect.getfile(inspect.currentframe())))
SCOPE_NAME = "goscope"
CONFIG_FILE = "goscope.ini"

def remove_ini_file():
    os.remove(TEST_DATA_DIR + "/" + CONFIG_FILE)

def copy_binary():
    binary = os.environ.get('GOPATH') + "/bin/" + SCOPE_NAME
    if not os.path.exists(binary):
        raise Exception("The binary %s does not exist." % binary)
    copyfile(binary, TEST_DATA_DIR + "/" + SCOPE_NAME)

def remove_binary():
    os.remove(TEST_DATA_DIR + "/" + SCOPE_NAME)

def prepare_ini_file():
    copy_and_set_go_path(TEST_DATA_DIR + "/" + CONFIG_FILE + ".in", TEST_DATA_DIR + "/" + CONFIG_FILE)

def copy_and_set_go_path(infile, outfile):
    with open(outfile, "wt") as fout:
        with open(infile, "rt") as fin:
            for line in fin:
                fout.write(line.replace("$GOPATH$", os.environ.get('GOPATH')))

class ResultsTest(ScopeHarnessTestCase):
    @classmethod
    def setUpClass(cls):
        copy_binary()
        prepare_ini_file()
        cls.harness = ScopeHarness.new_from_scope_list(Parameters([
            TEST_DATA_DIR + "/goscope.ini"
            ]))
        cls.view = cls.harness.results_view
        cls.view.active_scope = "goscope"
        cls.view.search_query = ""

    @classmethod
    def tearDownClass(cls):
        remove_ini_file()
        remove_binary()

    def test_basic_result(self):
        self.assertMatchResult(CategoryListMatcher()
            .has_at_least(1)
            .mode(CategoryListMatcherMode.BY_ID)
            .category(CategoryMatcher("category")
                    .has_at_least(1)
                    .mode(CategoryMatcherMode.BY_URI)
                    .title("Category")
                    .icon("")
                    )
            .match(self.view.categories)
            )

    def test_scope_properties(self):
        self.assertEqual(self.view.scope_id, 'goscope')
        self.assertEqual(self.view.display_name, 'mock.DisplayName')
        self.assertEqual(self.view.icon_hint, '/mock.Icon')
        self.assertEqual(self.view.description, 'mock.Description')
        self.assertEqual(self.view.search_hint, 'mock.SearchHint')
        self.assertEqual(self.view.shortcut, 'mock.HotKey')

        customizations = self.view.customizations
        self.assertTrue(len(customizations)!=0)

        header_customizations = customizations["page-header"]
        self.assertEqual(header_customizations["logo"], "http://assets.ubuntu.com/sites/ubuntu/1110/u/img/logos/logo-ubuntu-orange.svg")
        self.assertEqual(header_customizations["background"], "color://black")
        self.assertEqual(header_customizations["foreground-color"], "white")
        self.assertEqual(customizations["shape-images"], False)

    def test_result_data(self):
        self.assertMatchResult(CategoryListMatcher()
            .has_at_least(1)
            .mode(CategoryListMatcherMode.BY_ID)
            .category(CategoryMatcher("category")
                    .has_at_least(1)
                    .mode(CategoryMatcherMode.BY_URI)
                    .result(ResultMatcher("http://localhost/")
                    .properties({'test_value_bool': True})
                    .properties({'test_value_string': "test_value"})
                    .properties({'test_value_int': 1999})
                    .properties({'test_value_float': 1.999})
                    .dnd_uri("http://localhost_dnduri")
                    .art("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
                    ))
            .match(self.view.categories)
        )

        self.assertMatchResult(CategoryListMatcher()
            .has_at_least(1)
            .mode(CategoryListMatcherMode.BY_ID)
            .category(CategoryMatcher("category")
                    .has_at_least(1)
                    .mode(CategoryMatcherMode.BY_URI)
                    .result(ResultMatcher("http://localhost2/")
                    .properties({'test_value_bool': False})
                    .properties({'test_value_string': "test_value2"})
                    .properties({'test_value_int': 2000})
                    .properties({'test_value_float': 2.1})
                    .dnd_uri("http://localhost_dnduri2")
                    .properties({'test_value_map': {'value1':1,'value2':'string_value'}})
                    .properties({'test_value_array': [1999,"string_value"]})
                    .art("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
                    ))
            .match(self.view.categories)
        )

    def test_result_data_query(self):
        self.view.active_scope = 'goscope'
        test_query = "test_query"
        self.view.search_query = test_query

        self.assertMatchResult(CategoryListMatcher()
            .has_at_least(1)
            .mode(CategoryListMatcherMode.BY_ID)
            .category(CategoryMatcher("category")
                    .has_at_least(1)
                    .mode(CategoryMatcherMode.BY_URI)
                    .result(ResultMatcher("http://localhost/" + test_query)
                    .properties({'test_value_bool': True})
                    .properties({'test_value_string': "test_value" + test_query})
                    .properties({'test_value_int': 1999})
                    .properties({'test_value_float': 1.999})
                    .dnd_uri("http://localhost_dnduri" + test_query)
                    .art("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
                    ))
            .match(self.view.categories)
        )

        self.assertMatchResult(CategoryListMatcher()
            .has_at_least(1)
            .mode(CategoryListMatcherMode.BY_ID)
            .category(CategoryMatcher("category")
                    .has_at_least(1)
                    .mode(CategoryMatcherMode.BY_URI)
                    .result(ResultMatcher("http://localhost2/" + test_query)
                    .properties({'test_value_bool': False})
                    .properties({'test_value_string': "test_value2" + test_query})
                    .properties({'test_value_int': 2000})
                    .properties({'test_value_float': 2.1})
                    .dnd_uri("http://localhost_dnduri2" + test_query)
                    .properties({'test_value_map': {'value1':1,'value2':'string_value'}})
                    .properties({'test_value_array': [1999,"string_value"]})
                    .art("https://pbs.twimg.com/profile_images/1117820653/5ttls5.jpg.png")
                    ))
            .match(self.view.categories)
        )

class DepartmentsTest(ScopeHarnessTestCase):
    @classmethod
    def setUpClass(cls):
        copy_binary()
        prepare_ini_file()
        cls.harness = ScopeHarness.new_from_scope_list(Parameters([
                    TEST_DATA_DIR + "/goscope.ini"
            ]))
        cls.view = cls.harness.results_view
        cls.view.active_scope = "goscope"
        cls.view.search_query = ""

    @classmethod
    def tearDownClass(cls):
        remove_ini_file()
        remove_binary()

    def test_basic_departments(self):
        self.view.active_scope = 'goscope'
        self.view.search_query = ''

        departments = self.view.browse_department('')
        self.assertTrue(self.view.has_departments)
        self.assertEqual(len(departments), 2)
        # excercise different methods for getting children
        dep = departments[0]
        dep2 = departments.child(0)
        dep3 = departments.children[0]
        self.assertEqual(dep.id, dep2.id)
        self.assertEqual(dep2.id, dep3.id)

        self.assertMatchResult(DepartmentMatcher()
            .has_exactly(2)
            .label('Browse Music')
            .all_label('Browse Music Alt')
            .parent_id('')
            .parent_label('')
            .is_root(True)
            .is_hidden(False)
            .child(ChildDepartmentMatcher('Rock')
                   .label('Rock Music')
                   .has_children(True)
                   .is_active(False)
                   )
            .child(ChildDepartmentMatcher('Soul')
                   .label('Soul Music')
                   .has_children(True)
                   .is_active(False)
                   )
            .match(departments)
        )
    def test_child_department(self):
        self.view.active_scope = 'goscope'
        departments = self.view.browse_department('Rock')
        self.assertEqual(len(departments), 2)
        # excercise different methods for getting children
        dep = departments[0]
        dep2 = departments.child(0)
        dep3 = departments.children[0]
        self.assertEqual(dep.id, dep2.id)
        self.assertEqual(dep2.id, dep3.id)

        self.assertMatchResult(DepartmentMatcher()
            .has_exactly(2)
            .label('Rock Music')
            .all_label('Rock Music Alt')
            .parent_id('')
            .parent_label('Browse Music')
            .is_root(False)
            .is_hidden(False)
            .child(ChildDepartmentMatcher('60s')
                   .label('Rock from the 60s')
                   .has_children(False)
                   .is_active(False)
                   )
            .child(ChildDepartmentMatcher('70s')
                   .label('Rock from the 70s')
                   .has_children(False)
                   .is_active(False)
                   )
            .match(departments)
        )

class PreviewTest(ScopeHarnessTestCase):
    @classmethod
    def setUpClass(cls):
        copy_binary()
        prepare_ini_file()
        cls.harness = ScopeHarness.new_from_scope_list(Parameters([
                    TEST_DATA_DIR + "/goscope.ini"
            ]))
        cls.view = cls.harness.results_view
        cls.view.active_scope = "goscope"
        cls.view.search_query = ""


    def test_preview_layouts(self):
        pview = self.view.category(0).result(0).tap()
        self.assertIsInstance(pview, PreviewView)

        pview.column_count = 3
        self.assertMatchResult(PreviewColumnMatcher()
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("image")))
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("header"))
                         .widget(PreviewWidgetMatcher("summary"))
                         .widget(PreviewWidgetMatcher("actions")))
                 .column(PreviewMatcher()
                        ).match(pview.widgets))

        pview.column_count = 2
        self.assertMatchResult(PreviewColumnMatcher()
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("image")))
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("header"))
                         .widget(PreviewWidgetMatcher("summary"))
                         .widget(PreviewWidgetMatcher("actions"))
                        ).match(pview.widgets))

        pview.column_count = 1
        self.assertMatchResult(PreviewColumnMatcher()
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("image"))
                         .widget(PreviewWidgetMatcher("header"))
                         .widget(PreviewWidgetMatcher("summary"))
                         .widget(PreviewWidgetMatcher("actions"))
                        ).match(pview.widgets))

    def test_preview_action(self):
        pview = self.view.category(0).result(0).tap()
        self.assertIsInstance(pview, PreviewView)

        pview.column_count = 1
        self.assertMatchResult(PreviewColumnMatcher()
                 .column(PreviewMatcher()
                         .widget(PreviewWidgetMatcher("image"))
                         .widget(PreviewWidgetMatcher("header"))
                         .widget(PreviewWidgetMatcher("summary"))
                         .widget(PreviewWidgetMatcher("actions"))
                        ).match(pview.widgets))

        next_view = pview.widgets_in_first_column["actions"].trigger("hide", None)
        self.assertEqual(pview, next_view)

if __name__ == '__main__':
    unittest.main(argv = sys.argv[:1])
