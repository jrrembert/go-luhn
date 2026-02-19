#!/bin/bash

# Updates copyright year ranges in LICENSE and README.md.
# Called by the update-copyright-date GitHub Actions workflow.

current_year=$(date +"%Y")

files=("LICENSE" "README.md")
updated=0

for file in "${files[@]}"; do
    if [ ! -f "$file" ]; then
        echo "Skipping $file (not found)"
        continue
    fi

    # Update "Copyright (c) YYYY" or "Copyright (c) YYYY-YYYY" (LICENSE)
    # Update "© YYYY" or "© YYYY-YYYY" (README.md)
    sed -i.bak -E \
        -e "s/(Copyright \(c\) )([0-9]{4})([-][0-9]{4})?/\1\2-$current_year/" \
        -e "s/(© )([0-9]{4})([-][0-9]{4})?/\1\2-$current_year/" \
        "$file"
    rm -f "${file}.bak"

    if git diff --quiet "$file" 2>/dev/null; then
        echo "$file: already up to date"
    else
        echo "$file: updated copyright year to $current_year"
        updated=1
    fi
done

if [ "$updated" -eq 0 ]; then
    echo "No changes needed"
fi
