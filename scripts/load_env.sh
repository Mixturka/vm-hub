#!/bin/sh

if [ -f ".env.test" ]; then
  if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    echo "Windows detected: Using PowerShell"
    while IFS= read -r line; do
      if [[ "$line" =~ ^[^#]*=.* ]]; then
        $key=$(echo "$line" | cut -d '=' -f 1)
        $value=$(echo "$line" | cut -d '=' -f 2-)
        [ -n "$key" ] && [ -n "$value" ] && $env:$key=$value
      fi
    done < ".env.test"

    echo ".env loaded with resolved variables"
  else
    set -a
    source ./.env.test
    set +a

    echo ".env loaded with resolved variables"
    echo $POSTGRES_MIGRATIONS_PATH
  fi
else
  echo ".env not found"
fi
