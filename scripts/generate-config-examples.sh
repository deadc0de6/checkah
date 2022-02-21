#!/bin/bash

rl="readlink -f"
if ! ${rl} "${0}" >/dev/null 2>&1; then
  rl="realpath"

  if ! hash ${rl}; then
    echo "\"${rl}\" not found!" && exit 1
  fi
fi

cur=`dirname $(${rl} $0)`
cwd=`pwd`

cd ${cur}/..

fmt="json"
cfg="./configs/example.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} > ${cfg}

fmt="yaml"
cfg="./configs/example.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} > ${cfg}

fmt="json"
cfg="./configs/localhost.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} --local > ${cfg}

fmt="yaml"
cfg="./configs/localhost.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} --local > ${cfg}
