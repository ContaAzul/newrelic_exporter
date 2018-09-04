package newrelic

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListKeyTransactions(t *testing.T) {
	var require = require.New(t)
	server := responseStub(t, "key_transactions.json", http.StatusOK)
	defer server.Close()

	c := NewClient(server.URL, "key")
	result, err := c.ListKeyTransactions()

	require.NoError(err)
	require.Len(result, 2)
	require.Equal(KeyTransaction{
		ID:              11694,
		Name:            "/dashboards (POST)",
		TransactionName: "/dashboards (POST)",
		Reporting:       true,
		ApplicationSummary: ApplicationSummary{
			InstanceCount: 0,
			ResponseTime:  175,
			Throughput:    223,
			ErrorRate:     0.02,
			ApdexTarget:   0.21,
			ApdexScore:    0.87,
		},
	}, result[0])
}
