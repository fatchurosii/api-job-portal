#!/usr/bin/env bash
set -euo pipefail

# konfigurasi default
MIGRATE_CMD="${MIGRATE_CMD:-migrate}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-database/migrations}"
ROOT_DIR="${ROOT_DIR:-$(pwd)}"

# --- load .env (for DATABASE_URL) ------------------------------------------
if [[ -f "$ROOT_DIR/.env" ]]; then
  # simple loader: export vars from .env (lines like KEY=VALUE)
  # note: this is simple and assumes .env doesn't contain malicious shell code
  # if your .env has spaces/special chars in values, prefer to export DATABASE_URL directly.
  # shellcheck disable=SC2046
  export $(grep -E '^[A-Za-z_][A-Za-z0-9_]*=' "$ROOT_DIR/.env" | xargs) || true
fi

# build DATABASE_URL from DB_* env if DATABASE_URL not provided
: "${DATABASE_URL:=""}"

# set sensible defaults
: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_SSLMODE:=disable}"

# If DATABASE_URL still empty, build from DB_* (works for postgres)
if [[ -z "${DATABASE_URL}" ]]; then
  # require at least DB_USER and DB_NAME
  if [[ -z "${DB_USER:-}" || -z "${DB_NAME:-}" ]]; then
    echo "ERROR: DATABASE_URL not set and DB_USER/DB_NAME missing. Export DATABASE_URL or set DB_USER and DB_NAME."
    exit 2
  fi

  # NOTE: if DB_PASSWORD contains special chars, this simple concatenation may break.
  # Prefer to export DATABASE_URL directly in that case.
  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD:-}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
  export DATABASE_URL
fi

die() { echo "✖ $*" >&2; exit 1; }
need_db() { [[ -n "${DATABASE_URL:-}" ]] || die "DATABASE_URL not set (put it in .env or export it)"; }

run_migrate() {
  echo "+ ${MIGRATE_CMD} -path ${MIGRATIONS_DIR} -database ${DATABASE_URL} $*"
  "${MIGRATE_CMD}" -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" "$@"
}

usage() {
  cat <<EOF
Usage: $(basename "$0") <command> [args...]

Commands:
  create <name>    Create new migration files (up + down)
  up [n]           Apply migrations (or steps n)
  down [n]         Rollback n steps (default 1)
  rollback         Rollback 1 step
  goto <version>   Migrate to specific version
  status           Show current migration version
  force <version>  Force set migration version (clears dirty flag) -- use with caution
  drop             Drop everything (destructive)
  help             Show this help

Environment:
  You can either export DATABASE_URL directly, or provide DB_USER, DB_NAME (and optionally DB_PASSWORD, DB_HOST, DB_PORT, DB_SSLMODE) or a .env file in project root.

Notes:
  - 'force' only sets the recorded migration version and dirty flag in the database. It does NOT run SQL migrations.
  - Use 'force' when you have inspected/fixed partial migrations and want to clear the dirty state.
  - Always backup the database before running destructive commands like 'force' or 'drop'.

EOF
  exit 1
}

if [[ $# -lt 1 ]]; then
  usage
fi

cmd="$1"
shift || true

case "$cmd" in
  create)
    if [[ $# -lt 1 ]]; then
      echo "ERROR: create requires a name"
      usage
    fi
    name="$1"
    echo "Creating migration: $name"
    "${MIGRATE_CMD}" create -ext sql -dir "${MIGRATIONS_DIR}" -seq "${name}"
    ;;

  up)
    need_db
    if [[ $# -ge 1 ]]; then
      run_migrate steps "$1"
    else
      run_migrate up
    fi
    ;;

  down)
    need_db
    if [[ $# -ge 1 ]]; then
      run_migrate steps "-${1}"
    else
      run_migrate steps -1
    fi
    ;;

  rollback)
    need_db
    run_migrate steps -1
    ;;

  goto)
    need_db
    if [[ $# -lt 1 ]]; then
      echo "ERROR: goto requires a version number"
      usage
    fi
    run_migrate goto "$1"
    ;;

  status)
    need_db
    run_migrate version
    ;;

  force)
    need_db
    if [[ $# -lt 1 ]]; then
      echo "ERROR: force requires a version number"
      usage
    fi
    version="$1"
    echo "WARNING: you are about to force migration version => ${version}"
    echo "This will set the recorded migration version and clear dirty flag in the DB."
    echo "Make sure you've inspected/fixed the migration SQL and backed up the DB."
    read -r -p "Type 'FORCE' to proceed: " confirm
    if [[ "$confirm" != "FORCE" ]]; then
      echo "Aborted force."
      exit 0
    fi
    run_migrate force "${version}"
    ;;

  drop)
    need_db
    echo "WARNING: drop will delete all data. Type DROP to confirm:"
    read -r confirm
    if [[ "$confirm" == "DROP" ]]; then
      run_migrate drop
    else
      echo "Aborted."
      exit 0
    fi
    ;;

  help|-h|--help)
    usage
    ;;

  *)
    echo "Unknown command: $cmd"
    usage
    ;;
esac
