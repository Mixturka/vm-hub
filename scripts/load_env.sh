#!/bin/sh

#!/bin/sh

if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs -I {} echo {})
  for var in $(grep -o '^[A-Za-z_][A-Za-z0-9_]*' .env); do
    eval "export $var=$(eval echo \$$var)"
  done

  echo ".env loaded with resolved variables"
  echo $POSTGRES_MIGRATIONS_PATH
else
  echo ".env not found"
fi
