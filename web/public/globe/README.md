# Globe assets

Source video: `~/Videos/dark_planet.mp4` (same bytes as `dark-planet.mp4` here).

| File | Use |
|------|-----|
| `dark-planet.webm` | Primary Target HUD wallpaper (VP9 WebM, preferred) |
| `dark-planet.mp4` | Fallback loop for engines without WebM |
| `earth-dark.jpg` | three-globe earth texture (Live 3D) |

Mirror path: `/wallpaper/dark-planet.webm`, `.mp4`.

UI default is **video wallpaper** via `GlobeWallpaper` (`dark-planet.webm` → `.mp4` fallback).  
Operators can switch to **Live 3D** (WebGL `TalonGlobe`) when needed.

Auth proxy allows `/globe/*`, `/wallpaper/*`, and `*.webm` / `*.mp4` (see `web/src/proxy.ts`).
