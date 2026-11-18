package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	errDeviceCodeTerminalUnavailable = errors.New(
		"device_code authentication requires an interactive terminal; run Terraform or OpenTofu in a terminal, or sign in with 'az login --use-device-code' and use auth_method = \"azure_cli\"",
	)
	errDeviceCodeInstructions = errors.New("device code authentication instructions missing from Microsoft Entra ID")
)

func openDeviceCodeTerminal() (io.WriteCloser, error) {
	name := "/dev/tty"
	if runtime.GOOS == "windows" {
		name = "CONOUT$"
	}
	terminal, err := os.OpenFile(name, os.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("open device code terminal: %w", err)
	}
	return terminal, nil
}

func newDeviceCodePrompt(openTerminal func() (io.WriteCloser, error)) func(context.Context, azidentity.DeviceCodeMessage) error {
	return func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("device code prompt interrupted: %w", err)
		}
		instructions := message.Message
		if instructions == "" {
			if message.VerificationURL == "" || message.UserCode == "" {
				return errDeviceCodeInstructions
			}
			instructions = fmt.Sprintf("To sign in, open %s and enter the code %s.", message.VerificationURL, message.UserCode)
		}
		// Terraform captures provider stdout/stderr. Open the user's controlling
		// terminal separately so the challenge is visible while authentication polls.
		terminal, err := openTerminal()
		if err != nil {
			return fmt.Errorf("%w: %w", errDeviceCodeTerminalUnavailable, err)
		}
		_, writeErr := fmt.Fprintln(terminal, instructions)
		closeErr := terminal.Close()
		if writeErr != nil {
			return fmt.Errorf("display device code on the terminal: %w", writeErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close device code terminal: %w", closeErr)
		}
		return nil
	}
}
