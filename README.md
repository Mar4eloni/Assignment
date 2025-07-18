# Email Address Parser
 Go implementation of an RFC 5322 compliant email address parser that extracts and validates email addresses from text input, including handling display names and comments.
 
## Approach

### High-Level Design

1. Two-Phase Parsing :
   - First extracts potential email candidates using heuristics
   - then validates them agains the provided RFC 5322 regex pattern

2. Component Separation:
   - Cleans input by removing marks and normalizing whitespace
   - Handles displan names separatley from email addresses
   - Processes comments and quoted strings

3. Validation Pipline:
   - Input → Cleaning → Component Extraction → Validation → JSON Output

## Key Insights

### Core Solutions

1. Comment Handling:
   - Uses state tracking with 'commentDepth' and 'inQuote' to properly ignore nested comments while preserving quoted text
   - - Implements 'removeComments()' to strip only unquoted parentheses

2. Email Extraction:
   - Prioritezes angle-bracket notation but falls back to regex matching
   - Takes the last valid match when multiple candidates exist

3. Display Name Processing:
   - Handles both quoted '"Name Surname"' and unquoted formats
   - Properly unescapes special sequences like '\"' and '\\'
   - Preserves parentheses within quoted display names

## Weaknesses and Limitations

### Design Limitations

1. Internationalization:
   - Doesn't fully support international email adresses
   - limited handling of utf-8 in the names

2. Edge Cases:
   - Nat struggle with some obscure RFC 5322 formats
   - Incomplete handling of nested quoted strings
3. Error Recovery:
   - Fails on first error

## Other Comments

### Tradeoffs 
Choose a strict compliance over leniency

### Future Improvements:
Implementation of internatiponalization
Precompiled reges patterns
Single-pass parsing design
Support for multiple output formats like XML,CSV ...

## Usage

go build -o email-addr
./email-addr testdata/emails.txt or ./email-addr <file-name>
