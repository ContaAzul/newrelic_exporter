package newrelic

// KeyTransaction represents a New Relic key transaction.
type KeyTransaction struct {
	ID                 int64              `json:"id"`
	Name               string             `json:"name"`
	TransactionName    string             `json:"transaction_name"`
	Reporting          bool               `json:"reporting"`
	ApplicationSummary ApplicationSummary `json:"application_summary"`
}

type listKeyTransactionsResponse struct {
	KeyTransactions []KeyTransaction `json:"key_transactions"`
}

// ListKeyTransactions returns a paginated list of the key transactions associated with your
// New Relic account. The time range for summary data is the last 10 minutes.
func (c *Client) ListKeyTransactions() ([]KeyTransaction, error) {
	path := "v2/key_transactions.json"
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var response listKeyTransactionsResponse
	_, err = c.do(req, &response)
	return response.KeyTransactions, err
}
