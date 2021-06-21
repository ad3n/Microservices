#!/usr/bin/env bash
set -e

# shellcheck disable=SC2043
for name in MAIN_AUTH_SERVICE
do
    eval value=\$$name
    # shellcheck disable=SC2154
    sed -i "s|\${${name}}|${value}|g" krakend.json
done

# shellcheck disable=SC2043
for name in SECURITY_SERVICE
do
    eval value=\$$name
    # shellcheck disable=SC2154
    sed -i "s|\${${name}}|${value}|g" krakend.json
done

# shellcheck disable=SC2043
for name in SERVICE1_SERVICE
do
    eval value=\$$name
    # shellcheck disable=SC2154
    sed -i "s|\${${name}}|${value}|g" krakend.json
done

# shellcheck disable=SC2043
for name in SERVICE2_SERVICE
do
    eval value=\$$name
    # shellcheck disable=SC2154
    sed -i "s|\${${name}}|${value}|g" krakend.json
done

# shellcheck disable=SC2043
for name in SERVICE3_SERVICE
do
    eval value=\$$name
    # shellcheck disable=SC2154
    sed -i "s|\${${name}}|${value}|g" krakend.json
done

/usr/bin/api-gateway run -c krakend.json
