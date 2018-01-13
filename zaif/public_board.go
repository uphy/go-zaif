package zaif

func (p *PublicAPI) GetBoard(currencyPair string) (*Board, error) {
	type BoardData struct {
		Asks [][]float64
		Bids [][]float64
	}
	var boardData BoardData
	if err := p.getWithRetry("depth/"+currencyPair, &boardData); err != nil {
		return nil, err
	}
	return newBoard(convertBoardArray(boardData.Asks), convertBoardArray(boardData.Bids)), nil
}
