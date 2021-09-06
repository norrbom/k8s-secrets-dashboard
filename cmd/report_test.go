package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestReportRow(t *testing.T) {
	report := make(report, 0)
	env := "si1"
	var found bool
	var namespace string
	var secret string

	namespace = "namespace0"
	secret = "MY_METRICS_INFLUXDB_PASSWORD"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultInfluxDB, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultInfluxDB: %s", secret, report[namespace].VaultInfluxDB)
	}

	namespace = "namespace1"
	secret = "MY_VPN_SOLACE_PASSWORD"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultSolace, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultSolace: %s", secret, report[namespace].VaultSolace)
	}

	namespace = "namespace2"
	secret = "MY_VPN_SOLACE_USERNAME"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultSolace, secret); found {
		t.Errorf("Did not expect to find %s in reportRow.VaultSolace: %s", secret, report[namespace].VaultSolace)
	}

	namespace = "namespace3"
	secret = "MY_B2_DB_ORACLE_PASSWORD"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultOracle, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultOracle: %s", secret, report[namespace].VaultOracle)
	}

	namespace = "namespace4"
	secret = "MY_MYSQL_PASSWORD"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultMySQL, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultMySQL: %s", secret, report[namespace].VaultMySQL)
	}

	namespace = "namespace5"
	secret = "MY_COUCHBASE_PASSWORD"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].VaultCouchbase, secret); !found {
		t.Errorf("Expected to find %s in reportRow.VaultCouchbase: %s", secret, report[namespace].VaultCouchbase)
	}

	namespace = "namespace6"
	secret = "datasource.betstatspool.password"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretDatasource, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretDatasource: %s", secret, report[namespace].SealedSecretDatasource)
	}

	namespace = "namespace7"
	secret = "database.password"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretDatasource, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretDatasource: %s", secret, report[namespace].SealedSecretDatasource)
	}

	namespace = "namespace8"
	secret = "domainevents.solace.password"
	report[namespace] = NewReportRow(env, namespace, secret)
	report[namespace].addSecret("solace.password")
	if _, found = Find(report[namespace].SealedSecretSolace, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretSolace: %s", secret, report[namespace].SealedSecretSolace)
	}
	if _, found = Find(report[namespace].SealedSecretSolace, "solace.password"); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretSolace: %s", "solace.password", report[namespace].SealedSecretSolace)
	}

	namespace = "namespace9"
	secret = "influx.password"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretInfluxDB, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretInfluxDB: %s", secret, report[namespace].SealedSecretInfluxDB)
	}

	namespace = "namespace10"
	secret = "influx.Username"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretInfluxDB, secret); found {
		t.Errorf("Did not expect to find %s in reportRow.SealedSecretInfluxDB: %s", secret, report[namespace].SealedSecretInfluxDB)
	}

	namespace = "namespace11"
	secret = "api.password"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].Unknown, secret); !found {
		t.Errorf("Expected to find %s in reportRow.Unknown: %s", secret, report[namespace].Unknown)
	}

	namespace = "namespace12"
	secret = "curity.oauth.client-secret"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretCurity, secret); !found {
		t.Errorf("Expected to find %s in reportRow.SealedSecretCurity: %s", secret, report[namespace].SealedSecretCurity)
	}

	namespace = "namespace13"
	secret = "security.password.kindred.osvc"
	report[namespace] = NewReportRow(env, namespace, secret)
	if _, found = Find(report[namespace].SealedSecretCurity, secret); found {
		t.Errorf("Did not expect to find %s in reportRow.SealedSecretCurity: %s", secret, report[namespace].SealedSecretCurity)
	}

	// Test progress logic
	namespace = "progress"
	report[namespace] = NewReportRow(env, namespace, "api.service.aws.kindred.com")
	if report[namespace].isInProgress() != false {
		t.Errorf("Expected isInProgress to be false for: %+v", report[namespace])
	}
	if report[namespace].inScope() != false {
		t.Errorf("Expected inScope to be false for: %+v", report[namespace])
	}
	report[namespace].addSecret("influx")
	if report[namespace].inScope() != true {
		t.Errorf("Expected inScope to be true for: %+v", report[namespace])
	}
	report[namespace].addSecret("_MYSQL_PASSWORD")
	report[namespace].setProgress()
	if report[namespace].isCompleted() != false {
		t.Errorf("Expected isCompleted to be false for: %+v", report[namespace])
	}
	if report[namespace].isInProgress() != true || report[namespace].Progress != "in-progress" {
		t.Errorf("Expected Progress to be \"in-progress\" for: %+v", report[namespace])
	}
	namespace = "progress2"
	report[namespace] = NewReportRow(env, namespace, "_MYSQL_PASSWORD")
	report[namespace].addSecret("_INFLUXDB_PASSWORD")
	report[namespace].addSecret("api.service.aws.kindred.com")
	report[namespace].setProgress()
	if report[namespace].isCompleted() != true || report[namespace].Progress != "completed" {
		t.Errorf("Expected isCompleted to be true for: %+v", report[namespace])
	}

	environments := []string{"DATA", "PROD"}
	env_sorted := []string{"PROD", "DATA"}
	sort.Sort(byImportance(environments))
	if !reflect.DeepEqual(environments, env_sorted) {
		t.Errorf("Expected: " + strings.Join(environments, ",") + " to be equal to: " + strings.Join(env_sorted, ","))
	}

	environments = []string{"DATA.PT1", "DATA.PROD1", "DATA.SI1"}
	env_sorted = []string{"DATA.PROD1", "DATA.PT1", "DATA.SI1"}
	sort.Sort(byImportance(environments))
	if !reflect.DeepEqual(environments, env_sorted) {
		t.Errorf("Expected: " + strings.Join(environments, ",") + " to be equal to: " + strings.Join(env_sorted, ","))
	}

	environments = []string{"PT1", "STAGE1.US1", "DATA.PT1", "DATA.PROD1", "PROD1", "SI1", "DATA.SI1"}
	env_sorted = []string{"PROD1", "STAGE1.US1", "PT1", "SI1", "DATA.PROD1", "DATA.PT1", "DATA.SI1"}
	sort.Sort(byImportance(environments))
	if !reflect.DeepEqual(environments, env_sorted) {
		t.Errorf("Expected: " + strings.Join(environments, ",") + " to be equal to: " + strings.Join(env_sorted, ","))
	}
}
