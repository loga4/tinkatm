package atm

import "time"

type TinkResponse struct {
	TrackingID string `json:"trackingId"`
	Payload    struct {
		Hash   string `json:"hash"`
		Zoom   int    `json:"zoom"`
		Bounds struct {
			BottomLeft struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"bottomLeft"`
			TopRight struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"topRight"`
		} `json:"bounds"`
		Clusters []struct {
			ID     string `json:"id"`
			Hash   string `json:"hash"`
			Bounds struct {
				BottomLeft struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"bottomLeft"`
				TopRight struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"topRight"`
			} `json:"bounds"`
			Center struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"center"`
			Points []struct {
				ID    string `json:"id"`
				Brand struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					LogoFile    string `json:"logoFile"`
					RoundedLogo bool   `json:"roundedLogo"`
				} `json:"brand"`
				PointType string `json:"pointType"`
				Location  struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Address string   `json:"address"`
				Phone   []string `json:"phone"`
				Limits  []struct {
					Currency      string `json:"currency"`
					Max           int    `json:"max"`
					Denominations []int  `json:"denominations"`
					Amount        int    `json:"amount"`
				} `json:"limits"`
				WorkPeriods []struct {
					OpenDay   int    `json:"openDay"`
					OpenTime  string `json:"openTime"`
					CloseDay  int    `json:"closeDay"`
					CloseTime string `json:"closeTime"`
				} `json:"workPeriods"`
				InstallPlace string `json:"installPlace"`
				AtmInfo      struct {
					Available  bool `json:"available"`
					IsTerminal bool `json:"isTerminal"`
					Statuses   struct {
						CriticalFailure       bool `json:"criticalFailure"`
						QrOperational         bool `json:"qrOperational"`
						NfcOperational        bool `json:"nfcOperational"`
						CardReaderOperational bool `json:"cardReaderOperational"`
						CashInAvailable       bool `json:"cashInAvailable"`
					} `json:"statuses"`
					Limits []struct {
						Currency                string        `json:"currency"`
						Amount                  int           `json:"amount"`
						WithdrawMaxAmount       int           `json:"withdrawMaxAmount"`
						DepositionMaxAmount     int           `json:"depositionMaxAmount"`
						DepositionMinAmount     int           `json:"depositionMinAmount"`
						WithdrawDenominations   []interface{} `json:"withdrawDenominations"`
						DepositionDenominations []int         `json:"depositionDenominations"`
						OverTrustedLimit        bool          `json:"overTrustedLimit"`
					} `json:"limits"`
				} `json:"atmInfo"`
			} `json:"points"`
		} `json:"clusters"`
	} `json:"payload"`
	Time   time.Time `json:"time"`
	Status string    `json:"status"`
}
