#!/usr/bin/env python3
"""
Quick diagnostic script to examine the RSS feed structure and categories.
"""

import sys
from collections import Counter

try:
    import feedparser
except ImportError:
    print("Error: feedparser not installed. Run: pip install feedparser")
    sys.exit(1)

RSS_FEED_URL = "https://developer.microsoft.com/en-us/graph/changelog/rss"

print(f"Fetching RSS feed from {RSS_FEED_URL}...\n")

feed = feedparser.parse(RSS_FEED_URL)

print(f"Total entries: {len(feed.entries)}")
print(f"\n{'='*80}")
print("ANALYZING RSS FEED STRUCTURE")
print(f"{'='*80}\n")

# Look at first few entries
print("First 3 entries structure:")
for i, entry in enumerate(feed.entries[:3], 1):
    print(f"\n--- Entry {i} ---")
    print(f"Title: {entry.get('title', 'N/A')}")
    print(f"ID/GUID: {entry.get('id', 'N/A')}")
    print(f"Published: {entry.get('published', 'N/A')}")
    
    # Check for tags
    if hasattr(entry, 'tags'):
        print(f"Tags: {entry.tags}")
        for tag in entry.tags:
            print(f"  - term: {tag.term}, scheme: {tag.get('scheme', 'N/A')}, label: {tag.get('label', 'N/A')}")
    else:
        print("Tags: None")
    
    # Check for categories
    if hasattr(entry, 'category'):
        print(f"Category: {entry.category}")
    
    print(f"Description (first 200 chars): {entry.get('description', 'N/A')[:200]}")

# Collect all categories/tags
all_categories = []
all_terms = []

for entry in feed.entries:
    if hasattr(entry, 'tags'):
        for tag in entry.tags:
            if hasattr(tag, 'term'):
                all_terms.append(tag.term)
                # Separate API versions from categories
                if tag.term not in ['v1.0', 'beta', 'Prod', 'prd']:
                    all_categories.append(tag.term)

print(f"\n{'='*80}")
print("CATEGORY ANALYSIS")
print(f"{'='*80}\n")

print(f"Total unique terms: {len(set(all_terms))}")
print(f"Total unique categories (excluding API versions): {len(set(all_categories))}")

print("\nAll unique categories:")
category_counts = Counter(all_categories)
for category, count in sorted(category_counts.items(), key=lambda x: x[1], reverse=True):
    print(f"  - {category}: {count} occurrences")

print("\nAll unique terms (including API versions):")
term_counts = Counter(all_terms)
for term, count in sorted(term_counts.items(), key=lambda x: x[1], reverse=True):
    print(f"  - {term}: {count} occurrences")

