# AttendeeAttestor - Agent Architecture

This document describes the agentic architecture and development guidelines for the AttendeeAttestor system.

## System Overview

AttendeeAttestor is a privacy-preserving QR ticket verification system composed of three main apps;

1. **Issuer Service** (`issuer-service`) - Generates and cryptographically signs QR codes for event tickets
2. **Android Scanner** (`android`) - Verifies QR codes offline, for door staff, designed for Andoid.
3. **iOS Scanner** (`ios`) - Verifies QR codes offline, for door staff, designed for iOS/iPhone.

## App Responsibilities

### Issuer Service

**Primary Functions:**
- Generate unique ticket identifiers
- Create ticket payloads containing event ID, ticket ID, time window, and ticket type
- Cryptographically sign tickets using private key
- Encode signed tickets as QR codes
- Export public keys for scanner distribution

**Key Constraints:**
- Must never embed personal information (names, emails, phone numbers)
- Must generate cryptographically secure random ticket IDs
- Must support time-windowed validity periods
- Must produce signatures verifiable by corresponding public keys

### Android and iOS Scanner apps

**Primary Functions:**
- Scan QR codes from tickets
- Verify cryptographic signatures offline using public keys
- Validate ticket structure and required fields
- Check event ID matches expected event
- Validate ticket is within valid time window
- Maintain local cache of used ticket IDs for replay protection
- Optionally sync used ticket IDs when connectivity available
- Display clear admit/reject status to staff

**Key Constraints:**
- Must operate fully offline for core verification
- Must prevent double-entry using local cache
- Must provide simple, clear UI for non-technical staff
- Must handle network failures gracefully
- Must verify signatures without contacting central servers

## Development Guidelines

### Test-Driven Development (TDD)

All agents must follow strict test-driven development practices:

1. **Write tests first** - Define expected behavior through tests before implementation
2. **Red-Green-Refactor cycle** - Write failing test, implement minimum code to pass, refactor
3. **Test coverage** - Maintain high test coverage for all critical paths
4. **Integration tests** - Verify agent interactions and end-to-end workflows
5. **Property-based tests** - Use property-based testing for cryptographic operations

**Test Structure:**
- Unit tests for individual functions
- Integration tests for agent interactions
- Cryptographic tests for signature verification
- Offline operation tests for scanner agent
- Privacy tests to ensure no personal data leakage

### Cyclomatic Complexity

**Maximum complexity: 4**

All functions must maintain cyclomatic complexity below 4. This ensures:
- Functions remain simple and testable
- Code is easier to understand and maintain
- Edge cases are more manageable
- Refactoring is straightforward

**Strategies to reduce complexity:**
- Extract complex conditionals into separate functions
- Use early returns to reduce nesting
- Break large functions into smaller, focused functions
- Use strategy patterns for conditional logic
- Leverage helper functions for repeated patterns

**Example of acceptable complexity:**
```go
func validateTicketTimeWindow(ticket Ticket, currentTime time.Time) error {
    if ticket.ValidFrom.After(currentTime) {
        return ErrTicketNotYetValid
    }
    if ticket.ValidUntil.Before(currentTime) {
        return ErrTicketExpired
    }
    return nil
}
```

### Code Style

**Naming Conventions:**
- Use descriptive function names that clearly indicate purpose
- Use descriptive variable names that explain intent
- Avoid abbreviations unless widely understood
- Prefer clarity over brevity

**Function Naming:**
- Use verb-noun patterns: `generateTicketID()`, `verifySignature()`, `parseQRCode()`
- Boolean functions use `is`, `has`, `can` prefixes: `isTicketValid()`, `hasBeenUsed()`
- Error handling functions: `validateTicket()`, `checkTimeWindow()`

**Variable Naming:**
- Use nouns that describe the data: `ticketID`, `eventID`, `signature`, `publicKey`
- Avoid single-letter variables except in loops
- Use domain-specific terminology consistently

**No Code Comments:**
- Code must be self-documenting through naming
- Complex logic should be broken into smaller, well-named functions
- If code requires explanation, refactor to make intent clear
- Documentation comments are acceptable for exported APIs only

## Agent Interactions

### Ticket Issuance Flow

1. Issuer Agent generates ticket data structure
2. Issuer Agent signs ticket with private key
3. Issuer Agent encodes signed ticket as QR code
4. QR code distributed to attendee

### Ticket Verification Flow

1. Scanner Agent scans QR code
2. Scanner Agent decodes QR code to ticket data
3. Core Agent parses ticket structure
4. Core Agent verifies cryptographic signature
5. Scanner Agent validates event ID matches
6. Scanner Agent checks time window validity
7. Scanner Agent checks local cache for replay protection
8. Scanner Agent displays admit/reject status

### Replay Protection Sync Flow

1. Scanner Agent maintains local cache of used ticket IDs
2. When connectivity available, Scanner Agent syncs used IDs
3. Multiple scanners can share used ticket ID database
4. Sync prevents reuse across multiple entry points

## Privacy Requirements

All agents must adhere to strict privacy principles:

- **No personal data collection** - Never request or store names, emails, phone numbers
- **No tracking** - Do not create profiles or track attendee behavior
- **No accounts** - Attendees do not need accounts or registration
- **Cryptographic proof only** - Verify validity, not identity
- **Local-first** - Minimize network dependencies

## Security Requirements

- Use cryptographically secure random number generation for ticket IDs
- Implement proper key management for signing keys
- Validate all inputs before processing
- Handle cryptographic failures securely
- Protect against timing attacks in signature verification
- Ensure replay protection cannot be bypassed

## Error Handling

- All functions must return explicit error types
- Errors must be descriptive and actionable
- Cryptographic failures must be clearly distinguished from validation failures
- Network errors must not prevent offline verification
- Invalid input must be rejected with clear error messages

## Performance Requirements

- QR code generation must be fast enough for batch processing
- QR code scanning must provide near-instant feedback
- Signature verification must complete in milliseconds
- Local cache lookups must be efficient
- Offline operation must not degrade performance
