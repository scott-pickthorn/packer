package instance

import (
	"github.com/mitchellh/packer/packer"
	"io/ioutil"
	"os"
	"testing"
)

func testConfig() map[string]interface{} {
	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		panic(err)
	}

	return map[string]interface{}{
		"account_id":       "foo",
		"instance_type":    "m1.small",
		"region":           "us-east-1",
		"s3_bucket":        "foo",
		"source_ami":       "foo",
		"ssh_username":     "bob",
		"x509_cert_path":   tf.Name(),
		"x509_key_path":    tf.Name(),
		"x509_upload_path": "/foo",
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{}
	raw = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilderPrepare_AccountId(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["account_id"] = ""
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["account_id"] = "foo"
	err = b.Prepare(config)
	if err != nil {
		t.Errorf("err: %s", err)
	}

	config["account_id"] = "0123-0456-7890"
	err = b.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if b.config.AccountId != "012304567890" {
		t.Errorf("should strip hyphens: %s", b.config.AccountId)
	}
}

func TestBuilderPrepare_BundleDestination(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["bundle_destination"] = ""
	err := b.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if b.config.BundleDestination != "/tmp" {
		t.Fatalf("bad: %s", b.config.BundleDestination)
	}
}

func TestBuilderPrepare_BundlePrefix(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["bundle_prefix"] = ""
	err := b.Prepare(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if b.config.BundlePrefix != "image" {
		t.Fatalf("bad: %s", b.config.BundlePrefix)
	}
}

func TestBuilderPrepare_InvalidKey(t *testing.T) {
	var b Builder
	config := testConfig()

	// Add a random key
	config["i_should_not_be_valid"] = true
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_S3Bucket(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["s3_bucket"] = ""
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["s3_bucket"] = "foo"
	err = b.Prepare(config)
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestBuilderPrepare_X509CertPath(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["x509_cert_path"] = ""
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["x509_cert_path"] = "i/am/a/file/that/doesnt/exist"
	err = b.Prepare(config)
	if err == nil {
		t.Error("should have error")
	}

	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("error tempfile: %s", err)
	}
	defer os.Remove(tf.Name())

	config["x509_cert_path"] = tf.Name()
	err = b.Prepare(config)
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}
}

func TestBuilderPrepare_X509KeyPath(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["x509_key_path"] = ""
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}

	config["x509_key_path"] = "i/am/a/file/that/doesnt/exist"
	err = b.Prepare(config)
	if err == nil {
		t.Error("should have error")
	}

	tf, err := ioutil.TempFile("", "packer")
	if err != nil {
		t.Fatalf("error tempfile: %s", err)
	}
	defer os.Remove(tf.Name())

	config["x509_key_path"] = tf.Name()
	err = b.Prepare(config)
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}
}

func TestBuilderPrepare_X509UploadPath(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["x509_upload_path"] = ""
	err := b.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}
}
