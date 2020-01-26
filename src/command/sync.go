package command

import (
	"fmt"
	"log"

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

			err = stopServices(client, app)
			if err != nil {
				return err
			}

			err = syncFiles(client, app)
			if err != nil {
				return err
			}

			err = installPackages(client, app)
			if err != nil {
				return err
			}

			err = startServices(client, app)
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
		out, err := client.Cmd(cmd).SmartOutput()
		if err != nil {
			log.Print(string(out))
			return err
		}
	}

	return nil
}

// Syncs all files
func syncFiles(client *sshclient.Client, desiredAppState config.App) error {

	for _, file := range desiredAppState.Files {
		fmt.Printf("Creating file %s\n", file.Path)
		err := copyFile(client, file)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copies one file
// For simplicity, it will always override the remote files
func copyFile(client *sshclient.Client, file config.File) error {

	template := `
path="%s"
owner=%s
group=%s
mode=%d

cat <<EOF | sudo tee ${path}
%s
EOF

sudo chmod ${mode} ${path}
sudo chown ${owner}:${group} ${path}
`

	script := fmt.Sprintf(template,
		file.Path,
		file.Owner,
		file.Group,
		file.Mode,
		file.Content)
	out, err := client.Script(script).SmartOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}

	return nil
}

// Starts the services based on which packages are services
// For simplicity, I'm not checking if service is running, just restarting them
// because it works regardless of the service is running or not
func startServices(client *sshclient.Client, app config.App) error {
	for _, p := range app.Packages {
		if p.IsService {
			script := fmt.Sprintf("sudo service %s restart", p.Name)
			out, err := client.Script(script).SmartOutput()
			if err != nil {
				log.Print(string(out))
				return err
			}
			fmt.Printf("Service %s started\n", p.Name)
		}
	}

	return nil
}

// Stops the services
func stopServices(client *sshclient.Client, app config.App) error {
	for _, p := range app.Packages {
		if p.IsService {
			script := fmt.Sprintf("sudo service %s stop", p.Name)
			_, err := client.Script(script).SmartOutput()
			if err != nil {
				return err
			}
			fmt.Printf("Service %s stopped\n", p.Name)
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
