#!/bin/bash

if [ $# -ne 4 ]; then
  echo "Usage: `basename $0` {database} {pdf_file} {invoice_name} {page_range_i}"
  exit 65
fi

DATABASE=$1
PDF_FILE=$2
INVOICE_NAME=$3
PAGE_RANGE=$4
TEMP_PDF_I="temp_invoice_$INVOICE_NAME.pdf"

# Make sure INVOICE_NAME does not already exist in the database
if [ -n "$(sqlite3 $DATABASE "SELECT name FROM invoices WHERE name = '$INVOICE_NAME';")" ]; then
  echo "Invoice name $INVOICE_NAME already exists in the database."
  exit 1
fi

# Extract the specified page range from the itinerary PDF file
qpdf "$PDF_FILE" --pages . $PAGE_RANGE -- "$TEMP_PDF_I"

if [ ! -f "$TEMP_PDF_I" ]; then
    echo "Page extraction failed."
    exit 1
else
    sqlite3 $DATABASE "INSERT INTO invoices (name, pdf) VALUES ('$INVOICE_NAME', readfile('$TEMP_PDF_I'));"
fi

# Clean up the temporary PDF files
rm "$TEMP_PDF_I"
