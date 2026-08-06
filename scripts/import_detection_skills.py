#!/usr/bin/env python3
"""Import Detection-Skills into Talon's skills/ tree.

Converts Detection-Skills YAML-frontmatter SKILL.md files into Talon's
flat markdown format (# stage: / # category: comment headers parsed by
internal/core/skills.go parseSkillFile).

Type mapping:
  triage        → stage: recon, category: triage
  investigation → stage: exploit, category: investigation
  tuning        → stage: report, category: tuning
"""
from __future__ import annotations
import argparse, re, sys
from pathlib import Path

def parse_frontmatter(text: str) -> tuple[dict, str]:
    if not text.startswith("---"):
        return {}, text
    end = text.find("\n---", 3)
    if end == -1:
        return {}, text
    fm = text[3:end].strip("\n")
    body = text[end + 4:].lstrip("\n")
    meta: dict = {"labels": []}
    current = None
    for raw in fm.split("\n"):
        line = raw.rstrip()
        if not line.strip() or line.strip().startswith("#"):
            continue
        if not line.startswith((" ", "\t")) and ":" in line:
            key, _, val = line.partition(":")
            key, val = key.strip(), val.strip()
            if val == "" and key in ("metadata",):
                current = "metadata"
            elif val == "" and key in ("labels",):
                current = "labels"
            elif val:
                # Clean up YAML multi-line indicator
                if val.startswith("|-"):
                    val = val[2:].strip()
                meta[key] = val.strip('"').strip("'")
                current = None
        elif current == "labels" and line.strip().startswith("-"):
            meta["labels"].append(line.strip()[1:].strip())
        elif current == "metadata" and ":" in line:
            key, _, val = line.partition(":")
            meta[key.strip()] = val.strip().strip('"').strip("'")
    return meta, body

TYPE_STAGE = {
    "triage": "recon",
    "investigation": "exploit",
    "tuning": "report",
}

def slug(name: str) -> str:
    s = re.sub(r'[^a-z0-9]+', '-', name.lower()).strip('-')
    return s

def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--src", default="/home/anubhavg/Desktop/Detection-Skills")
    ap.add_argument("--dst", default="skills")
    args = ap.parse_args()
    src = Path(args.src)
    dst = Path(args.dst) / "detection"
    written: list[Path] = []
    for type_dir in sorted((src / "skills").iterdir()):
        if not type_dir.is_dir() or type_dir.name not in TYPE_STAGE:
            continue
        skill_type = type_dir.name
        stage = TYPE_STAGE[skill_type]
        for skill_dir in sorted(type_dir.iterdir()):
            skill_md = skill_dir / "SKILL.md"
            if not skill_dir.is_dir() or not skill_md.exists():
                continue
            meta, body = parse_frontmatter(skill_md.read_text("utf-8"))
            name = meta.get("name", skill_dir.name)
            description = meta.get("description", "")
            version = meta.get("version", "1.0.0")
            author = meta.get("author", "")
            labels = meta.get("labels", [])
            s = slug(name)
            content = f"# stage: {stage}\n# category: {skill_type}\n\n"
            content += f"> {description}\n\n"
            if labels:
                content += f"**Labels:** {', '.join(labels)}\n\n"
            if author:
                content += f"**Author:** {author} · **Version:** {version}\n\n"
            content += "---\n\n"
            content += body.strip() + "\n"
            target = dst / skill_type / f"{s}.md"
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_text(content, "utf-8")
            written.append(target)
    print(f"Wrote {len(written)} detection skill files under {dst}/")
    for p in written:
        print(f"  {p}")
    return 0

if __name__ == "__main__":
    sys.exit(main())
