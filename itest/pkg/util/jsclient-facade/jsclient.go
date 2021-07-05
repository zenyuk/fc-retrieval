package jsclient_facade

import (
	"fmt"
	"os"
	"os/exec"
)

func CallClientJsInstall() error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "../../../client-js/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func CallClientJsE2E(key string, walletKey string, registerAPI string, lotusAP string, lotusToken string) error {
	cmd := exec.Command("npm", "run", "test-e2e")
	if os.Getenv("RELOAD_JS_TESTS") == "yes" {
		cmd = exec.Command("npm", "run", "test-e2e-watch")
	}
	cmd.Dir = "../../../client-js"

	cmd.Env = append(os.Environ(),
		"ESTABLISHMENT_TTL=101",
		fmt.Sprintf("FCR_BLOCKCHAIN_PUBLIC_KEY=%s", key),
		fmt.Sprintf("FCR_REGISTER_API_URL=%s", registerAPI),
		fmt.Sprintf("FCR_WALLET_PRIVATE_KEY=%s", walletKey),
		fmt.Sprintf("FCR_LOTUS_AP=%s", lotusAP),
		fmt.Sprintf("FCR_LOTUS_AUTH_TOKEN=%s", lotusToken),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
