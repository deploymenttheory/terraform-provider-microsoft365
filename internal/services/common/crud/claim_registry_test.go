package crud

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestClaimKey_FirstAndSecondClaim(t *testing.T) {
	ResetClaims("test_claim_first")
	if ClaimKey("test_claim_first", "key") {
		t.Fatal("first claim reported the key as already claimed")
	}
	if !ClaimKey("test_claim_first", "key") {
		t.Fatal("second claim did not report the key as already claimed")
	}
}

func TestClaimKey_FoldsCase(t *testing.T) {
	ResetClaims("test_claim_case")
	if ClaimKey("test_claim_case", "AABBCCDD-1111") {
		t.Fatal("first claim reported the key as already claimed")
	}
	if !ClaimKey("test_claim_case", "aabbccdd-1111") {
		t.Fatal("a claim that differs only by case was not seen as already claimed")
	}
}

func TestClaimKey_ScopesAreIndependent(t *testing.T) {
	ResetClaims("test_claim_scope_a")
	ResetClaims("test_claim_scope_b")
	RecordClaim("test_claim_scope_a", "key")
	if ClaimKey("test_claim_scope_b", "key") {
		t.Fatal("a claim in one scope was seen by another scope")
	}
	if !ClaimKey("test_claim_scope_a", "KEY") {
		t.Fatal("a recorded claim was not seen in its own scope")
	}
}

func TestRecordClaim_IsSeenByClaimKey(t *testing.T) {
	ResetClaims("test_claim_record")
	RecordClaim("test_claim_record", "Key")
	RecordClaim("test_claim_record", "key")
	if !ClaimKey("test_claim_record", "KEY") {
		t.Fatal("a recorded claim was not seen by ClaimKey")
	}
}

func TestClaimKey_ConcurrentClaimsHaveOneWinner(t *testing.T) {
	const claimers = 100
	ResetClaims("test_claim_concurrent")
	var winners int32
	var start sync.WaitGroup
	var done sync.WaitGroup
	start.Add(1)
	for index := range claimers {
		done.Add(1)
		go func() {
			defer done.Done()
			start.Wait()
			key := "aabbccdd-2222"
			if index%2 == 0 {
				key = "AABBCCDD-2222"
			}
			if !ClaimKey("test_claim_concurrent", key) {
				atomic.AddInt32(&winners, 1)
			}
		}()
	}
	start.Done()
	done.Wait()
	if winners != 1 {
		t.Fatalf("winners = %d, want 1", winners)
	}
}

func TestResetClaims_ForgetsOnlyItsScope(t *testing.T) {
	RecordClaim("test_claim_reset", "key")
	RecordClaim("test_claim_reset_kept", "key")
	ResetClaims("test_claim_reset")
	if ClaimKey("test_claim_reset", "key") {
		t.Fatal("a reset scope still held its claim")
	}
	if !ClaimKey("test_claim_reset_kept", "key") {
		t.Fatal("resetting one scope forgot another scope's claim")
	}
}
