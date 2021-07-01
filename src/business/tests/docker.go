package tests

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

type Container struct {
	ID   string
	Host string
}

func startContainer(t *testing.T, image string, port string, args ...string) *Container {
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, image)

	cmd := exec.Command("docker", arg...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start the container: %s: %v", image, err)
	}

	id := out.String()[:12]

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect the container: %s: %v", id, err)
	}

	var doc []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	ip, randPort := extractRandIPPort(t, doc, port)

	c := Container{
		ID:   id,
		Host: net.JoinHostPort(ip, randPort),
	}

	t.Logf("Image: %s", image)
	t.Logf("ContainerID: %s", c.ID)
	t.Logf("Host: %s", c.Host)

	return &c
}

func stopContainer(t *testing.T, id string) {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}

	t.Log("Stopped:", id)

	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}

	t.Log("Removed:", id)
}

func dumpContainerLogs(t *testing.T, id string) {
	out, err := exec.Command("docker", "logs", id).CombinedOutput()

	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}

	t.Logf("Logs for %s:\n%s", id, out)
}

func extractRandIPPort(t *testing.T, doc []map[string]interface{}, port string) (string, string) {
	nw, exists := doc[0]["NetworkSettings"]
	if !exists {
		t.Fatal("could not get a network settings")
	}

	ports, exists := nw.(map[string]interface{})["Ports"]
	if !exists {
		t.Fatal("could not get a port settings")
	}

	tcp, exists := ports.(map[string]interface{})[port]
	if !exists {
		t.Fatalf("could not get port %s settings", port)
	}

	list, exists := tcp.([]interface{})
	if !exists {
		t.Fatalf("could not get network port %s list settings", port)
	}

	if len(list) == 0 {
		t.Fatalf("could not get network port %s list settings", port)
	}

	data, exists := list[0].(map[string]interface{})
	if !exists {
		t.Fatalf("could not get network port %s list settings", port)
	}

	return data["HostIp"].(string), data["HostPort"].(string)
}
