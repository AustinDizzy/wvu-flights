#!/bin/bash

if [ $# -ne 5 ]; then
  echo "Usage: `basename $0` {database} {pdf_file} {target_id} {page_range_i} {page_range_j}"
  exit 65
fi

DATABASE=$1
TARGET_ID=$3
PDF_FILE=$2
PAGE_RANGE_I=$4
PAGE_RANGE_J=$5
TEMP_PDF_I="temp_itinerary_$TARGET_ID.pdf"
TEMP_PDF_R="temp_reservation_$TARGET_ID.pdf"

# Check if $TARGET_ID exists in the database
if [ -z "$(sqlite3 $DATABASE "SELECT id FROM trips WHERE id = '$TARGET_ID';")" ]; then
  echo "Trip ID $TARGET_ID does not exist in the database."
  exit 1
fi

# Extract the specified page range from the itinerary PDF file
qpdf "$PDF_FILE" --pages . $PAGE_RANGE_I -- "$TEMP_PDF_I"

# Extract the specified page range from the plane reservation PDF file
qpdf "$PDF_FILE" --pages . $PAGE_RANGE_J -- "$TEMP_PDF_R"

if [ ! -f "$TEMP_PDF_I" ]; then
    echo "Page extraction failed."
    exit 1
else
  sqlite3 $DATABASE "UPDATE trips SET itinerary = readfile('$TEMP_PDF_I') WHERE id = '$TARGET_ID';"
fi

if [ ! -f "$TEMP_PDF_R" ]; then
    echo "Page extraction failed."
    exit 1
else
  sqlite3 $DATABASE "UPDATE trips SET reservation = readfile('$TEMP_PDF_R') WHERE id = '$TARGET_ID';"
fi

# Clean up the temporary PDF files
rm "$TEMP_PDF_I"
rm "$TEMP_PDF_R"
