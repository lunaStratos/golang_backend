
/**
 * 게시판 페이징
 *
 * */
func BoardPaging(total int32, nowPage int32, cntPage int32) Page {
	var obj = Page{0, 0, 0, 0, 0, 0, 0, 0}

	obj.NowPage = nowPage
	obj.CntPerPage = cntPage
	obj.Total = total

	//제일 마지막 페이지 계산
	obj.LastPage = int32(math.Ceil(float64(float64(total) / float64(cntPage))))

	//시작, 끝 페이지 계산
	obj.EndPage = int32(math.Ceil(float64(float64(nowPage)/float64(cntPage)))) * cntPage

	if obj.LastPage < obj.EndPage {
		obj.EndPage = obj.LastPage
	}

	obj.StartPage = obj.NowPage - cntPage + 1

	if obj.StartPage < 1 {
		obj.StartPage = 1
	}

	// start === offset 값
	obj.End = nowPage*cntPage - 1
	obj.Start = obj.End - cntPage + 1
	return obj
}

type Page struct {
	NowPage    int32 // 현재 페이지
	StartPage  int32 // 시작 페이지
	EndPage    int32 // 종료 페이지
	Total      int32 // 총 글 갯수(all count)
	CntPerPage int32 // 1페이지 당 게시글 수 (설정값)
	LastPage   int32 // 마지막 페이지
	Start      int32 // 시작 위치 (offset값)
	End        int32 // 종료위치 (사용안함)
}
