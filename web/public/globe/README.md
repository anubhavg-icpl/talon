# Globe assets

Source video: `~/Videos/dark_planet.mp4` (same bytes as `dark-planet.mp4` here).

| File | Use |
|------|-----|
| `dark-planet.mp4` | Primary Target HUD / engagement wallpaper (looped video) |
| `operator-globe-hud.webp` | Poster still + fallback if video fails |
| `earth-dark.jpg` | three-globe earth texture (Live 3D) |
| `earth-topology.png` | three-globe bump map (Live 3D) |
| `night-sky.png` | Background-variant sky (Live 3D) |

Mirror path (same media): `/wallpaper/dark-planet.mp4`, `/wallpaper/operator-globe-hud.webp`.

UI default is **video wallpaper** via `GlobeWallpaper` (`VIDEO_SRC=/globe/dark-planet.mp4`).  
Operators can switch to **Live 3D** (WebGL `TalonGlobe`) when needed.

Auth proxy must allow `/globe/*` and `*.mp4` (see `web/src/proxy.ts`) so the video loads before login and after.
