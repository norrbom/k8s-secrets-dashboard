package main

import (
	"testing"
)

func TestReportRow(t *testing.T) {
	report := make(report, 0)
	var found bool
	var namespace string
	var secret string

	namespace = "namespace0"
	secret = "MY_METRICS_INFLUXDB_PASSWORD"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultInfluxDB, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultInfluxDB: %s", secret, report[namespace].VaultInfluxDB)
	}

	namespace = "namespace1"
	secret = "MY_VPN_SOLACE_PASSWORD"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultSolace, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultSolace: %s", secret, report[namespace].VaultSolace)
	}

	namespace = "namespace2"
	secret = "MY_VPN_SOLACE_USERNAME"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultSolace, secret); found {
		t.Errorf("Did not expect to find %s in reportRow.VaultSolace: %s", secret, report[namespace].VaultSolace)
	}

	namespace = "namespace3"
	secret = "MY_B2_DB_ORACLE_PASSWORD"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultOracle, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultOracle: %s", secret, report[namespace].VaultOracle)
	}

	namespace = "namespace4"
	secret = "MY_MYSQL_PASSWORD"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultMySQL, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultMySQL: %s", secret, report[namespace].VaultMySQL)
	}

	namespace = "namespace5"
	secret = "MY_COUCHBASE_PASSWORD"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].VaultCouchbase, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultCouchbase: %s", secret, report[namespace].VaultCouchbase)
	}

	namespace = "namespace6"
	secret = "datasource.betstatspool.password"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].SealedSecretDatasource, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretDatasource: %s", secret, report[namespace].SealedSecretDatasource)
	}

	namespace = "namespace7"
	secret = "database.password"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].SealedSecretDatasource, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretDatasource: %s", secret, report[namespace].SealedSecretDatasource)
	}

	namespace = "namespace8"
	secret = "domainevents.solace.password"
	report.add(namespace, secret)
	report.add(namespace, "solace.password")
	if _, found = Find(report[namespace].SealedSecretSolace, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretSolace: %s", secret, report[namespace].SealedSecretSolace)
	}
	if _, found = Find(report[namespace].SealedSecretSolace, "solace.password"); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretSolace: %s", "solace.password", report[namespace].SealedSecretSolace)
	}

	namespace = "namespace9"
	secret = "influx.password"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].SealedSecretInfluxDB, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretInfluxDB: %s", secret, report[namespace].SealedSecretInfluxDB)
	}

	namespace = "namespace10"
	secret = "influx.Username"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].SealedSecretInfluxDB, secret); found {
		t.Errorf("Did not expect to find %s in reportRow.SealedSecretInfluxDB: %s", secret, report[namespace].SealedSecretInfluxDB)
	}

	namespace = "namespace11"
	secret = "api.password"
	report.add(namespace, secret)
	if _, found = Find(report[namespace].Unknown, secret); !found {
		t.Errorf("Expected to find %s in reportRow.Unknown: %s", secret, report[namespace].Unknown)
	}

	// Test progress logic
	namespace = "progress"
	report.add(namespace, "api.service.aws.kindred.com")
	if report[namespace].isInProgress() != false {
		t.Errorf("Expected isInProgress to be false for: %s", report[namespace])
	}
	if report[namespace].inScope() != false {
		t.Errorf("Expected inScope to be false for: %s", report[namespace])
	}
	report.add(namespace, "influx")
	if report[namespace].inScope() != true {
		t.Errorf("Expected inScope to be true for: %s", report[namespace])
	}
	report.add(namespace, "_MYSQL_PASSWORD")
	rrInProgress := report[namespace]
	if rrInProgress.isCompleted() != false {
		t.Errorf("Expected isCompleted to be false for: %s", report[namespace])
	}
	if report[namespace].isInProgress() != true {
		t.Errorf("Expected isInProgress to be true for: %s", report[namespace])
	}
	namespace = "progress2"
	report.add(namespace, "_MYSQL_PASSWORD")
	report.add(namespace, "_INFLUXDB_PASSWORD")
	report.add(namespace, "api.service.aws.kindred.com")
	rrComplete := report[namespace]
	if rrComplete.isCompleted() != true {
		t.Errorf("Expected isCompleted to be true for: %s", report[namespace])
	}
}
