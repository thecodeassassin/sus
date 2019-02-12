#!/bin/bash

export GIT_HOOKS_PATH=./.git/hooks

echo "#!/bin/sh
make test" > ${GIT_HOOKS_PATH}/pre-commit

chmod +x ${GIT_HOOKS_PATH}/pre-commit
