#!/usr/bin/env bash
set -e

CHANNEL="mychannel"
ORDERER="localhost:7050"
ORDERER_TLS_CA="${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
QUERY_RETRIES=3
QUERY_DELAY=2

declare -A ORG_CONFIG
ORG_CONFIG[Org1MSP]="localhost:7051|${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt|${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
ORG_CONFIG[Org2MSP]="localhost:9051|${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt|${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp"

set_org_env() {
    local org="$1"
    local cfg="${ORG_CONFIG[$org]}"
    if [[ -z "$cfg" ]]; then
        echo "ERROR: Unknown organization '$org'" >&2
        return 1
    fi
    IFS='|' read -r addr tls msp <<< "$cfg"
    if [[ ! -e "$tls" ]]; then echo "ERROR: TLS cert not found: $tls" >&2; return 1; fi
    if [[ ! -e "$msp" ]]; then echo "ERROR: MSP directory not found: $msp" >&2; return 1; fi
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="$org"
    export CORE_PEER_TLS_ROOTCERT_FILE="$tls"
    export CORE_PEER_MSPCONFIGPATH="$msp"
    export CORE_PEER_ADDRESS="$addr"
}

invoke_chaincode() {
    local cc="$1" func="$2" args="$3"
    local peer_args=""
    for o in Org1MSP Org2MSP; do
        IFS='|' read -r a t _ <<< "${ORG_CONFIG[$o]}"
        peer_args+=" --peerAddresses $a --tlsRootCertFiles $t"
    done
    local cmd="peer chaincode invoke -o $ORDERER --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_TLS_CA -C $CHANNEL -n $cc $peer_args -c '{\"function\":\"$func\",\"Args\":$args}'"
    echo "INVOKE: $cc.$func($args)"
    local output
    output=$(eval "$cmd" 2>&1) || true
    if echo "$output" | grep -qE "(status:200|Response.*200|committed with status VALID)"; then
        echo "STATUS: SUCCESS"
        return 0
    else
        echo "STATUS: FAILED"
        echo "$output" | tail -n 5 | sed 's/^/  /' >&2
        return 1
    fi
}

query_with_retry() {
    sleep 5
    local cc="$1" func="$2" args="$3"
    echo ""
    echo "QUERY: $cc.$func($args)"
    for attempt in $(seq 1 $QUERY_RETRIES); do
        local result
        result=$(peer chaincode query -C "$CHANNEL" -n "$cc" -c "{\"function\":\"$func\",\"Args\":$args}" 2>&1) || true
        if echo "$result" | grep -qiE "(error|failed|not found)"; then
            if (( attempt < QUERY_RETRIES )); then
                echo "  [RETRY] $attempt/$QUERY_RETRIES - Ledger sync pending, waiting ${QUERY_DELAY}s..."
                sleep "$QUERY_DELAY"
                continue
            else
                echo "  [FAIL] Query failed after $QUERY_RETRIES attempts"
                echo "  $result" | tail -n 2 | sed 's/^/  /' >&2
                return 1
            fi
        else
            echo "  [OK] Query succeeded"
            if command -v jq &>/dev/null; then echo "$result" | jq . 2>/dev/null || echo "$result"
            else echo "$result"; fi
            return 0
        fi
    done
}

print_usage() {
    cat <<'EOF'
Usage: invoke_cc.sh [OPTIONS]
Required:
  -c CONTRACT   Chaincode name
  -o ORG        Organization (Org1MSP or Org2MSP)
  -f FUNCTION   Function to invoke
  -a ARGS       Arguments as JSON array
Optional:
  -q FUNC ARGS  Override query function and arguments
  -n            Skip automatic query after invoke
  -h            Show help
EOF
}

CC=""; ORG=""; FUNC=""; ARGS=""
QFUNC="GetInstanceState"; QARGS=""; NO_QUERY=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        -c) CC="$2"; shift 2 ;;
        -o) ORG="$2"; shift 2 ;;
        -f) FUNC="$2"; shift 2 ;;
        -a) ARGS="$2"; shift 2 ;;
        -q) QFUNC="$2"; QARGS="$3"; shift 3 ;;
        -n) NO_QUERY=true; shift ;;
        -h) print_usage; exit 0 ;;
        *) echo "ERROR: Unknown option '$1'"; exit 1 ;;
    esac
done

if [[ -z "$CC" || -z "$ORG" || -z "$FUNC" || -z "$ARGS" ]]; then
    echo "ERROR: Missing required parameters" >&2
    print_usage >&2
    exit 1
fi

if ! echo "$ARGS" | python3 -c "import sys,json; json.loads(sys.stdin.read())" 2>/dev/null; then
    echo "ERROR: Arguments must be a valid JSON array" >&2
    exit 1
fi

[[ -z "$QARGS" ]] && QARGS="$ARGS"

echo "========================================"
echo "Contract: $CC | Org: $ORG | Channel: $CHANNEL"
echo "========================================"

if [[ ! -e "$ORDERER_TLS_CA" ]]; then echo "ERROR: Orderer TLS CA not found" >&2; exit 1; fi
set_org_env "$ORG" || exit 1

invoke_chaincode "$CC" "$FUNC" "$ARGS"
INVOKE_STATUS=$?

if (( INVOKE_STATUS == 0 )) && [[ "$NO_QUERY" == false ]]; then
    query_with_retry "$CC" "$QFUNC" "$QARGS" || true
fi

echo ""
echo "========================================"
echo "Overall: $([ $INVOKE_STATUS -eq 0 ] && echo 'SUCCESS' || echo 'FAILED')"
echo "========================================"
exit $INVOKE_STATUS
