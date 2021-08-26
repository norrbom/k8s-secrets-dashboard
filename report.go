package main

import (
	"strings"
)

type reportRow struct {
	environment            string
	namespace              string
	Unknown                []string
	VaultSolace            []string
	VaultInfluxDB          []string
	VaultOracle            []string
	VaultMySQL             []string
	VaultCouchbase         []string
	SealedSecretDatasource []string
	SealedSecretInfluxDB   []string
	SealedSecretSolace     []string
	SealedSecretCouchbase  []string
	Progress               string
}

/*
factory method for reportRow, usage rr := NewReportRow("prod1", "service")
*/
func NewReportRow(environment string, namespace string, secret string) *reportRow {
	rr := new(reportRow)
	rr.environment = environment
	rr.namespace = namespace
	rr.addSecret(secret)
	return rr
}

type report map[string]*reportRow

func (rr *reportRow) inScope() bool {
	if len(rr.SealedSecretCouchbase) > 0 || len(rr.SealedSecretDatasource) > 0 || len(rr.SealedSecretInfluxDB) > 0 || len(rr.SealedSecretSolace) > 0 ||
		len(rr.VaultSolace) > 0 || len(rr.VaultInfluxDB) > 0 || len(rr.VaultOracle) > 0 || len(rr.VaultMySQL) > 0 || len(rr.VaultCouchbase) > 0 {
		return true
	}
	return false
}

func (rr *reportRow) setProgress() {
	switch {
	case (*rr).isCompleted():
		(*rr).Progress = "completed"
	case (*rr).isInProgress():
		(*rr).Progress = "in-progress"
	default:
		(*rr).Progress = "not-started"
	}
}

func (rr *reportRow) isInProgress() bool {
	vault_count := len(rr.VaultInfluxDB) + len(rr.VaultMySQL) + len(rr.VaultOracle) + len(rr.VaultCouchbase) + len(rr.VaultSolace)
	ss_count := len(rr.SealedSecretInfluxDB) + len(rr.SealedSecretDatasource) + len(rr.SealedSecretCouchbase) + len(rr.SealedSecretSolace)
	if ss_count > 0 && vault_count > 0 {
		return true
	}
	return false
}

func (rr *reportRow) isCompleted() bool {
	if len(rr.SealedSecretDatasource) <= (len(rr.VaultMySQL)+len(rr.VaultOracle)) &&
		len(rr.SealedSecretSolace) <= len(rr.VaultSolace) &&
		len(rr.SealedSecretInfluxDB) <= len(rr.VaultInfluxDB) &&
		len(rr.SealedSecretCouchbase) <= len(rr.VaultCouchbase) {
		return true
	}
	return false
}

func (rr *reportRow) addSecret(secret string) {
	switch {
	case strings.Index(strings.ToLower(secret), "username") > -1:
		break // we are not intrested in usernames
	case strings.HasSuffix(secret, "_INFLUXDB_PASSWORD"):
		(*rr).VaultInfluxDB = append((*rr).VaultInfluxDB, secret)
	case strings.HasSuffix(secret, "_SOLACE_PASSWORD"):
		(*rr).VaultSolace = append((*rr).VaultSolace, secret)
	case strings.HasSuffix(secret, "_ORACLE_PASSWORD"):
		(*rr).VaultOracle = append((*rr).VaultOracle, secret)
	case strings.HasSuffix(secret, "_MYSQL_PASSWORD"):
		(*rr).VaultMySQL = append((*rr).VaultMySQL, secret)
	case strings.HasSuffix(secret, "_COUCHBASE_PASSWORD"):
		(*rr).VaultCouchbase = append((*rr).VaultCouchbase, secret)
	case strings.Index(secret, "datasource") > -1 || strings.Index(secret, "database") > -1:
		(*rr).SealedSecretDatasource = append((*rr).SealedSecretDatasource, secret)
	case strings.Index(secret, "influx") > -1:
		(*rr).SealedSecretInfluxDB = append((*rr).SealedSecretInfluxDB, secret)
	case strings.Index(secret, "solace") > -1:
		(*rr).SealedSecretSolace = append((*rr).SealedSecretSolace, secret)
	case strings.Index(secret, "couch") > -1:
		(*rr).SealedSecretCouchbase = append((*rr).SealedSecretCouchbase, secret)
	default:
		(*rr).Unknown = append((*rr).Unknown, secret)
	}
}
