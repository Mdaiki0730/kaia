#!/bin/bash

# テストの確率を測定するスクリプト
# go test ./blockchain を1000000回実行して失敗率を計算

TOTAL_RUNS=100
SUCCESS_COUNT=0
FAILURE_COUNT=0

echo "Starting $TOTAL_RUNS test runs..."
echo "This may take a while..."
echo ""

# 進捗表示のための変数
PROGRESS_INTERVAL=5
START_TIME=$(date +%s)

for i in $(seq 1 $TOTAL_RUNS); do
    if make test-others; then
        ((SUCCESS_COUNT++))
    else
        ((FAILURE_COUNT++))
    fi
    go clean -cache
    
    # 進捗表示
    if [ $((i % PROGRESS_INTERVAL)) -eq 0 ]; then
        CURRENT_TIME=$(date +%s)
        ELAPSED=$((CURRENT_TIME - START_TIME))
        
        if [ $ELAPSED -gt 0 ]; then
            RATE=$((i / ELAPSED))
            if [ $RATE -gt 0 ]; then
                REMAINING=$(( (TOTAL_RUNS - i) / RATE ))
                echo "[$i/$TOTAL_RUNS] Success: $SUCCESS_COUNT, Failure: $FAILURE_COUNT, Elapsed: ${ELAPSED}s, ETA: ${REMAINING}s"
            else
                echo "[$i/$TOTAL_RUNS] Success: $SUCCESS_COUNT, Failure: $FAILURE_COUNT, Elapsed: ${ELAPSED}s, ETA: calculating..."
            fi
        else
            echo "[$i/$TOTAL_RUNS] Success: $SUCCESS_COUNT, Failure: $FAILURE_COUNT, Elapsed: ${ELAPSED}s"
        fi
    fi
done

# 結果を計算
TOTAL=$((SUCCESS_COUNT + FAILURE_COUNT))
SUCCESS_RATE=$(echo "scale=6; $SUCCESS_COUNT * 100 / $TOTAL" | bc)
FAILURE_RATE=$(echo "scale=6; $FAILURE_COUNT * 100 / $TOTAL" | bc)

END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

echo ""
echo "=========================================="
echo "Test Results:"
echo "=========================================="
echo "Total runs:     $TOTAL_RUNS"
echo "Success:        $SUCCESS_COUNT ($SUCCESS_RATE%)"
echo "Failure:        $FAILURE_COUNT ($FAILURE_RATE%)"
echo "Total time:     ${TOTAL_TIME}s"
echo "=========================================="

