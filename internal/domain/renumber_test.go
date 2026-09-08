package domain

import "testing"

// Tái hiện đúng outline mà model đã lưu trong lần chạy hỏng: cả bốn cung đều index 0.
func TestRenumberVolumesFixesDuplicateZeroIndices(t *testing.T) {
	volumes := []VolumeOutline{
		{Index: 0, Arcs: []ArcOutline{{Index: 0}, {Index: 0}, {Index: 0}, {Index: 0}}},
		{Index: 1, Arcs: []ArcOutline{{Index: 0}, {Index: 0}}},
	}
	RenumberVolumes(volumes)

	for vi, v := range volumes {
		if v.Index != vi+1 {
			t.Fatalf("tập %d: index=%d, mong đợi %d", vi, v.Index, vi+1)
		}
		for ai, a := range v.Arcs {
			if a.Index != ai+1 {
				t.Fatalf("tập %d cung %d: index=%d, mong đợi %d", vi, ai, a.Index, ai+1)
			}
		}
	}
}
