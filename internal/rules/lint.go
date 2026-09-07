package rules

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Lint 内置产品底线检查：扫描正文中的机制残留，与用户规则无关，commit 时始终执行。
// 与 Check 同契约——仅返事实（铁律一），不阻断流程，由评审/用户裁定。
//
// 当前三类（全部来自真实长跑产物的实证缺陷）：
//   - markdown_residue：正文残留 ** 加粗、首行之外的 # 标题行（导出 txt 会裸露符号）
//   - non_cjk_fragments / cjk_leak：文字混杂片段，方向按正文主文字自动判定
//     （中文正文报拉丁片段；越南语等拉丁文字正文报汉字片段）
func Lint(text string) []Violation {
	var vs []Violation
	vs = appendMarkdownResidue(vs, text)
	vs = appendScriptMixing(vs, text)
	vs = appendBrokenWords(vs, text)
	return vs
}

func appendMarkdownResidue(vs []Violation, text string) []Violation {
	if n := strings.Count(text, "**"); n > 0 {
		vs = append(vs, Violation{
			Rule:     "markdown_residue",
			Target:   "**",
			Actual:   n,
			Severity: SeverityWarning,
		})
	}
	headings := 0
	seenContent := false
	for line := range strings.SplitSeq(text, "\n") {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		// 第一个非空行的 # 标题是章文件的合法格式（不按行号写死，容忍前导空行）
		first := !seenContent
		seenContent = true
		if !first && strings.HasPrefix(t, "#") {
			headings++
		}
	}
	if headings > 0 {
		vs = append(vs, Violation{
			Rule:     "markdown_residue",
			Target:   "#",
			Actual:   headings,
			Severity: SeverityWarning,
		})
	}
	return vs
}

var (
	latinFragmentRe = regexp.MustCompile(`[A-Za-z]{2,}`)
	// 含 CJK 标点：漏进越南语正文的不只是汉字，还有「。，、？！」这类全角标点。
	// 实测第 4 章正文里一个孤立的「。」——纯汉字规则完全看不见。
	cjkFragmentRe   = regexp.MustCompile(`[\p{Han}\x{3000}-\x{303f}\x{ff01}-\x{ff5e}]+`)
)

// appendScriptMixing 报告正文里「另一种文字」的混入片段。
//
// 哪种文字算混入由正文自身决定，不读配置：中文正文里裸混 "pattern" 是缺陷，
// 而越南语正文每个词都是拉丁字母，按同一条规则会在每章误报上千次，把评审
// 上下文淹掉——实测一章命中 1021 次而全章并无一个汉字。反过来，越南语正文
// 里漏出「根系之力」才是真缺陷，旧规则完全看不见。
//
// 因此先按字符占比判定正文主文字，再只报少数派。两种题材各自成立，且配置
// 写错也不会失效。合法外来词（品牌名/缩写）仍会命中——warning 级事实，由评审裁定。
func appendScriptMixing(vs []Violation, text string) []Violation {
	latin := latinFragmentRe.FindAllString(text, -1)
	han := cjkFragmentRe.FindAllString(text, -1)

	// 比的是「汉字个数」与「拉丁词个数」，不是两边的字符数：一个汉字约等于一个词，
	// 而一个拉丁词有好几个字母。按字符数比，中文正文里混几个 "pattern"/"DNA"
	// 就会把拉丁字符数顶过汉字数，从而误判正文语种。
	rule, matches := "non_cjk_fragments", latin
	if runeCount(han) <= len(latin) {
		// 正文是拉丁文字（越南语等）：汉字才是混入。
		rule, matches = "cjk_leak", han
	}
	if len(matches) == 0 {
		return vs
	}

	seen := make(map[string]struct{})
	var examples []string
	for _, m := range matches {
		if _, ok := seen[m]; ok {
			continue
		}
		seen[m] = struct{}{}
		if len(examples) < 3 {
			examples = append(examples, m)
		}
	}
	return append(vs, Violation{
		Rule:     rule,
		Target:   strings.Join(examples, "、"),
		Actual:   len(matches),
		Severity: SeverityWarning,
	})
}

func runeCount(ss []string) int {
	n := 0
	for _, s := range ss {
		n += utf8.RuneCountInString(s)
	}
	return n
}

// brokenWordRe 匹配被段落分隔切断的单词：一行以字母结尾，跨过空行后又以小写字母开头。
// 越南语正文实测：「n Tông, Ng」+ 空行 +「ọc Lâm dừng bước」——人名 Ngọc 被劈成两半。
// 拉丁文字里一个词不会跨段落，因此这个形状没有正当写法，可直接判为缺陷。
var brokenWordRe = regexp.MustCompile(`(?m)[\p{L}]\n\s*\n[[:space:]]*[\p{Ll}]`)

func appendBrokenWords(vs []Violation, text string) []Violation {
	matches := brokenWordRe.FindAllString(text, -1)
	if len(matches) == 0 {
		return vs
	}
	return append(vs, Violation{
		Rule:     "broken_word",
		Target:   strings.Join(strings.Fields(matches[0]), "⏎"),
		Actual:   len(matches),
		Severity: SeverityWarning,
	})
}
