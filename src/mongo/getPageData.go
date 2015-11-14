package mongo

import (
	l "logger"

	"errors"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
)

func GetIndexPageData() (*Index, error) {
	l.Output(nil, l.I_M_GetPage, l.Info)

	var teamArray []string
	for name, _ := range teams {
		teamArray = append(teamArray, name)
	}

	sort.Strings(teamArray)
	index := &Index{
		Team:   teamArray,
		Length: len(teams),
	}
	return index, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {
	l.Output(nil, l.I_M_GetPage, l.Info)

	var leadersBoard LeadersBoard
	var userScore UserScore

	for _, player := range players {
		var totalPar, totalScore, passedHoleCnt int
		for holeIndex, playerScore := range player.Score {
			holeNum := holeIndex + 1
			if playerScore["total"].(int) != 0 {
				totalPar += fields[holeNum].Par
				total, _ := playerScore["total"].(int)
				totalScore += total
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore - totalPar,
			Total: totalScore,
			Name:  users[player.UserId].Name,
			Hole:  passedHoleCnt,
		}
		leadersBoard.Ranking = append(leadersBoard.Ranking, userScore)
	}
	sort.Sort(sortByScore(leadersBoard.Ranking))

	return &leadersBoard, nil
}

func GetScoreEntrySheetPageData(teamName string, holeString string) (*ScoreEntrySheet, error) {
	l.Output(
		logrus.Fields{"Team Name": teamName, "Hole String": holeString},
		l.I_M_GetPage,
		l.Info,
	)

	if len(holeString) == 0 {
		l.PutErr(nil, l.Trace(), l.E_Nil, teamName)
		return nil, errors.New("hole string is nil")
	}
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	field := fields[holeNum]

	userIds := teams[teamName].UserIds
	member := make([]string, len(userIds))
	total := make([]int, len(userIds))
	putt := make([]int, len(userIds))
	for i, userId := range userIds {
		member[i] = users[userId].Name
		total[i] = players[userId].Score[holeIndex]["total"].(int)
		putt[i] = players[userId].Score[holeIndex]["putt"].(int)
	}

	scoreEntrySheet := ScoreEntrySheet{
		Team:    teamName,
		Hole:    holeNum,
		Member:  member,
		UserIds: userIds,
		Par:     field.Par,
		Yard:    field.Yard,
		Total:   total,
		Putt:    putt,
		Excnt:   excnt[teamName][holeNum],
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {
	l.Output(
		logrus.Fields{"Team Name": teamName},
		l.I_M_GetPage,
		l.Info,
	)

	userIds := teams[teamName].UserIds
	member := make([]string, len(userIds))
	apply := make([]int, len(userIds))
	totalScore := make([]int, len(userIds))
	totalPutt := make([]int, len(userIds))
	outTotalScore := make([]int, len(userIds))
	outTotalPutt := make([]int, len(userIds))
	inTotalScore := make([]int, len(userIds))
	inTotalPutt := make([]int, len(userIds))
	holes := make([]Hole, len(fields))

	var totalPar int
	var outTotalPar int
	var inTotalPar int
	for holeNum, field := range fields {
		holeIndex := holeNum - 1

		totalPar += field.Par
		if holeNum < 10 {
			outTotalPar += field.Par
		} else {
			inTotalPar += field.Par
		}
		score := make([]int, len(userIds))
		for playerIndex, userId := range userIds {
			scoreAHole := players[userId].Score[holeIndex]["total"].(int)
			puttAHole := players[userId].Score[holeIndex]["putt"].(int)

			score[playerIndex] = scoreAHole
			if holeNum < 10 {
				outTotalScore[playerIndex] += scoreAHole
				outTotalPutt[playerIndex] += puttAHole
			} else {
				inTotalScore[playerIndex] += scoreAHole
				inTotalPutt[playerIndex] += puttAHole
			}
			totalScore[playerIndex] += scoreAHole
			totalPutt[playerIndex] += puttAHole
			if holeIndex == 0 {
				member[playerIndex] = users[userId].Name
				apply[playerIndex] = players[userId].Apply
			}
		}
		holes[holeIndex] = Hole{
			Hole:  holeNum,
			Par:   field.Par,
			Yard:  field.Yard,
			Score: score,
		}
	}
	outSum := Sum{
		Par:   outTotalPar,
		Score: outTotalScore,
		Putt:  outTotalPutt,
	}
	inSum := Sum{
		Par:   inTotalPar,
		Score: inTotalScore,
		Putt:  inTotalPutt,
	}
	sum := Sum{
		Par:   totalPar,
		Score: totalScore,
		Putt:  totalPutt,
	}
	scoreViewSheet := ScoreViewSheet{
		Team:    teamName,
		Member:  member,
		UserIds: userIds,
		Apply:   apply,
		Hole:    holes,
		OutSum:  outSum,
		InSum:   inSum,
		Sum:     sum,
		Defined: teams[teamName].Defined,
	}

	return &scoreViewSheet, nil
}

func GetEntireScorePageData() (*EntireScore, error) {
	l.Output(nil, l.I_M_GetPage, l.Info)

	//rows[*][0] ホール rows[*][1] パー rows[*][2->n] PlayerName
	//rows[0] チーム名
	//rows[1] プレーヤー名
	//rows[2] IN
	//rows[3] ホール1
	//rows[4] ホール2
	//rows[5] ホール3
	//rows[6] ホール4
	//rows[7] ホール5
	//rows[8] ホール6
	//rows[9] ホール7
	//rows[10] ホール8
	//rows[11] ホール9
	//rows[12] OUT
	//rows[13] ホール1
	//rows[14] ホール2
	//rows[15] ホール3
	//rows[16] ホール4
	//rows[17] ホール5
	//rows[18] ホール6
	//rows[19] ホール7
	//rows[20] ホール8
	//rows[21] ホール9
	//rows[22] Gross
	//rows[23] Net
	//rows[24] 申請
	//rows[25] スコア差
	//rows[26] 順位
	rowSize := 27
	columnSize := len(players) + 2

	rows := make([][]string, rowSize)
	for i := 0; i < rowSize; i++ {
		if i == 0 {
			rows[i] = make([]string, len(teams)*2)
		} else {
			rows[i] = make([]string, columnSize)
		}
	}

	rows[1][0], rows[1][1] = "ホール", "パー"
	rows[2][0], rows[2][1] = "OUT", "ー"
	rows[12][0], rows[12][1] = "IN", "ー"
	rows[22][0], rows[22][1] = "Gross", "ー"
	rows[23][0], rows[23][1] = "Net", "ー"
	rows[24][0], rows[24][1] = "申請", "ー"
	rows[25][0], rows[25][1] = "スコア差", "ー"
	rows[26][0], rows[26][1] = "順位", "ー"

	var passedPlayerNum int
	var teamIndex int
	var headingColumnNum = 2
	var inColumNum = 1
	var outColumNum = 1
	var userIdArrayOrder []string

	var keys []string
	for k := range teams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	finalRankings := []FinalRanking{}
	for _, teamName := range keys {
		rows[0][teamIndex*2] = strconv.Itoa(len(teams[teamName].UserIds))
		rows[0][teamIndex*2+1] = "TEAM " + teamName
		userIds := teams[teamName].UserIds
		for playerIndex, userId := range userIds {
			userIdArrayOrder = append(userIdArrayOrder, userId)
			var gross int
			var net int
			userDataIndex := playerIndex + passedPlayerNum + headingColumnNum
			rows[1][userDataIndex] = users[userId].Name

			for holeNum, field := range fields {
				holeIndex := holeNum - 1
				halfHoleNum := len(fields) / 2
				holeRowNum := holeIndex + headingColumnNum + inColumNum
				inOutBorderline := headingColumnNum + inColumNum + halfHoleNum
				if holeRowNum >= inOutBorderline {
					holeRowNum = holeRowNum + outColumNum
				}
				if playerIndex == 0 {
					if holeNum > halfHoleNum {
						holeNum = holeNum - halfHoleNum
					}
					if field.Ignore {
						rows[holeRowNum][0] = "-i" + strconv.Itoa(holeNum)
					}
					rows[holeRowNum][0] = strconv.Itoa(holeNum)
					rows[holeRowNum][1] = strconv.Itoa(field.Par)
				}
				playerTotal := players[userId].Score[holeIndex]["total"].(int)
				gross += playerTotal
				if field.Ignore {
					net += field.Par
				} else {
					net += playerTotal
				}
				rows[holeRowNum][userDataIndex] = strconv.Itoa(playerTotal)
			}
			rows[22][userDataIndex] = strconv.Itoa(gross)
			rows[23][userDataIndex] = strconv.Itoa(net)
			rows[24][userDataIndex] = strconv.Itoa(players[userId].Apply)
			scoreDiff := players[userId].Apply - net
			if scoreDiff < 0 {
				scoreDiff = scoreDiff * -1
			}
			rows[25][userDataIndex] = strconv.Itoa(scoreDiff)
			finalRanking := FinalRanking{
				UserId:    userId,
				ScoreDiff: scoreDiff,
				Gross:     gross,
			}
			finalRankings = append(finalRankings, finalRanking)
		}
		passedPlayerNum += len(userIds)
		teamIndex++
	}

	sort.Sort(sortByRank(finalRankings))
	for userIdIndex, userId := range userIdArrayOrder {
		userDataIndex := userIdIndex + 2
		for rankingIndex, ranking := range finalRankings {
			if userId != ranking.UserId {
				continue
			}
			rows[26][userDataIndex] = strconv.Itoa(rankingIndex + 1)
		}
	}

	EntireScore := EntireScore{
		Rows: rows,
	}

	return &EntireScore, nil
}

func GetTimeLinePageData() (*TimeLine, error) {
	l.Output(nil, l.I_M_GetPage, l.Info)

	var timeLine TimeLine
	var tmpThreads []Thread
	var tmpThread Thread
	var tmpReactions []Reaction
	var tmpReaction Reaction

	threadsKeys := sortByDate{}
	for k, v := range threads {
		t := ThreadDate{k, v.CreatedAt}
		threadsKeys = append(threadsKeys, t)
	}

	sort.Sort(threadsKeys)

	for _, threadKey := range threadsKeys {
		threadId := threadKey.ThreadId
		if len(threads[threadId].ImgUrl) == 0 {
			continue
		}
		for _, reaction := range threads[threadId].Reactions {
			tmpReaction.Name = reaction["name"].(string)
			tmpReaction.ContentType = reaction["contenttype"].(int)
			tmpReaction.Content = reaction["content"].(string)
			tmpReaction.UserId = reaction["userid"].(string)
			tmpReaction.DateTime = reaction["datetime"].(string)
			tmpReactions = append(tmpReactions, tmpReaction)
		}
		tmpThread.ThreadId = threadId
		tmpThread.UserId = threads[threadId].UserId
		tmpThread.UserName = threads[threadId].UserName
		tmpThread.Msg = threads[threadId].Msg
		tmpThread.ImgUrl = threads[threadId].ImgUrl
		tmpThread.ColorCode = threads[threadId].ColorCode
		tmpThread.Positive = threads[threadId].Positive
		tmpThread.CreatedAt = threads[threadId].CreatedAt
		tmpThread.Reactions = tmpReactions
		tmpThreads = append(tmpThreads, tmpThread)
	}

	timeLine.Threads = tmpThreads

	return &timeLine, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
// 1 : Score Sort
// 2 : Passed Hole Sort
func (s sortByScore) Len() int      { return len(s) }
func (s sortByScore) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return s[i].Hole > s[j].Hole
	}
	return s[i].Score < s[j].Score
}

func (s sortByRank) Len() int      { return len(s) }
func (s sortByRank) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortByRank) Less(i, j int) bool {
	if s[i].ScoreDiff == s[j].ScoreDiff {
		return s[i].Gross < s[j].Gross
	}
	return s[i].ScoreDiff < s[j].ScoreDiff
}

func (s sortByDate) Len() int           { return len(s) }
func (s sortByDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByDate) Less(i, j int) bool { return s[i].CreatedAt > s[j].CreatedAt }
