#!/bin/bash

cur=`dirname $(readlink -f $0)`
cwd=`pwd`

cd ${cur}/..

fmt="json"
cfg="configs/example.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} > ${cfg}

fmt="yaml"
cfg="configs/example.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} > ${cfg}

fmt="json"
cfg="configs/localhost.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} --local > ${cfg}

fmt="yaml"
cfg="configs/localhost.${fmt}"
echo "generate \"${cfg}\" in ${fmt}"
bin/checkah example --format=${fmt} --local > ${cfg}
