set -eu

function is_prime()
{
    local n=$1
    if [ $n -lt 2 ]; then
        return 1
    fi
    for ((i = 2; i * i <= n; i++)); do
        if [ $((n % i)) -eq 0 ]; then
            return 1
        fi
    done
    return 0
}

for ((j = 2; j < 300; j++)); do
    if is_prime $j; then
        go run ./cmd/prove ${j} >small/${j}.json
        echo $j
    fi

done
