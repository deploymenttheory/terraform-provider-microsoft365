package client

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

type testDeviceCodeTerminal struct {
	bytes.Buffer
	closed             bool
	writeErr, closeErr error
}

func (t *testDeviceCodeTerminal) Write(data []byte) (int, error) {
	if t.writeErr != nil {
		return 0, t.writeErr
	}
	return t.Buffer.Write(data)
}
func (t *testDeviceCodeTerminal) Close() error { t.closed = true; return t.closeErr }

func TestUnitDeviceCodePromptReachesTerminalWithoutPluginStdout(t *testing.T) {
	terminal := &testDeviceCodeTerminal{}
	callback := newDeviceCodePrompt(func() (io.WriteCloser, error) { return terminal, nil })
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = original; reader.Close(); writer.Close() }()
	err = callback(context.Background(), azidentity.DeviceCodeMessage{Message: "Open https://microsoft.com/devicelogin and enter UNIT-CODE-ONLY"})
	require.NoError(t, writer.Close())
	os.Stdout = original
	stdout, readErr := io.ReadAll(reader)
	require.NoError(t, readErr)
	require.NoError(t, err)
	require.Contains(t, terminal.String(), "UNIT-CODE-ONLY", "the prompt must reach the controlling terminal while plugin stdout is captured")
	require.Empty(t, string(stdout), "authentication instructions must not use the plugin stdout stream")
	require.True(t, terminal.closed)
}

func TestUnitDeviceCodePromptHeadlessFailsImmediately(t *testing.T) {
	callback := newDeviceCodePrompt(func() (io.WriteCloser, error) { return nil, os.ErrNotExist })
	err := callback(context.Background(), azidentity.DeviceCodeMessage{Message: "UNIT-CODE-ONLY"})
	require.ErrorIs(t, err, errDeviceCodeTerminalUnavailable)
	require.ErrorContains(t, err, "azure_cli")
	require.NotContains(t, err.Error(), "UNIT-CODE-ONLY", "the sign-in code must not leak through provider error logs")
}

func TestUnitDeviceCodePromptCancellationAndIOErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	opened := false
	callback := newDeviceCodePrompt(func() (io.WriteCloser, error) { opened = true; return &testDeviceCodeTerminal{}, nil })
	require.ErrorIs(t, callback(ctx, azidentity.DeviceCodeMessage{Message: "UNIT-CODE-ONLY"}), context.Canceled)
	require.False(t, opened)
	for _, tc := range []struct {
		name               string
		writeErr, closeErr error
	}{{name: "write", writeErr: io.ErrShortWrite}, {name: "close", closeErr: io.ErrClosedPipe}} {
		t.Run(tc.name, func(t *testing.T) {
			terminal := &testDeviceCodeTerminal{writeErr: tc.writeErr, closeErr: tc.closeErr}
			callback := newDeviceCodePrompt(func() (io.WriteCloser, error) { return terminal, nil })
			err := callback(context.Background(), azidentity.DeviceCodeMessage{Message: "UNIT-CODE-ONLY"})
			require.Error(t, err)
			require.True(t, terminal.closed)
		})
	}
}

func TestUnitDeviceCodePromptFallbackInstructions(t *testing.T) {
	terminal := &testDeviceCodeTerminal{}
	callback := newDeviceCodePrompt(func() (io.WriteCloser, error) { return terminal, nil })
	require.NoError(t, callback(context.Background(), azidentity.DeviceCodeMessage{UserCode: "UNIT-CODE-ONLY", VerificationURL: "https://microsoft.com/devicelogin"}))
	require.Contains(t, terminal.String(), "https://microsoft.com/devicelogin")
	require.Contains(t, terminal.String(), "UNIT-CODE-ONLY")
	require.ErrorIs(t, callback(context.Background(), azidentity.DeviceCodeMessage{}), errDeviceCodeInstructions)
}

// Run explicitly from a terminal to exercise the real OS console with synthetic
// instructions; default unit runs remain independent of an attached terminal.
func TestUnitDeviceCodePromptNativeTerminal(t *testing.T) {
	if os.Getenv("M365_TEST_DEVICE_CODE_TERMINAL") != "1" {
		t.Skip("native terminal probe requires an explicitly attached terminal")
	}
	callback := newDeviceCodePrompt(openDeviceCodeTerminal)
	require.NoError(t, callback(context.Background(), azidentity.DeviceCodeMessage{
		Message: "M365_NATIVE_TERMINAL_PROBE_OK (synthetic instructions; no authentication)",
	}))
}
