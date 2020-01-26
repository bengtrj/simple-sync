package command

import (
	"fmt"

	"github.com/bengtrj/simple-sync/config"
	"github.com/bengtrj/simple-sync/sshclient"
)

//Sync runs the synchronization logic
func Sync(config *config.Sync) error {

	for _, server := range config.DesiredState.Servers {
		for _, app := range config.DesiredState.Apps {

			fmt.Printf("Updating/installing app %s\n", app.Name)
			address := fmt.Sprintf("%s:%s", server.IP, "22")
			client, err := sshclient.DialWithPasswd(address, config.User, config.Password)
			if err != nil {
				return err
			}
			defer client.Close()

			err = installPackages(client, app)
			if err != nil {
				return err
			}

		}
	}

	return nil

}

// Installs all packages
// Idempontent since apt-get won't change already installed packages
func installPackages(client *sshclient.Client, desiredAppState config.App) error {

	err := aptUpdate(client)
	if err != nil {
		return err
	}

	for _, p := range desiredAppState.Packages {
		fmt.Printf("Installing new package %s\n", p.Name)

		cmd := fmt.Sprintf("sudo apt-get install %s -y", p.Name)
		_, err := client.Cmd(cmd).SmartOutput()
		if err != nil {
			return err
		}
	}

	return nil
}

func aptUpdate(client *sshclient.Client) error {
	_, err := client.Cmd("sudo apt-get update -y").SmartOutput()
	if err != nil {
		return err
	}
	return nil
}
