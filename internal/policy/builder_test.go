package policy_test

import (
	"testing"

	"github.com/brobridge/sentinel/internal/policy"
	"sigs.k8s.io/yaml"
)

func TestBuildProcessPolicy(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:    "block-shells",
		Process: []policy.ProcessRule{{Binaries: []string{"/bin/bash", "/bin/sh"}}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got.Metadata.Name != "block-shells" {
		t.Errorf("name = %q, want block-shells", got.Metadata.Name)
	}
	if got.Kind != "TracingPolicy" {
		t.Errorf("kind = %q, want TracingPolicy", got.Kind)
	}
	if len(got.Spec.KProbes) != 1 {
		t.Fatalf("kprobes len = %d, want 1", len(got.Spec.KProbes))
	}
	kp := got.Spec.KProbes[0]
	if kp.Call != "sys_execve" {
		t.Errorf("call = %q, want sys_execve", kp.Call)
	}
	if !kp.Syscall {
		t.Error("syscall should be true for process rule")
	}
	if kp.Selectors[0].MatchActions[0].Action != "Post" {
		t.Errorf("action = %q, want Post", kp.Selectors[0].MatchActions[0].Action)
	}
}

func TestBuildFilePolicy_Write(t *testing.T) {
	input := policy.PolicyFormInput{
		Name: "watch-etc",
		File: []policy.FileRule{{Paths: []string{"/etc/passwd"}, Operation: "write"}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if len(got.Spec.KProbes) != 1 {
		t.Fatalf("kprobes len = %d, want 1", len(got.Spec.KProbes))
	}
	if got.Spec.KProbes[0].Call != "sys_write" {
		t.Errorf("call = %q, want sys_write", got.Spec.KProbes[0].Call)
	}
}

func TestBuildFilePolicy_Read(t *testing.T) {
	input := policy.PolicyFormInput{
		Name: "watch-read",
		File: []policy.FileRule{{Paths: []string{"/etc/shadow"}, Operation: "read"}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got.Spec.KProbes[0].Call != "sys_read" {
		t.Errorf("call = %q, want sys_read", got.Spec.KProbes[0].Call)
	}
}

func TestBuildFilePolicy_UnknownOperation(t *testing.T) {
	input := policy.PolicyFormInput{
		Name: "bad",
		File: []policy.FileRule{{Paths: []string{"/tmp"}, Operation: "delete"}},
	}

	_, err := policy.Build(input, policy.ActionPost)
	if err == nil {
		t.Error("Build() should return error for unknown operation")
	}
}

func TestBuildNetworkPolicy(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:    "block-external",
		Network: []policy.NetworkRule{{Protocol: "TCP", CIDR: "10.0.0.0/8", Port: 8080}},
	}

	got, err := policy.Build(input, policy.ActionSigkill)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	kp := got.Spec.KProbes[0]
	if kp.Call != "tcp_connect" {
		t.Errorf("call = %q, want tcp_connect", kp.Call)
	}
	if kp.Selectors[0].MatchArgs[0].Values[0] != "10.0.0.0/8:8080" {
		t.Errorf("value = %q, want 10.0.0.0/8:8080", kp.Selectors[0].MatchArgs[0].Values[0])
	}
	if kp.Selectors[0].MatchActions[0].Action != "Sigkill" {
		t.Errorf("action = %q, want Sigkill", kp.Selectors[0].MatchActions[0].Action)
	}
}

func TestBuildNetworkPolicy_NoPort(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:    "any-port",
		Network: []policy.NetworkRule{{Protocol: "TCP", CIDR: "192.168.0.0/16"}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got.Spec.KProbes[0].Selectors[0].MatchArgs[0].Values[0] != "192.168.0.0/16" {
		t.Errorf("expected CIDR only, got %q", got.Spec.KProbes[0].Selectors[0].MatchArgs[0].Values[0])
	}
}

func TestBuildNamespacedPolicy(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:      "ns-policy",
		Namespace: "production",
		Process:   []policy.ProcessRule{{Binaries: []string{"/bin/bash"}}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got.Kind != "TracingPolicyNamespaced" {
		t.Errorf("kind = %q, want TracingPolicyNamespaced", got.Kind)
	}
	if got.Metadata.Namespace != "production" {
		t.Errorf("namespace = %q, want production", got.Metadata.Namespace)
	}
}

func TestBuildPodSelector(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:        "scoped",
		PodSelector: map[string]string{"app": "myapp"},
		Process:     []policy.ProcessRule{{Binaries: []string{"/bin/sh"}}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if got.Spec.PodSelector == nil {
		t.Fatal("expected non-nil PodSelector")
	}
	if got.Spec.PodSelector.MatchLabels["app"] != "myapp" {
		t.Errorf("podSelector.app = %q, want myapp", got.Spec.PodSelector.MatchLabels["app"])
	}
}

func TestBuildMultipleRules(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:    "multi",
		Process: []policy.ProcessRule{{Binaries: []string{"/bin/bash"}}},
		File:    []policy.FileRule{{Paths: []string{"/etc"}, Operation: "open"}},
		Network: []policy.NetworkRule{{Protocol: "TCP", CIDR: "0.0.0.0/0"}},
	}

	got, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	if len(got.Spec.KProbes) != 3 {
		t.Errorf("kprobes len = %d, want 3", len(got.Spec.KProbes))
	}
}

func TestBuildYAMLRoundtrip(t *testing.T) {
	input := policy.PolicyFormInput{
		Name:    "roundtrip",
		Process: []policy.ProcessRule{{Binaries: []string{"/bin/sh"}}},
	}

	tp, err := policy.Build(input, policy.ActionPost)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}
	b, err := yaml.Marshal(tp)
	if err != nil {
		t.Fatalf("yaml.Marshal() error: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty YAML")
	}

	var out policy.TracingPolicy
	if err := yaml.Unmarshal(b, &out); err != nil {
		t.Fatalf("yaml.Unmarshal() error: %v", err)
	}
	if out.Metadata.Name != "roundtrip" {
		t.Errorf("name after roundtrip = %q, want roundtrip", out.Metadata.Name)
	}
}
