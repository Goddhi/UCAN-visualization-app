# UCAN Visualization Web App

UCAN Visualization Web App is a comprehensive developer tool for parsing, validating, and visualizing UCAN (User Controlled Authorization Networks) delegation chains. This application addresses the core challenge of debugging decentralized authorization by providing visual feedback and detailed error reporting for UCAN tokens.

### Features

**Delegation Chain Visualizer**

Visual tree/graph showing who delegated what to whom
Interactive exploration of trust relationships
Color-coded validity indicators (valid/expired/invalid)
See the complete authorization path from root to end user

**Invocation Inspector**

Paste a UCAN invocation and see exactly what it's trying to do
View the capability being invoked (resource, ability, caveats)
Inspect all attached proofs in the invocation
Understand invocation metadata (issuer, audience, timestamps)

**Capability Breakdown**

Parse and display resources, abilities, and attenuations in human-readable format
See technical capability structure translated to plain English
Understand what permissions are granted and what restrictions apply
Identify how capabilities were narrowed through delegation (attenuation)

**Proof Chain Validator**

Check if a UCAN chain is valid and where it breaks if invalid
Cryptographic signature verification
Time-based validation (expiration, not-before)
Capability escalation detection
Clear root cause analysis with fix suggestions