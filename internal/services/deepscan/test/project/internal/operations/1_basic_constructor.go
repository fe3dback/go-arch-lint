package operations

type (
	myFetcher interface {
		Fetch()
	}

	PublicFetcherForDI = myFetcher

	myFetcherAlias = myFetcher
	myFetcherNamed myFetcher

	processor1 struct {
		fetcher myFetcher
	}
)

func NewProcessorBasic1(myFetcher myFetcher) *processor1 {
	return &processor1{
		fetcher: myFetcher,
	}
}

func NewProcessorNames1(myFetcher myFetcherNamed) *processor1 {
	return &processor1{
		fetcher: myFetcher,
	}
}

func NewProcessorAlias1(myFetcher myFetcherAlias) *processor1 {
	return &processor1{
		fetcher: myFetcher,
	}
}

func NewProcessorBasicDual1(fetcher1 myFetcherAlias, fetcher2 myFetcherNamed) *processor1 {
	_ = fetcher2
	return &processor1{
		fetcher: fetcher1,
	}
}

func NewProcessorBasicSpreadNames1(fetcher1, fetcher2 myFetcherNamed) *processor1 {
	_ = fetcher2
	return &processor1{
		fetcher: fetcher1,
	}
}

func NewProcessorBasicSpreadNamesAnonim1(fetcher1, _, _, fetcher4 myFetcherNamed) *processor1 {
	_ = fetcher4
	return &processor1{
		fetcher: fetcher1,
	}
}

func NewProcessorBasicSpreadTypes1(fetchers ...myFetcherNamed) *processor1 {
	return &processor1{
		fetcher: fetchers[0],
	}
}

func NewProcessorBasicSlice1(fetchers []PublicFetcherForDI) *processor1 {
	return &processor1{
		fetcher: fetchers[0],
	}
}

func NewProcessorBasicArray1(fetchers [5]PublicFetcherForDI) *processor1 {
	return &processor1{
		fetcher: fetchers[0],
	}
}
