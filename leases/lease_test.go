package leases

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/teamhephy/k8s-claimer/testutil"
)

var (
	clusterName = "my test cluster"
)

func TestParseLease(t *testing.T) {
	tme := time.Now()
	tStr := tme.Format(TimeFormat)
	jsonStr := testutil.LeaseJSON(clusterName, tme, TimeFormat)
	lease, err := ParseLease(jsonStr)
	assert.NoErr(t, err)
	assert.Equal(t, lease.ClusterName, clusterName, "cluster name")
	assert.Equal(t, lease.LeaseExpirationTime, tStr, "expiration time string")
}

func TestNewLease(t *testing.T) {
	tme := time.Now()
	l := NewLease(clusterName, tme)
	assert.Equal(t, l.ClusterName, clusterName, "cluster name")
	assert.Equal(t, l.LeaseExpirationTime, tme.Format(TimeFormat), "lease expiration time")
}

func TestExpirationTime(t *testing.T) {
	tme := time.Now()
	l := NewLease(clusterName, tme)
	ex, err := l.ExpirationTime()
	assert.NoErr(t, err)
	assert.Equal(t, ex.Format(TimeFormat), tme.Format(TimeFormat), "expiration time")

	l = &Lease{ClusterName: clusterName, LeaseExpirationTime: "invalid expiration time"}
	ex, err = l.ExpirationTime()
	assert.True(t, err != nil, "error wasn't returned when it should have been")
	assert.True(t, ex.IsZero(), "non-zero time was returned when it should have been zero")
}
