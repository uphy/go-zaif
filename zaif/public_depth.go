package zaif

func (p *PublicAPI) GetDepth(currencyPair string) (*Depth, error) {
	type DepthData struct {
		Asks [][]float64
		Bids [][]float64
	}
	var depthData DepthData
	if err := p.getWithRetry("depth/"+currencyPair, &depthData); err != nil {
		return nil, err
	}
	return newDepth(convertDepthArray(depthData.Asks), convertDepthArray(depthData.Bids)), nil
}
