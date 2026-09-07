package domain

import "strings"
import "testing"

// Tái hiện đại cương thật đã sinh ra truyện "viết lại rồi thêm một tí":
// 11/12 chương dùng chung một móc câu.
func TestStalledOutlineCatchesRepeatedHook(t *testing.T) {
	same := "Một tu sĩ trẻ xuất hiện với lời mời gia nhập tông phái"
	var entries []OutlineEntry
	for i := 1; i <= 12; i++ {
		e := OutlineEntry{Chapter: i, Title: "Chương", CoreEvent: "Giả Mộc học được điều thứ " + string(rune('a'+i))}
		if i > 1 {
			e.Hook = same
		}
		entries = append(entries, e)
	}
	got := StalledOutline(entries)
	if !strings.Contains(got, "hook") {
		t.Fatalf("phải bắt được hook lặp, nhận %q", got)
	}
}

func TestStalledOutlineAcceptsHealthyOutline(t *testing.T) {
	entries := []OutlineEntry{
		{Chapter: 1, CoreEvent: "Mất mùa, phải nộp giống", Hook: "Cây mai rung lên trong đêm"},
		{Chapter: 2, CoreEvent: "Gặp người lạ ngoài chợ", Hook: "Người lạ biết tên mẹ hắn"},
		{Chapter: 3, CoreEvent: "Bị đuổi khỏi thôn", Hook: "Hắn mang theo một hạt giống không rõ nguồn"},
	}
	if got := StalledOutline(entries); got != "" {
		t.Fatalf("đại cương lành mạnh không được báo lỗi, nhận %q", got)
	}
}

// Tái hiện sự cố thật: Tập 1 còn hai cung khung xương (38 chương) mà sách
// vẫn được tuyên bố hoàn thành sau 15 chương của Tập 2.
func TestSkeletonArcsSeesUnexpandedVolume(t *testing.T) {
	volumes := []VolumeOutline{
		{Index: 1, Arcs: []ArcOutline{
			{Index: 1, Title: "Nguyên Bản Dược Lâm", EstimatedChapters: 25},
			{Index: 2, Title: "Thảm Họa Của Linh Căn", EstimatedChapters: 13},
		}},
		{Index: 2, Arcs: []ArcOutline{
			{Index: 1, Title: "Linh Căn Giới Hạn", Chapters: []OutlineEntry{{Chapter: 1}}},
		}},
	}
	got := SkeletonArcs(volumes)
	if len(got) != 2 {
		t.Fatalf("phải thấy 2 cung khung xương, nhận %d: %v", len(got), got)
	}
	if !strings.Contains(got[0], "Nguyên Bản Dược Lâm") {
		t.Fatalf("phải nêu tên cung, nhận %q", got[0])
	}
}

func TestSkeletonArcsSilentWhenAllExpanded(t *testing.T) {
	volumes := []VolumeOutline{{Index: 1, Arcs: []ArcOutline{
		{Index: 1, Chapters: []OutlineEntry{{Chapter: 1}}},
	}}}
	if got := SkeletonArcs(volumes); len(got) != 0 {
		t.Fatalf("không được báo khi mọi cung đã khai triển, nhận %v", got)
	}
}
