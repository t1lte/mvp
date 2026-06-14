#!/usr/bin/env bash
set -e

export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export FABRIC_LOGGING_SPEC=ERROR
export CORE_LOGGING_LEVEL=ERROR
export ORDERER_GENERAL_LOGLEVEL=ERROR

CHANNEL="mychannel"
CC_LANG="go"
CC_VERSION="1"
CC_SEQUENCE="1"
CONTRACTS=("model1" "model2" "model3" "model4")
LOG_FILE="deploy_$(date +%Y%m%d_%H%M%S).log"

timestamp() {
    date "+%Y-%m-%d %H:%M:%S"
}

run_quiet() {
    local label="$1"
    shift
    local start_time=$(date +%s)
    
    printf "[%s] %-30s [PENDING]\n" "$(timestamp)" "$label"
    
    if "$@" >/dev/null 2>&1; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        printf "[%s] %-30s [OK] (%ds)\n" "$(timestamp)" "$label" "$duration"
        return 0
    else
        printf "[%s] %-30s [FAILED]\n" "$(timestamp)" "$label"
        echo "ERROR: Command failed. See $LOG_FILE for details." >&2
        "$@" >> "$LOG_FILE" 2>&1 || true
        return 1
    fi
}

echo "=============================================="
echo "Hyperledger Fabric Test Network Deployment"
echo "Started: $(timestamp)"
echo "Channel: $CHANNEL"
echo "Contracts: ${CONTRACTS[*]}"
echo "=============================================="
echo ""

if [ ! -f "./network.sh" ]; then
    echo "[$(timestamp)] ERROR: ./network.sh not found. Run from test-network directory." >&2
    exit 1
fi

run_quiet "Network up + create channel" \
    ./network.sh up createChannel -c "$CHANNEL" -ca

echo ""
echo "--- Chaincode Deployment ---"

for cc_name in "${CONTRACTS[@]}"; do
    run_quiet "Deploy $cc_name" \
        ./network.sh deployCC \
            -ccn "$cc_name" \
            -ccp "./${cc_name}" \
            -ccl "$CC_LANG" \
            -c "$CHANNEL" \
            -ccv "$CC_VERSION" \
            -ccs "$CC_SEQUENCE"
done

echo ""
echo "=============================================="
echo "Deployment Summary"
echo "Completed: $(timestamp)"
echo "=============================================="
echo ""
echo "Deployed chaincodes on channel '$CHANNEL':"
echo "---"

peer lifecycle chaincode querycommitted \
    --channelID "$CHANNEL" \
    --output json 2>/dev/null | \
    python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    for cc in data.get('chaincode_definitions', []):
        print(f\"  {cc['name']:20s} v{cc['version']} (seq: {cc['sequence']})\")
except:
    print('  (unable to parse chaincode list)')
" 2>/dev/null || echo "  (run 'peer lifecycle chaincode querycommitted --channelID $CHANNEL' for details)"

echo ""
echo "Log file: $LOG_FILE"
echo "To attach CLI: export \$(./network.sh env)"
