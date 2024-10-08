package dto

type EsReIndexInfo struct {
	EsConnect int `json:"es_connect"`
	UrlValues struct {
		Timeout             int    `json:"timeout"`
		RequestsPerSecond   int    `json:"requests_per_second"`
		Slices              int    `json:"slices"`
		Scroll              string `json:"scroll"`
		WaitForActiveShards string `json:"wait_for_active_shards"`
		Refresh             *bool  `json:"refresh"`
		WaitForCompletion   *bool  `json:"wait_for_completion"`
	} `json:"url_values"`
	Body map[string]interface{} `json:"body"`
}
