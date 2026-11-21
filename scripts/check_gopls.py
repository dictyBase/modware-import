"""Repository gopls gate. Invoke through bash scripts/check-gopls.sh.

Lex Go comments and string/rune literals to recognize complete generated line
comments anywhere in a file. Do not mistake markers in strings or block comments
for generated provenance. Generated client bodies are never rewritten.
"""
from collections import Counter
from pathlib import Path
import re
import subprocess
import sys

GO_TOKENS = re.compile(
    r'//[^\n\r]*|/\*.*?\*/|"(?:\\.|[^"\\])*"|`[^`]*`|\'(?:\\.|[^\'\\])*\'',
    re.DOTALL,
)
GENERATED_MARKER = re.compile(r'// Code generated .* DO NOT EDIT\.')
DIAGNOSTIC = re.compile(r'^(.+\.go):\d+:\d+(?:-\d+(?::\d+)?)?: (.+)$')

# Exact file + symbol + message, bounded by audited baseline occurrences.
# The generated service interface requires the deprecated request types/fields;
# Attributes is an existing legacy input fallback. No other deprecations qualify.
LEGACY_FILE = 'internal/server/feature_annotation_deprecated.go'
DEPRECATION = (' is deprecated: Marked as deprecated in '
               'dictybase/feature_annotation/feature_annotation.proto.')
LEGACY_LIMITS = {
    ('internal/server/feature_annotation.go', 'req.Attributes' + DEPRECATION): 2,
    (LEGACY_FILE, 'feature.RemoveTagRequest' + DEPRECATION): 1,
    (LEGACY_FILE, 'feature.UpdateTagRequest' + DEPRECATION): 2,
    (LEGACY_FILE, 'req.Id' + DEPRECATION): 12,
    (LEGACY_FILE, 'req.Tag' + DEPRECATION): 17,
    (LEGACY_FILE, 'req.Tag.Tag' + DEPRECATION): 6,
    (LEGACY_FILE, 'req.Tag.UpdatedBy' + DEPRECATION): 3,
    (LEGACY_FILE, 'req.Tag.Value' + DEPRECATION): 2,
}


class DiagnosticsError(Exception):
    """Unexpected or over-budget gopls diagnostics."""


def is_generated(source):
    for token in GO_TOKENS.finditer(source):
        if GENERATED_MARKER.fullmatch(token[0]):
            prefix = source[source.rfind('\n', 0, token.start()) + 1:token.start()]
            if not prefix.strip():
                return True
    return False


def select_files(root, candidates):
    return [name for name in candidates if name and not name.startswith('vendor/')
            and not is_generated((root / name).read_text())]


def diagnostic_key(line, root):
    match = DIAGNOSTIC.fullmatch(line)
    if match is None:
        return None
    name = Path(match[1])
    if name.is_absolute():
        try:
            name = name.relative_to(root)
        except ValueError:
            return None
    return name.as_posix(), match[2]


def audit_diagnostics(output, root, report):
    counts = Counter()
    failures = []
    accepted = 0
    for line in output.splitlines():
        if not line.strip():
            continue
        key = diagnostic_key(line, root)
        counts[key] += 1
        if counts[key] <= LEGACY_LIMITS.get(key, 0):
            accepted += 1
        else:
            failures.append(line)
    print(f'audited compatibility diagnostics: {accepted}', file=report)
    for key in sorted(k for k in counts if k in LEGACY_LIMITS):
        print(f'  {counts[key]}/{LEGACY_LIMITS[key]} {key[0]}|{key[1]}', file=report)
    if failures:
        raise DiagnosticsError('unexpected gopls diagnostics:\n' + '\n'.join(failures))


def run(root):
    listed = subprocess.run(
        ['git', 'ls-files', '-z', '--cached', '--others', '--exclude-standard', '*.go'],
        cwd=root, check=True, capture_output=True,
    )
    selected = select_files(root, listed.stdout.decode().rstrip('\x00').split('\x00'))
    result = subprocess.run(['gopls', 'check', '-severity=hint', *selected],
                            cwd=root, capture_output=True, text=True)
    if result.returncode:
        raise DiagnosticsError(f'gopls failed ({result.returncode}):\n{result.stdout}{result.stderr}')
    print(f'gopls checked {len(selected)} nongenerated first-party files')
    audit_diagnostics(result.stdout + result.stderr, root, sys.stdout)


if __name__ == '__main__':
    try:
        run(Path(__file__).resolve().parents[1])
    except (DiagnosticsError, OSError, subprocess.CalledProcessError) as error:
        print(error, file=sys.stderr)
        sys.exit(1)
