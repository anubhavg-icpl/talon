# stage: recon
# category: core

> Information gathering workflow — passive + active reconnaissance

# Information Gathering Skill

Perform passive and active information gathering to build a target profile and an attack surface map.

## Execution Steps

### 1. Passive Reconnaissance
- Access the target via the fetch tool and collect HTTP response headers
- Identify the server type, version, and WAF
- Analyze technology stack indicators in the HTML source

### 2. Active Reconnaissance
- Probe common web ports
- Enumerate directories and paths
- Check sensitive files (robots.txt, .env, .git)
- Discover API endpoints

### 3. Technology Stack Identification
- Front-end frameworks (React/Vue/Angular/jQuery)
- Back-end frameworks (Express/Django/Flask/Spring)
- CMS systems (WordPress/Joomla/custom)
- Database type

### 4. Output
- Target profile (IP/domain/ports/services/technology stack)
- Attack surface map (accessible paths, APIs, admin entry points)
