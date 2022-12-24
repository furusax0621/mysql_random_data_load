package tableparser

import (
	"testing"
	"time"

	"github.com/furusax0621/mysql_random_data_load/testutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
)

func TestParse(t *testing.T) {
	db := testutils.GetMySQLConnection(t)
	v := testutils.GetMinorVersion(t, db)
	var want *Table

	// Patch part of version is stripped by GetMinorVersion so for these test
	// it is .0
	sampleFiles := map[string]string{
		"5.6.0": "table001.json",
		"5.7.0": "table002.json",
		"8.0.0": "table003.json",
	}
	sampleFile, ok := sampleFiles[v.String()]
	if !ok {
		t.Fatalf("Unknown MySQL version %s", v.String())
	}

	table, err := NewTable(db, "sakila", "film")
	if err != nil {
		t.Error(err)
	}
	if testutils.UpdateSamples() {
		testutils.WriteJson(t, sampleFile, table)
	}
	testutils.LoadJson(t, sampleFile, &want)

	testutils.Equals(t, want, table)
}

func TestGetIndexes(t *testing.T) {
	db := testutils.GetMySQLConnection(t)
	want := make(map[string]Index)
	testutils.LoadJson(t, "indexes.json", &want)

	idx, err := getIndexes(db, "sakila", "film_actor")
	if testutils.UpdateSamples() {
		testutils.WriteJson(t, "indexes.json", idx)
	}
	testutils.Ok(t, err)
	testutils.Equals(t, want, idx)
}

func TestGetTriggers(t *testing.T) {
	db := testutils.GetMySQLConnection(t)
	want := []Trigger{}
	v572, _ := version.NewVersion("5.7.2")
	v800, _ := version.NewVersion("8.0.0")

	sampleFile := "trigers-8.0.0.json"
	if testutils.GetVersion(t, db).LessThan(v800) {
		sampleFile = "trigers-5.7.2.json"
	}
	if testutils.GetVersion(t, db).LessThan(v572) {
		sampleFile = "trigers-5.7.1.json"
	}

	testutils.LoadJson(t, sampleFile, &want)

	triggers, err := getTriggers(db, "sakila", "rental")
	// fake timestamp to make it constant/testeable
	triggers[0].Created.Time = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	if testutils.UpdateSamples() {
		logrus.Info("Updating sample file: " + sampleFile)
		testutils.WriteJson(t, sampleFile, triggers)
	}
	testutils.Ok(t, err)
	testutils.Equals(t, want, triggers)
}
