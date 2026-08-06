#!/usr/bin/env python3
"""Import VulnClaw skills + warstories into Talon's skills/ tree.

Converts VulnClaw YAML-frontmatter SKILL.md files into Talon's flat markdown
format (`# stage:` / `# category:` comment headers parsed by
internal/core/skills.go parseSkillFile).

Usage:
    python3 scripts/import_vulnclaw_skills.py [--src /path/to/VulnClaw] [--dst skills]

Translation of Chinese bodies is NOT done here; the script prints a list of
files containing CJK text at the end so they can be translated in place.
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

CJK_RE = re.compile(r"[一-鿿぀-ヿ가-힯]")

# routing.phases / name heuristics -> Talon stage
PHASE_MAP = {
    "recon": "recon",
    "discovery": "recon",
    "vuln_discovery": "recon",
    "exploit": "exploit",
    "exploitation": "exploit",
    "attack": "exploit",
    "post": "post_exploit",
    "post_exploit": "post_exploit",
    "post-exploitation": "post_exploit",
    "report": "report",
    "reporting": "report",
}

NAME_STAGE_HINTS = [
    ("recon", "recon"),
    ("osint", "recon"),
    ("intake", "recon"),
    ("cve-lookup", "recon"),
    ("cve-triage", "recon"),
    ("postex", "post_exploit"),
    ("post-exploitation", "post_exploit"),
    ("intranet", "post_exploit"),
    ("report", "report"),
]


def parse_frontmatter(text: str) -> tuple[dict, str]:
    """Minimal YAML-frontmatter parser for VulnClaw's simple schema.

    Returns (meta, body). meta has 'name', 'description', 'routing' (dict of
    key -> list[str]). Unknown/complex lines are ignored.
    """
    if not text.startswith("---"):
        return {}, text
    end = text.find("\n---", 3)
    if end == -1:
        return {}, text
    fm = text[3:end].strip("\n")
    body = text[end + 4 :].lstrip("\n")
    meta: dict = {"routing": {}}
    current_section = None
    for raw in fm.split("\n"):
        line = raw.rstrip()
        if not line.strip() or line.strip().startswith("#"):
            continue
        if not line.startswith((" ", "\t")) and ":" in line:
            key, _, val = line.partition(":")
            key = key.strip()
            val = val.strip()
            if val == "":
                current_section = key
                if key != "routing":
                    meta.setdefault(key, {})
                continue
            current_section = None
            meta[key] = val.strip('"').strip("'")
        elif current_section == "routing" and ":" in line:
            key, _, val = line.partition(":")
            key = key.strip()
            val = val.strip()
            items = [v.strip().strip('"').strip("'") for v in val.strip("[]").split(",") if v.strip()]
            meta["routing"][key] = items
    return meta, body


def map_stage(name: str, routing: dict) -> str:
    phases = routing.get("phases") or routing.get("phase") or []
    for ph in phases:
        ph = ph.lower().strip()
        if ph in PHASE_MAP:
            return PHASE_MAP[ph]
    low = name.lower()
    for hint, stage in NAME_STAGE_HINTS:
        if hint in low:
            return stage
    return "exploit"


def category_for(name: str, source: str) -> str:
    if source == "warstories":
        return "warstories"
    if source == "core":
        return "core"
    if name.startswith("redteam-"):
        return "redteam"
    if name.startswith("ctf-"):
        return "ctf"
    return "specialized"


def emit_skill(name: str, description: str, stage: str, category: str, body: str,
               references: list[tuple[str, str]]) -> str:
    out = [f"# stage: {stage}", f"# category: {category}", ""]
    if description:
        out.append(f"> {description}")
        out.append("")
    out.append(body.strip())
    for ref_name, ref_body in references:
        out.append("")
        out.append(f"## References — {ref_name}")
        out.append("")
        out.append(ref_body.strip())
    return "\n".join(out).rstrip() + "\n"


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--src", default="/home/anubhavg/Desktop/VulnClaw")
    ap.add_argument("--dst", default="skills")
    args = ap.parse_args()

    src = Path(args.src)
    dst = Path(args.dst) / "vulnclaw"
    written: list[Path] = []
    cjk_files: list[Path] = []

    def write_skill(name: str, source: str, text: str, references: list[tuple[str, str]]):
        meta, body = parse_frontmatter(text)
        description = meta.get("description", "")
        stage = map_stage(name, meta.get("routing", {}))
        category = category_for(name, source)
        content = emit_skill(name, description, stage, category, body, references)
        target = dst / category / f"{name}.md"
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")
        written.append(target)
        if CJK_RE.search(content):
            cjk_files.append(target)

    # core/*.md
    for f in sorted((src / "vulnclaw/skills/core").glob("*.md")):
        write_skill(f.stem, "core", f.read_text(encoding="utf-8"), [])

    # specialized/*/SKILL.md (+ references/)
    for d in sorted((src / "vulnclaw/skills/specialized").iterdir()):
        skill_md = d / "SKILL.md"
        if not d.is_dir() or not skill_md.exists():
            continue
        refs = []
        ref_dir = d / "references"
        if ref_dir.is_dir():
            for rf in sorted(ref_dir.glob("*.md")):
                refs.append((rf.stem, rf.read_text(encoding="utf-8")))
        write_skill(d.name, "specialized", skill_md.read_text(encoding="utf-8"), refs)

    # warstories (skip README template index)
    ws = src / "vulnclaw/warstories"
    for f in sorted(ws.glob("*.md")):
        if f.stem == "README":
            continue
        write_skill(f.stem, "warstories", f.read_text(encoding="utf-8"), [])

    print(f"Wrote {len(written)} skill files under {dst}/")
    for p in written:
        print(f"  {p}")
    if cjk_files:
        print(f"\n{len(cjk_files)} files contain CJK text (need translation):")
        for p in cjk_files:
            print(f"  {p}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
