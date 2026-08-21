# RedTeam Framework

A complete, open‑source C2 platform for educational and authorized red team engagements. Runs entirely without root.

## ⚠️ Legal Warning
**Only use on systems you own or have explicit written permission to test. Unauthorized access is illegal.**

## Features
- Web dashboard with real‑time client management
- Modular client: Recon, Persistence, Lateral Movement, Exfil, Shell
- No root required
- Cross‑platform (Windows, Linux, macOS)
- Docker support

## Quick Start
1. **Server**: `cd server && go run main.go` (or `docker-compose up`)
2. **Client**: Edit `client/main.go` – set `serverURL` to your server IP.
3. **Build client**: `cd client && ./build.sh`
4. **Deploy** client binary on target and run it.
5. Open browser to `http://server-ip:8080`

## Modules
- **Recon**: Host info, IPs, processes, container detection
- **Persistence**: User‑level startup (cron, Windows Startup)
- **Lateral**: SSH key scanner (bruteforce stub)
- **Exfil**: File read + base64 encoding
- **Shell**: Raw command execution

## Customization
Add your own modules in `client/modules/` and register them in `main.go` switch.
