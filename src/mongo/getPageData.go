package mongo

import (
	"errors"
	"log"
	"sort"
	"strconv"
)

func GetIndexPageData() (*Index, error) {

	var teamArray []string
	for name, _ := range teams {
		teamArray = append(teamArray, name)
	}

	sort.Strings(teamArray)
	index := &Index{
		Team:   teamArray,
		Length: len(teams),
	}
	log.Println(index)
	return index, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {

	var leadersBoard LeadersBoard
	var userScore UserScore

	for _, player := range players {
		var totalPar, totalScore, passedHoleCnt int
		for holeIndex, playerScore := range player.Score {
			if playerScore["total"] != 0 {
				totalPar += fields[holeIndex].Par
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
	if len(holeString) == 0 {
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
		Team:   teamName,
		Hole:   holeNum,
		Member: member,
		Par:    field.Par,
		Yard:   field.Yard,
		Total:  total,
		Putt:   putt,
		Excnt:  excnt[teamName][holeNum],
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {

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
	rankMap := make(map[string]bool)

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
	rows[20][0], rows[20][1] = "Gross", "ー"
	rows[21][0], rows[21][1] = "Net", "ー"
	rows[22][0], rows[22][1] = "申請", "ー"
	rows[23][0], rows[23][1] = "スコア差", "ー"
	rows[24][0], rows[24][1] = "順位", "ー"

	var passedPlayerNum int
	var teamIndex int
	var headingColumnNum = 2
	var inColumNum = 1
	var outColumNum = 1
	for teamName, team := range teams {
		rows[0][teamIndex*2] = strconv.Itoa(len(team.UserIds))
		rows[0][teamIndex*2+1] = "TEAM " + teamName
		userIds := team.UserIds
		for playerIndex, userId := range userIds {
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
						rows[holeRowNum][0] = "-i"
					}
					rows[holeRowNum][0] += strconv.Itoa(holeNum)
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
			rows[20][userDataIndex] = strconv.Itoa(gross)
			rows[21][userDataIndex] = strconv.Itoa(net)
			rows[22][userDataIndex] = strconv.Itoa(players[userId].Apply)
			scoreDiff := players[userId].Apply - net
			if scoreDiff < 0 {
				scoreDiff = scoreDiff * -1
			}
			rows[23][userDataIndex] = strconv.Itoa(scoreDiff)
			rankMap[strconv.Itoa(scoreDiff)] = true
		}
		passedPlayerNum += len(userIds)
		teamIndex++
	}

	var rank []int
	for k := range rankMap {
		intK, _ := strconv.Atoi(k)
		rank = append(rank, intK)
	}

	sort.Ints(rank)
	for i := 0; i < len(rows[24])-2; i++ {
		userDataIndex := i + 2
		for j := 0; j < len(rank); j++ {
			if rows[23][userDataIndex] != strconv.Itoa(rank[j]) {
				continue
			} else {
				rows[24][userDataIndex] = strconv.Itoa(j + 1)
			}
		}
	}

	EntireScore := EntireScore{
		Rows: rows,
	}

	return &EntireScore, nil
}

func GetTimeLinePageData() (*TimeLine, error) {

	var timeLine TimeLine
	var tmpThreads []Thread
	var tmpThread Thread
	var tmpReactions []Reaction
	var tmpReaction Reaction

	for threadId, thread := range threads {
		for _, reaction := range thread.Reactions {
			tmpReaction.Name = reaction["name"].(string)
			tmpReaction.ContentType = reaction["contenttype"].(int)
			tmpReaction.Name = reaction["content"].(string)
			tmpReaction.DateTime = reaction["datetime"].(string)
			tmpReactions = append(tmpReactions, tmpReaction)
		}
		tmpThread.ThreadId = threadId
		tmpThread.Msg = thread.Msg
		tmpThread.ImgUrl = thread.ImgUrl
		tmpThread.ColorCode = thread.ColorCode
		tmpThread.Positive = thread.Positive
		tmpThread.CreatedAt = thread.CreatedAt
		tmpThread.Reactions = tmpReactions
		tmpThreads = append(tmpThreads, tmpThread)
	}

	timeLine.Threads = tmpThreads

	return &timeLine, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
